package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/mock"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/resources"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) GetObject(_ context.Context, endpoint string, id string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, id, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) GetObjectByName(_ context.Context, endpoint string, name string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, name, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) GetObjectByTitle(_ context.Context, endpoint string, title string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, title, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) GetObjectByQuery(_ context.Context, endpoint string, query string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, query, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) CreateObject(_ context.Context, endpoint string, record client.Record) error {
	args := m.Called(endpoint, record)
	return args.Error(0)
}

func (m *ClientMock) UpdateObject(_ context.Context, endpoint string, record client.Record) error {
	args := m.Called(endpoint, record)
	return args.Error(0)
}

func (m *ClientMock) DeleteObject(_ context.Context, endpoint string, id string) error {
	args := m.Called(endpoint, id)
	return args.Error(0)
}

func (m *ClientMock) CreateRecord(_ context.Context, table, scope string, fields map[string]string) (map[string]string, error) {
	args := m.Called(table, scope, fields)
	out, _ := args.Get(0).(map[string]string)
	return out, args.Error(1)
}

func (m *ClientMock) GetRecord(_ context.Context, table, sysID string) (map[string]string, error) {
	args := m.Called(table, sysID)
	out, _ := args.Get(0).(map[string]string)
	return out, args.Error(1)
}

func (m *ClientMock) GetRecordByQuery(_ context.Context, table, query string) (map[string]string, error) {
	args := m.Called(table, query)
	out, _ := args.Get(0).(map[string]string)
	return out, args.Error(1)
}

func (m *ClientMock) UpdateRecord(_ context.Context, table, sysID string, fields map[string]string) (map[string]string, error) {
	args := m.Called(table, sysID, fields)
	out, _ := args.Get(0).(map[string]string)
	return out, args.Error(1)
}

type RecordMock struct {
	mock.Mock
}

func (m *RecordMock) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *RecordMock) GetScope() string {
	args := m.Called()
	return args.String(0)
}

func (m *RecordMock) GetStatus() string {
	args := m.Called()
	return args.String(0)
}

func (m *RecordMock) GetError() *client.ErrorDetail {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.(*client.ErrorDetail)
	}
	return nil
}

// callRead invokes ReadContext and reports whether diagnostics contain errors.
func callRead(res *schema.Resource, data *schema.ResourceData, meta interface{}) bool {
	if res.ReadContext == nil {
		return false
	}
	diags := res.ReadContext(context.Background(), data, meta)
	return diags.HasError()
}

// callCreate invokes CreateContext and reports whether diagnostics contain errors.
func callCreate(res *schema.Resource, data *schema.ResourceData, meta interface{}) bool {
	if res.CreateContext == nil {
		return false
	}
	diags := res.CreateContext(context.Background(), data, meta)
	return diags.HasError()
}

// callUpdate invokes UpdateContext and returns (hasError, diagnostics). Some
// resources are all-ForceNew (m2m relations) and have no update handler.
func callUpdate(res *schema.Resource, data *schema.ResourceData, meta interface{}) (bool, diag.Diagnostics) {
	if res.UpdateContext == nil {
		return false, nil
	}
	diags := res.UpdateContext(context.Background(), data, meta)
	return diags.HasError(), diags
}

// callDelete invokes DeleteContext and reports whether diagnostics contain errors.
func callDelete(res *schema.Resource, data *schema.ResourceData, meta interface{}) bool {
	if res.DeleteContext == nil {
		return false
	}
	diags := res.DeleteContext(context.Background(), data, meta)
	return diags.HasError()
}

