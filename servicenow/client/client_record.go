package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// envelopeStatusKey is the JSONv2 status field returned in every response.
const envelopeStatusKey = "__status"

// envelopeErrorKey is the JSONv2 error detail field, present on failures.
const envelopeErrorKey = "__error"

// CreateRecord inserts a row into the given table and returns every column of the
// server-rendered record. `scope` is the application scope sys_id ("" or "global"
// for the global scope). `fields` carries the columns the caller wants to set;
// every other column will be defaulted by ServiceNow and reflected back in the
// returned map.
func (client *Client) CreateRecord(ctx context.Context, table, scope string, fields map[string]string) (map[string]string, error) {
	path := table + ".do?JSONv2&sysparm_action=insert"
	if scope != "" && scope != "global" {
		path += "&sysparm_record_scope=" + url.QueryEscape(scope)
	}
	body := cloneFields(fields)
	raw, err := client.requestJSON(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}
	return decodeSingleRecord(raw)
}

// GetRecord fetches the row identified by sys_id and returns all of its columns.
// Returns *NotFoundError when no record exists for that sys_id.
func (client *Client) GetRecord(ctx context.Context, table, sysID string) (map[string]string, error) {
	path := table + ".do?JSONv2&sysparm_query=sys_id=" + url.QueryEscape(sysID)
	raw, err := client.requestJSON(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return decodeSingleRecord(raw)
}

// GetRecordByQuery fetches the single row matching an encoded sysparm_query.
// Returns *NotFoundError on zero matches and an error on multi-match.
func (client *Client) GetRecordByQuery(ctx context.Context, table, query string) (map[string]string, error) {
	path := table + ".do?JSONv2&sysparm_query=" + url.QueryEscape(query)
	raw, err := client.requestJSON(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return decodeSingleRecord(raw)
}

// UpdateRecord patches the row identified by sys_id with the given fields and
// returns the post-update column values. Fields that are not in `fields` are
// left untouched by ServiceNow; pass an empty string to clear a column.
func (client *Client) UpdateRecord(ctx context.Context, table, sysID string, fields map[string]string) (map[string]string, error) {
	path := table + ".do?JSONv2&sysparm_action=update&sysparm_query=sys_id=" + url.QueryEscape(sysID)
	body := cloneFields(fields)
	raw, err := client.requestJSON(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}
	return decodeSingleRecord(raw)
}

// cloneFields makes a defensive copy of the caller's map so we don't mutate it
// when injecting envelope-level keys.
func cloneFields(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// decodeSingleRecord parses a JSONv2 envelope expecting exactly one record and
// returns it as a flat string→string map after validating the __status field.
func decodeSingleRecord(raw []byte) (map[string]string, error) {
	envelope := struct {
		Records []map[string]json.RawMessage `json:"records"`
	}{}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("decode JSONv2 envelope: %w", err)
	}
	if len(envelope.Records) == 0 {
		return nil, &NotFoundError{Reason: "no records found"}
	}
	if len(envelope.Records) > 1 {
		return nil, fmt.Errorf("more than one record received")
	}

	rec := envelope.Records[0]

	// Status validation: __status must equal "success".
	if statusRaw, ok := rec[envelopeStatusKey]; ok {
		var status string
		if err := json.Unmarshal(statusRaw, &status); err == nil && status != "" && status != "success" {
			detail := unwrapErrorDetail(rec[envelopeErrorKey])
			return nil, fmt.Errorf("error from ServiceNow -> %s: %s", detail.Message, detail.Reason)
		}
	}

	out := make(map[string]string, len(rec))
	for k, v := range rec {
		if k == envelopeStatusKey || k == envelopeErrorKey {
			continue
		}
		// Most fields are strings; references may rarely appear as objects when
		// display_value is enabled (we do not set it). Fall back gracefully.
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			out[k] = s
			continue
		}
		// Non-string value: re-marshal so the caller can still inspect it.
		out[k] = string(v)
	}
	return out, nil
}

// unwrapErrorDetail best-effort decodes a JSONv2 __error payload.
func unwrapErrorDetail(raw json.RawMessage) ErrorDetail {
	if len(raw) == 0 {
		return ErrorDetail{Message: "unknown", Reason: ""}
	}
	var detail ErrorDetail
	_ = json.Unmarshal(raw, &detail)
	return detail
}
