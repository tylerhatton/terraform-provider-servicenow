package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	servicenow "github.com/tylerhatton/terraform-provider-servicenow/servicenow"
)

// testAccPreCheck validates required environment variables are set before running acceptance tests.
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if v := os.Getenv("SERVICENOW_INSTANCE_URL"); v == "" {
		t.Fatal("SERVICENOW_INSTANCE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICENOW_USERNAME"); v == "" {
		t.Fatal("SERVICENOW_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICENOW_PASSWORD"); v == "" {
		t.Fatal("SERVICENOW_PASSWORD must be set for acceptance tests")
	}
}

// testAccProviderFactories returns a map of provider factories for acceptance tests.
func testAccProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"servicenow": func() (*schema.Provider, error) {
			return servicenow.Provider(), nil
		},
	}
}

// providerBlock returns the provider configuration block for use in acceptance test configs.
func providerBlock() string {
	return fmt.Sprintf(`
provider "servicenow" {
  instance_url = %q
  username     = %q
  password     = %q
}
`,
		os.Getenv("SERVICENOW_INSTANCE_URL"),
		os.Getenv("SERVICENOW_USERNAME"),
		os.Getenv("SERVICENOW_PASSWORD"),
	)
}

// checkExists returns a resource.TestCheckFunc that verifies a resource exists in state with a non-empty ID.
func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is empty for: %s", resourceName)
		}
		return nil
	}
}

// checkDestroy returns a resource.TestCheckFunc that verifies all resources of the given type
// have been removed from state after destruction.
func checkDestroy(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			if rs.Primary.ID != "" {
				return fmt.Errorf(
					"resource %s with ID %s still exists in state after destroy",
					resourceType, rs.Primary.ID,
				)
			}
		}
		return nil
	}
}
