# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
make build       # Build terraform-provider-servicenow binary
make install     # Build and install to ~/.terraform.d/plugins/
make release     # Cross-compile for all supported platforms

# Test
make test        # Unit tests (4 parallel, 30s timeout)
make testacc     # Acceptance tests (requires TF_ACC=1, 120m timeout)
go test ./servicenow/resources/ -run TestResourceAlias -v  # Single test
```

There is no lint target in the Makefile. Tests use `testify/mock` for mocking the `ServiceNowClient` interface.

## Architecture

This is a Terraform Provider SDK v1 provider for the ServiceNow REST API. The codebase has three layers:

```
main.go                        # plugin.Serve() entry point
servicenow/
  provider.go                  # Provider schema, resource/datasource registration, client init
  client/
    client_base.go             # ServiceNowClient interface, Client struct, HTTP logic
    client_*.go                # One file per resource: endpoint constant + response struct
  resources/
    common.go                  # Shared schema helpers (setOnlyRequiredSchema, etc.)
    resource_*.go              # CRUD resource definitions
    data_source_*.go           # Read-only data source definitions
    resources_test.go          # All unit tests (mocked)
```

### ServiceNow API Client

`client/client_base.go` contains the full HTTP client implementation. All requests target `<instance_url>/<table>.do?JSONv2` using ServiceNow's legacy JSON API format:

- **Read**: `GET ?JSONv2&sysparm_query=sys_id={id}`
- **Create**: `POST ?JSONv2&sysparm_action=insert` with JSON body
- **Update**: `POST ?JSONv2&sysparm_action=update&sysparm_query=sys_id={id}` with JSON body
- **Delete**: `POST ?JSONv2&sysparm_action=deleteRecord&sysparm_sys_id={id}`

Authentication is HTTP Basic (base64-encoded). Responses have the shape `{"records": [...]}` and include a `__status` field that must equal `"success"`.

The `ServiceNowClient` interface enables mock substitution in tests. Individual client files contain only:
1. An endpoint constant (e.g., `EndpointAlias = "sys_alias.do"`)
2. A struct embedding `BaseResult` with JSON tags matching ServiceNow field names

### Resource Pattern

Every resource follows this exact pattern. To add a new resource, create three files:

**`client/client_thing.go`**:
```go
const EndpointThing = "sys_thing.do"
type Thing struct {
    BaseResult
    Name string `json:"name"`
    // ServiceNow boolean fields use the `,string` tag:
    Active bool `json:"active,string"`
}
```

**`resources/resource_thing.go`**:
```go
func ResourceThing() *schema.Resource { ... }  // schema definition
func createThing(data *schema.ResourceData, serviceNowClient interface{}) error { ... }
func readThing(data *schema.ResourceData, serviceNowClient interface{}) error { ... }
func updateThing(data *schema.ResourceData, serviceNowClient interface{}) error { ... }
func deleteThing(data *schema.ResourceData, serviceNowClient interface{}) error { ... }
func resourceFromThing(data *schema.ResourceData, thing *client.Thing) { ... }  // API → state
func resourceToThing(data *schema.ResourceData) *client.Thing { ... }            // state → API
```

**`resources/data_source_thing.go`** (optional):
- Reuses `ResourceThing().Schema` but wraps all fields with `setOnlyRequiredSchema()`
- Uses `client.GetObjectByName` for lookup
- Only implements a `read` handler

Register both in `servicenow/provider.go` under `ResourcesMap` and `DataSourcesMap`.

### Boolean Fields

ServiceNow returns booleans as `"true"`/`"false"` strings. Always use the `,string` tag:
```go
Active bool `json:"active,string"`
```

### Scoped Applications

Resources in a specific application scope include `sys_scope` in their struct and pass `sysparm_record_scope` as a query parameter during create. The client handles this automatically when `GetScope()` returns a non-empty value.
