package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const (
	dataSourceRecordSysID = "sys_id"
	dataSourceRecordQuery = "query"
)

// DataSourceRecord looks up a single row in any ServiceNow table.
//
// Provide either `sys_id` for a direct lookup, or `query` for an encoded
// sysparm_query expression that must match exactly one record.
func DataSourceRecord() *schema.Resource {
	return &schema.Resource{
		Description: "Look up a single row in any ServiceNow table, either by sys_id or by an encoded sysparm_query that must match exactly one record. Returns all of the record's columns in `output`.",

		ReadContext: readDataSourceRecord,

		Schema: map[string]*schema.Schema{
			recordTable: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ServiceNow table to look up (e.g. `incident`, `sys_user_group`, `u_my_table`).",
			},
			dataSourceRecordSysID: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "sys_id of the record to look up. Mutually exclusive with `query`.",
				ExactlyOneOf: []string{dataSourceRecordSysID, dataSourceRecordQuery},
			},
			dataSourceRecordQuery: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Encoded sysparm_query expression (e.g. `name=Admin^active=true`). Must match exactly one record.",
				ExactlyOneOf: []string{dataSourceRecordSysID, dataSourceRecordQuery},
			},
			recordOutput: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "All columns ServiceNow returned for the matched record, as a map of column name → string value.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func readDataSourceRecord(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	snowClient := meta.(client.ServiceNowClient)
	table := data.Get(recordTable).(string)

	var (
		result map[string]string
		err    error
	)
	if sysID := data.Get(dataSourceRecordSysID).(string); sysID != "" {
		result, err = snowClient.GetRecord(ctx, table, sysID)
	} else {
		query := data.Get(dataSourceRecordQuery).(string)
		result, err = snowClient.GetRecordByQuery(ctx, table, query)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	sysID, ok := result["sys_id"]
	if !ok || sysID == "" {
		return diag.Errorf("ServiceNow returned a record without a sys_id")
	}
	data.SetId(sysID)
	if err := data.Set(dataSourceRecordSysID, sysID); err != nil {
		return diag.FromErr(err)
	}

	output := make(map[string]interface{}, len(result))
	for k, v := range result {
		output[k] = v
	}
	if err := data.Set(recordOutput, output); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