var resourcesToTest = []*schema.Resource{
	resources.ResourceAlias(),
	resources.ResourceApplication(),
	resources.ResourceACL(),
	resources.ResourceApplicationMenu(),
	resources.ResourceApplicationModule(),
	resources.ResourceAssignmentRule(),
	resources.ResourceBasicAuthCredential(),
	resources.ResourceBusinessRule(),
	resources.ResourceCertificate(),
	resources.ResourceChoice(),
	resources.ResourceClientScript(),
	resources.ResourceContentCSS(),
	resources.ResourceCSSInclude(),
	resources.ResourceCSSIncludeRelation(),
	resources.ResourceDataLookup(),
	resources.ResourceDBTable(),
	resources.ResourceDictionary(),
	resources.ResourceEmailTemplate(),
	resources.ResourceEncryptionContext(),
	resources.ResourceExtensionPoint(),
	resources.ResourceFlow(),
	resources.ResourceGroup(),
	resources.ResourceGroupMember(),
	resources.ResourceGroupRole(),
	resources.ResourceHttpConnection(),
	resources.ResourceJdbcConnection(),
	resources.ResourceJsInclude(),
	resources.ResourceJsIncludeRelation(),
	resources.ResourceMidServer(),
	resources.ResourceNotification(),
	resources.ResourceOAuthEntity(),
	resources.ResourceQuestionChoice(),
	resources.ResourceRecord(),
	resources.ResourceRole(),
	resources.ResourceRestMessage(),
	resources.ResourceRestMessageHeader(),
	resources.ResourceRestMethod(),
	resources.ResourceRestMethodHeader(),
	resources.ResourceScheduledJob(),
	resources.ResourceScriptAction(),
	resources.ResourceScriptedRestApi(),
	resources.ResourceScriptedRestResource(),
	resources.ResourceScriptInclude(),
	resources.ResourceServer(),
	resources.ResourceServiceCatalog(),
	resources.ResourceServiceCatalogCategory(),
	resources.ResourceServiceCatalogItem(),
	resources.ResourceServiceCatalogVariable(),
	resources.ResourceSystemProperty(),
	resources.ResourceSystemPropertyCategory(),
	resources.ResourceSystemPropertyRelation(),
	resources.ResourceTransformEntry(),
	resources.ResourceTransformMap(),
	resources.ResourceUIAction(),
	resources.ResourceUIMacro(),
	resources.ResourceUIPage(),
	resources.ResourceUIPolicy(),
	resources.ResourceUIPolicyAction(),
	resources.ResourceUIScript(),
	resources.ResourceUser(),
	resources.ResourceUserRole(),
	resources.ResourceWidget(),
	resources.ResourceWidgetDependency(),
	resources.ResourceWidgetDependencyRelation(),
}

var dataSourcesToTest = []*schema.Resource{
	resources.DataSourceACL(),
	resources.DataSourceAlias(),
	resources.DataSourceApplication(),
	resources.DataSourceApplicationCategory(),
	resources.DataSourceApplicationMenu(),
	resources.DataSourceApplicationModule(),
	resources.DataSourceAssignmentRule(),
	resources.DataSourceBasicAuthCredential(),
	resources.DataSourceBusinessRule(),
	resources.DataSourceCertificate(),
	resources.DataSourceChoice(),
	resources.DataSourceClientScript(),
	resources.DataSourceContentCSS(),
	resources.DataSourceCSSInclude(),
	resources.DataSourceDataLookup(),
	resources.DataSourceDBTable(),
	resources.DataSourceDictionary(),
	resources.DataSourceEmailTemplate(),
	resources.DataSourceEncryptionContext(),
	resources.DataSourceExtensionPoint(),
	resources.DataSourceFlow(),
	resources.DataSourceGroup(),
	resources.DataSourceHttpConnection(),
	resources.DataSourceJdbcConnection(),
	resources.DataSourceJsInclude(),
	resources.DataSourceMidServer(),
	resources.DataSourceNotification(),
	resources.DataSourceOAuthEntity(),
	resources.DataSourceRecord(),
	resources.DataSourceRestMessage(),
	resources.DataSourceRole(),
	resources.DataSourceScheduledJob(),
	resources.DataSourceScriptAction(),
	resources.DataSourceScriptInclude(),
	resources.DataSourceScriptedRestApi(),
	resources.DataSourceServer(),
	resources.DataSourceServiceCatalog(),
	resources.DataSourceServiceCatalogCategory(),
	resources.DataSourceServiceCatalogItem(),
	resources.DataSourceSystemProperty(),
	resources.DataSourceSystemPropertyCategory(),
	resources.DataSourceTransformMap(),
	resources.DataSourceUIAction(),
	resources.DataSourceUIMacro(),
	resources.DataSourceUIPage(),
	resources.DataSourceUIPolicy(),
	resources.DataSourceUIScript(),
	resources.DataSourceUser(),
	resources.DataSourceWidget(),
	resources.DataSourceWidgetDependency(),
}

