package resources_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// TestMain is the entry point used by `go test -sweep=<region>` to run the
// sweeper functions registered below. When invoked without -sweep it behaves
// like a normal test run.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sweeperClient builds a ServiceNow API client from the standard acceptance
// test environment variables. Sweepers fail fast when the credentials are not
// available so we never accidentally hit a production instance.
func sweeperClient() (*client.Client, error) {
	instance := os.Getenv("SERVICENOW_INSTANCE_URL")
	user := os.Getenv("SERVICENOW_USERNAME")
	pass := os.Getenv("SERVICENOW_PASSWORD")
	if instance == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("SERVICENOW_INSTANCE_URL/USERNAME/PASSWORD must be set for sweepers")
	}
	return client.NewClient(instance, user, pass), nil
}

// sweeperHTTPClient returns an HTTP client with a sane timeout for sweeper
// list operations. We use raw HTTP here because the typed client only exposes
// single-record lookups, while sweepers need to enumerate by prefix.
func sweeperHTTPClient() *http.Client {
	return &http.Client{Timeout: 60 * time.Second}
}

// sweeperRecord is the minimal projection of a JSONv2 list response we need
// to locate orphaned records. Each ServiceNow endpoint returns a Records
// array; sys_id is always present and uniquely addresses the record.
type sweeperRecord struct {
	SysID string `json:"sys_id"`
}

// sweeperListResponse mirrors the JSONv2 list envelope returned by the
// ServiceNow API.
type sweeperListResponse struct {
	Records []sweeperRecord `json:"records"`
}

// listRecordsByQuery performs a raw GET against the JSONv2 API using the
// supplied encoded sysparm_query and returns the sys_id values that match.
// It returns an empty slice (and nil error) when there are no matching
// records — that is the success path for a clean sweeper run.
func listRecordsByQuery(c *client.Client, endpoint, query string) ([]string, error) {
	requestURL := fmt.Sprintf("%s/%s?JSONv2&sysparm_query=%s",
		c.BaseURL, endpoint, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build sweeper request for %s: %w", endpoint, err)
	}
	req.Header.Set("Authorization", c.Auth)
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := sweeperHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("sweeper list %s failed: %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sweeper read response from %s: %w", endpoint, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sweeper list %s returned %s: %s", endpoint, resp.Status, body)
	}

	var parsed sweeperListResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("sweeper parse response from %s: %w", endpoint, err)
	}

	ids := make([]string, 0, len(parsed.Records))
	for _, r := range parsed.Records {
		if r.SysID != "" {
			ids = append(ids, r.SysID)
		}
	}
	return ids, nil
}

// sweepByFieldPrefix lists records whose `field` starts with the given
// `prefix` and deletes each one through the normal client DeleteObject API.
// A delete failure on a single record is logged but does not abort the
// remaining sweep so a partially corrupt environment can still be cleaned.
func sweepByFieldPrefix(endpoint, field, prefix string) error {
	c, err := sweeperClient()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("%sSTARTSWITH%s", field, prefix)
	ids, err := listRecordsByQuery(c, endpoint, query)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		log.Printf("[INFO] sweeper: no %s records matched %s=%s*", endpoint, field, prefix)
		return nil
	}
	for _, id := range ids {
		log.Printf("[INFO] sweeper deleting %s id=%s (matched %s=%s*)", endpoint, id, field, prefix)
		if delErr := c.DeleteObject(context.Background(), endpoint, id); delErr != nil {
			log.Printf("[WARN] sweeper failed to delete %s id=%s: %v", endpoint, id, delErr)
		}
	}
	return nil
}

// sweepByNamePrefix is the common case: sweep records whose `name` column
// begins with the given prefix.
func sweepByNamePrefix(endpoint, prefix string) error {
	return sweepByFieldPrefix(endpoint, "name", prefix)
}

// sweepByTitlePrefix sweeps records whose `title` column begins with the
// given prefix (used by the service catalog endpoint).
func sweepByTitlePrefix(endpoint, prefix string) error {
	return sweepByFieldPrefix(endpoint, "title", prefix)
}

