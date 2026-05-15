package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const (
	recordTable  = "table"
	recordFields = "fields"
	recordOutput = "output"
)

// ResourceRecord manages a row in an arbitrary ServiceNow table.
//
// Unlike the rest of the provider which targets one fixed table per resource,
// servicenow_record works against any table — built-in (incident, sys_user,
// cmdb_ci_server, ...) or custom (u_*). Field values are passed as a
// map(string) so the resource also accepts columns that are not known at
// compile time (e.g. u_* columns added via servicenow_dictionary).
func ResourceRecord() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a row in any ServiceNow table. Works with built-in tables (incident, sys_user, ...) and custom tables, including custom columns. Use the `fields` map to declare the columns Terraform owns; ServiceNow-populated columns (number, sys_created_on, etc.) are exposed via `output`.",

		CreateContext: createResourceRecord,
		ReadContext:   readResourceRecord,
		UpdateContext: updateResourceRecord,
		DeleteContext: deleteResourceRecord,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceRecord,
		},

		Schema: map[string]*schema.Schema{
			recordTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ServiceNow table name to create the record in (e.g. `incident`, `sys_user_group`, `u_my_table`).",
			},
			commonScope: getScopeSchema(),
			recordFields: {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Columns Terraform owns on this record, as a map of column name → string value. All values are serialized as strings on the wire — quote booleans (`\"true\"`) and integers (`\"42\"`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			recordOutput: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Every other column ServiceNow returned for this record (auto-generated numbers, computed fields, system metadata). Read-only.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func createResourceRecord(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	snowClient := meta.(client.ServiceNowClient)
	table := data.Get(recordTable).(string)
	scope := data.Get(commonScope).(string)
	fields := stringMap(data.Get(recordFields))

	result, err := snowClient.CreateRecord(ctx, table, scope, fields)
	if err != nil {
		return diag.FromErr(err)
	}
	if id, ok := result["sys_id"]; ok && id != "" {
		data.SetId(id)
	} else {
		return diag.Errorf("ServiceNow create succeeded but the response did not contain a sys_id")
	}
	return splitFieldsAndOutput(data, fields, result)
}

func readResourceRecord(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	snowClient := meta.(client.ServiceNowClient)
	table := data.Get(recordTable).(string)

	result, err := snowClient.GetRecord(ctx, table, data.Id())
	if err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	ownedFields := stringMap(data.Get(recordFields))
	return splitFieldsAndOutput(data, ownedFields, result)
}

func updateResourceRecord(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	snowClient := meta.(client.ServiceNowClient)
	table := data.Get(recordTable).(string)

	plannedFields := stringMap(data.Get(recordFields))

	// When the user removes a key from the `fields` map, we still need to tell
	// ServiceNow to clear it. Compute the diff against the prior state.
	body := make(map[string]string, len(plannedFields))
	for k, v := range plannedFields {
		body[k] = v
	}
	if data.HasChange(recordFields) {
		oldRaw, newRaw := data.GetChange(recordFields)
		oldMap := stringMap(oldRaw)
		newMap := stringMap(newRaw)
		for k := range oldMap {
			if _, stillOwned := newMap[k]; !stillOwned {
				body[k] = ""
			}
		}
	}

	result, err := snowClient.UpdateRecord(ctx, table, data.Id(), body)
	if err != nil {
		return diag.FromErr(err)
	}
	return splitFieldsAndOutput(data, plannedFields, result)
}

func deleteResourceRecord(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	snowClient := meta.(client.ServiceNowClient)
	table := data.Get(recordTable).(string)
	if err := snowClient.DeleteObject(ctx, table+".do", data.Id()); err != nil {
		// Tolerate "already gone" on destroy.
		if client.IsNotFound(err) {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}

// importResourceRecord supports `terraform import servicenow_record.x <table>/<sys_id>`.
// It seeds the `table` attribute and the resource ID so the subsequent Read
// can populate `output`. `fields` will start empty — the user adds columns
// they want Terraform to manage afterwards.
func importResourceRecord(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(data.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("import id must be in the form <table>/<sys_id>, got %q", data.Id())
	}
	if err := data.Set(recordTable, parts[0]); err != nil {
		return nil, err
	}
	data.SetId(parts[1])
	return []*schema.ResourceData{data}, nil
}

// splitFieldsAndOutput refreshes both the user-owned `fields` map and the
// computed `output` map from the server response. Keys present in `owned`
// land in `fields`; everything else (minus the sys_id which is the resource
// id) lands in `output`.
func splitFieldsAndOutput(data *schema.ResourceData, owned map[string]string, server map[string]string) diag.Diagnostics {
	fields := make(map[string]interface{}, len(owned))
	output := make(map[string]interface{}, len(server))
	for k, v := range server {
		if _, isOwned := owned[k]; isOwned {
			fields[k] = v
			continue
		}
		output[k] = v
	}
	// Preserve any owned keys the server response did not include (rare, but
	// keeps state stable if ServiceNow ever returns a sparse record).
	for k, v := range owned {
		if _, present := fields[k]; !present {
			fields[k] = v
		}
	}
	if err := data.Set(recordFields, fields); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set(recordOutput, output); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// stringMap coerces a schema.TypeMap value (which the SDK gives us as
// map[string]interface{}) into a map[string]string.
func stringMap(v interface{}) map[string]string {
	raw, ok := v.(map[string]interface{})
	if !ok {
		return map[string]string{}
	}
	out := make(map[string]string, len(raw))
	for k, val := range raw {
		if s, ok := val.(string); ok {
			out[k] = s
		} else if val != nil {
			out[k] = fmt.Sprintf("%v", val)
		}
	}
	return out
}