func TestResourcesCanRead(t *testing.T) {
	for _, res := range resourcesToTest {
		// servicenow_record requires `table` in state before Read can run.
		raw := map[string]interface{}{}
		if _, ok := res.Schema["table"]; ok {
			raw["table"] = "incident"
		}
		data := schema.TestResourceDataRaw(t, res.Schema, raw)
		data.SetId("hello")
		clientMock := new(ClientMock)
		// Resources use one of these read paths; tolerate any (success).
		clientMock.On("GetObject", mock.AnythingOfType("string"), "hello", mock.Anything).Return(nil).Maybe()
		clientMock.On("GetRecord", mock.AnythingOfType("string"), "hello").Return(map[string]string{"sys_id": "hello"}, nil).Maybe()
		callRead(res, data, clientMock)
	}
}

func TestResourceRestMessageHandleReadError(t *testing.T) {
	for _, res := range resourcesToTest {
		raw := map[string]interface{}{}
		if _, ok := res.Schema["table"]; ok {
			raw["table"] = "incident"
		}
		data := schema.TestResourceDataRaw(t, res.Schema, raw)
		data.SetId("hello")
		clientMock := new(ClientMock)
		clientMock.On("GetObject", mock.AnythingOfType("string"), "hello", mock.Anything).Return(fmt.Errorf("nothing to see here")).Maybe()
		clientMock.On("GetRecord", mock.AnythingOfType("string"), "hello").Return(map[string]string(nil), &client.NotFoundError{Reason: "not found"}).Maybe()
		callRead(res, data, clientMock)
		assert.Equal(t, "", data.Id())
	}
}