// init registers a sweeper for every resource that uses a deterministic
// `tf-*` / `TF Acc*` test-name prefix. Sweepers run via:
//
//	go test ./servicenow/resources -sweep=default
//
// Each sweeper is tolerant of "no records found" — that is the expected
// outcome when there are no orphaned test resources to clean up. The
// prefix match is the primary safety mechanism that keeps real customer
// data out of scope.
func init() {
	resource.AddTestSweepers("servicenow_alias", &resource.Sweeper{
		Name: "servicenow_alias",
		F: func(region string) error {
			return sweepByNamePrefix("sys_alias.do", "tf-acc-")
		},
	})

	resource.AddTestSweepers("servicenow_application", &resource.Sweeper{
		Name: "servicenow_application",
		F: func(region string) error {
			return sweepByNamePrefix("sys_app.do", "tf")
		},
	})

	resource.AddTestSweepers("servicenow_role", &resource.Sweeper{
		Name: "servicenow_role",
		F: func(region string) error {
			return sweepByNamePrefix("sys_user_role.do", "tf_acc_")
		},
	})

	resource.AddTestSweepers("servicenow_rest_message", &resource.Sweeper{
		Name: "servicenow_rest_message",
		F: func(region string) error {
			return sweepByNamePrefix("sys_rest_message.do", "tf-acc-")
		},
	})

	resource.AddTestSweepers("servicenow_script_include", &resource.Sweeper{
		Name: "servicenow_script_include",
		F: func(region string) error {
			return sweepByNamePrefix("sys_script_include.do", "TfAcc")
		},
	})

	resource.AddTestSweepers("servicenow_scripted_rest_api", &resource.Sweeper{
		Name: "servicenow_scripted_rest_api",
		F: func(region string) error {
			return sweepByNamePrefix("sys_ws_definition.do", "tf")
		},
	})

	resource.AddTestSweepers("servicenow_ui_macro", &resource.Sweeper{
		Name: "servicenow_ui_macro",
		F: func(region string) error {
			return sweepByNamePrefix("sys_ui_macro.do", "tf_acc_")
		},
	})

	resource.AddTestSweepers("servicenow_ui_page", &resource.Sweeper{
		Name: "servicenow_ui_page",
		F: func(region string) error {
			return sweepByNamePrefix("sys_ui_page.do", "tf_acc_")
		},
	})

	resource.AddTestSweepers("servicenow_ui_script", &resource.Sweeper{
		Name: "servicenow_ui_script",
		F: func(region string) error {
			return sweepByNamePrefix("sys_ui_script.do", "tf_acc_")
		},
	})

	resource.AddTestSweepers("servicenow_system_property", &resource.Sweeper{
		Name: "servicenow_system_property",
		F: func(region string) error {
			return sweepByNamePrefix("sys_properties.do", "tf.acc.")
		},
	})

	resource.AddTestSweepers("servicenow_system_property_category", &resource.Sweeper{
		Name: "servicenow_system_property_category",
		F: func(region string) error {
			return sweepByNamePrefix("sys_properties_category.do", "TF Acc")
		},
	})

	resource.AddTestSweepers("servicenow_extension_point", &resource.Sweeper{
		Name: "servicenow_extension_point",
		F: func(region string) error {
			return sweepByNamePrefix("sys_extension_point.do", "TF Acc")
		},
	})

	resource.AddTestSweepers("servicenow_service_catalog", &resource.Sweeper{
		Name: "servicenow_service_catalog",
		F: func(region string) error {
			return sweepByTitlePrefix("sc_catalog.do", "TF Acc")
		},
	})

	resource.AddTestSweepers("servicenow_widget", &resource.Sweeper{
		Name: "servicenow_widget",
		F: func(region string) error {
			return sweepByNamePrefix("sp_widget.do", "TF Acc")
		},
	})

	resource.AddTestSweepers("servicenow_widget_dependency", &resource.Sweeper{
		Name: "servicenow_widget_dependency",
		F: func(region string) error {
			return sweepByNamePrefix("sp_dependency.do", "tf-acc-")
		},
	})

	resource.AddTestSweepers("servicenow_oauth_entity", &resource.Sweeper{
		Name: "servicenow_oauth_entity",
		F: func(region string) error {
			return sweepByNamePrefix("oauth_entity.do", "tf-acc-")
		},
	})

	resource.AddTestSweepers("servicenow_basic_auth_credential", &resource.Sweeper{
		Name: "servicenow_basic_auth_credential",
		F: func(region string) error {
			return sweepByNamePrefix("basic_auth_credentials.do", "tf-acc-")
		},
	})

	resource.AddTestSweepers("servicenow_http_connection", &resource.Sweeper{
		Name: "servicenow_http_connection",
		F: func(region string) error {
			return sweepByNamePrefix("http_connection.do", "tf-acc-")
		},
	})

	resource.AddTestSweepers("servicenow_db_table", &resource.Sweeper{
		Name: "servicenow_db_table",
		F: func(region string) error {
			return sweepByNamePrefix("sys_db_object.do", "u_tf")
		},
	})
}