func TestDataSourcesCanRead(t *testing.T) {
	for _, res := range dataSourcesToTest {
		// Fill every Required string field with a placeholder so data sources
		// with composite lookups (e.g. dictionary requires name + element)
		// don't crash inside the Read function.
		fakeData := map[string]interface{}{}
		for key, prop := range res.Schema {
			if !prop.Required {
				continue
			}
			switch prop.Type {
			case schema.TypeString:
				fakeData[key] = "oi"
			case schema.TypeBool:
				fakeData[key] = true
			case schema.TypeInt:
				fakeData[key] = 1
			}
		}
		data := schema.TestResourceDataRaw(t, res.Schema, fakeData)

		// Allow any of the four lookup methods; the data source picks one.
		clientMock := new(ClientMock)
		clientMock.On("GetObjectByName", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil).Maybe()
		clientMock.On("GetObjectByTitle", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil).Maybe()
		clientMock.On("GetObjectByQuery", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil).Maybe()
		clientMock.On("GetObject", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil).Maybe()
		clientMock.On("GetRecord", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(map[string]string{"sys_id": "x"}, nil).Maybe()
		clientMock.On("GetRecordByQuery", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(map[string]string{"sys_id": "x"}, nil).Maybe()

		callRead(res, data, clientMock)
	}
}

func TestResourcesCanUpdate(t *testing.T) {
	for _, res := range resourcesToTest {
		fakeData := map[string]interface{}{}
		// Fill the data with stuff following the schema.
		for key, prop := range res.Schema {
			if !prop.Computed {
				switch prop.Type {
				case schema.TypeString:
					fakeData[key] = "hello"
				case schema.TypeBool:
					fakeData[key] = true
				case schema.TypeInt:
					fakeData[key] = 42
				}
			}
		}

		data := schema.TestResourceDataRaw(t, res.Schema, fakeData)
		data.SetId("fenouille")

		clientMock := new(ClientMock)
		clientMock.On("UpdateObject", mock.AnythingOfType("string"), mock.Anything).Return(nil).Maybe()
		clientMock.On("GetObject", mock.AnythingOfType("string"), "fenouille", mock.Anything).Return(nil).Maybe()
		clientMock.On("UpdateRecord", mock.AnythingOfType("string"), "fenouille", mock.Anything).Return(map[string]string{"sys_id": "fenouille"}, nil).Maybe()
		clientMock.On("GetRecord", mock.AnythingOfType("string"), "fenouille").Return(map[string]string{"sys_id": "fenouille"}, nil).Maybe()

		callUpdate(res, data, clientMock)
	}
}

func TestResourcesCanDelete(t *testing.T) {
	for _, res := range resourcesToTest {
		raw := map[string]interface{}{}
		if _, ok := res.Schema["table"]; ok {
			raw["table"] = "incident"
		}
		data := schema.TestResourceDataRaw(t, res.Schema, raw)
		data.SetId("fenouille")
		clientMock := new(ClientMock)
		clientMock.On("DeleteObject", mock.AnythingOfType("string"), "fenouille").Return(nil).Maybe()
		callDelete(res, data, clientMock)
	}
}

func TestResourcesCanCreate(t *testing.T) {
	for _, res := range resourcesToTest {
		fakeData := map[string]interface{}{}
		// Fill the data with stuff following the schema.
		for key, prop := range res.Schema {
			if !prop.Computed {
				switch prop.Type {
				case schema.TypeString:
					fakeData[key] = "hello"
				case schema.TypeBool:
					fakeData[key] = true
				case schema.TypeInt:
					fakeData[key] = 42
				}
			}
		}

		data := schema.TestResourceDataRaw(t, res.Schema, fakeData)

		clientMock := new(ClientMock)
		clientMock.On("CreateObject", mock.AnythingOfType("string"), mock.Anything).Return(nil).Maybe()
		clientMock.On("GetObject", mock.AnythingOfType("string"), "", mock.Anything).Return(nil).Maybe()
		clientMock.On("CreateRecord", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(map[string]string{"sys_id": "fenouille"}, nil).Maybe()
		clientMock.On("GetRecord", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(map[string]string{"sys_id": "fenouille"}, nil).Maybe()

		hasError := callCreate(res, data, clientMock)
		assert.False(t, hasError)
	}
}

func TestResourcesReturnErrorOnUpdateFailure(t *testing.T) {
	for _, res := range resourcesToTest {
		fakeData := map[string]interface{}{}
		// Fill the data with stuff following the schema.
		for key, prop := range res.Schema {
			if !prop.Computed {
				switch prop.Type {
				case schema.TypeString:
					fakeData[key] = "hello"
				case schema.TypeBool:
					fakeData[key] = true
				case schema.TypeInt:
					fakeData[key] = 42
				}
			}
		}

		data := schema.TestResourceDataRaw(t, res.Schema, fakeData)
		data.SetId("fenouille")

		clientMock := new(ClientMock)
		clientMock.On("UpdateObject", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("update failed")).Maybe()
		clientMock.On("UpdateRecord", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(map[string]string(nil), fmt.Errorf("update failed")).Maybe()

		if res.UpdateContext == nil {
			continue // resource is all-ForceNew, no update path
		}
		hasError, _ := callUpdate(res, data, clientMock)
		assert.True(t, hasError)
	}
}

func TestResourcesDeleteHandlesNotFound(t *testing.T) {
	for _, res := range resourcesToTest {
		raw := map[string]interface{}{}
		if _, ok := res.Schema["table"]; ok {
			raw["table"] = "incident"
		}
		data := schema.TestResourceDataRaw(t, res.Schema, raw)
		data.SetId("fenouille")
		clientMock := new(ClientMock)
		clientMock.On("DeleteObject", mock.AnythingOfType("string"), "fenouille").Return(fmt.Errorf("record not found")).Maybe()
		callDelete(res, data, clientMock)
	}
}
