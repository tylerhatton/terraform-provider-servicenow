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

func (m *ClientMock) GetObject(endpoint string, id string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, id, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) GetObjectByName(endpoint string, name string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, name, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) GetObjectByTitle(endpoint string, title string, responseObjectOut client.Record) error {
	args := m.Called(endpoint, title, responseObjectOut)
	return args.Error(0)
}

func (m *ClientMock) CreateObject(endpoint string, record client.Record) error {
	args := m.Called(endpoint, record)
	return args.Error(0)
}

func (m *ClientMock) UpdateObject(endpoint string, record client.Record) error {
	args := m.Called(endpoint, record)
	return args.Error(0)
}

func (m *ClientMock) DeleteObject(endpoint string, id string) error {
	args := m.Called(endpoint, id)
	return args.Error(0)
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

// callRead invokes ReadContext if set, otherwise falls back to Read.
// Returns true if any diagnostics have errors (or the error return is non-nil).
func callRead(res *schema.Resource, data *schema.ResourceData, meta interface{}) bool {
	if res.ReadContext != nil {
		diags := res.ReadContext(context.Background(), data, meta)
		return diags.HasError()
	}
	err := res.Read(data, meta)
	return err != nil
}

// callCreate invokes CreateContext if set, otherwise falls back to Create.
func callCreate(res *schema.Resource, data *schema.ResourceData, meta interface{}) bool {
	if res.CreateContext != nil {
		diags := res.CreateContext(context.Background(), data, meta)
		return diags.HasError()
	}
	err := res.Create(data, meta)
	return err != nil
}

// callUpdate invokes UpdateContext if set, otherwise falls back to Update.
func callUpdate(res *schema.Resource, data *schema.ResourceData, meta interface{}) (bool, diag.Diagnostics) {
	if res.UpdateContext != nil {
		diags := res.UpdateContext(context.Background(), data, meta)
		return diags.HasError(), diags
	}
	err := res.Update(data, meta)
	return err != nil, nil
}

// callDelete invokes DeleteContext if set, otherwise falls back to Delete.
func callDelete(res *schema.Resource, data *schema.ResourceData, meta interface{}) bool {
	if res.DeleteContext != nil {
		diags := res.DeleteContext(context.Background(), data, meta)
		return diags.HasError()
	}
	err := res.Delete(data, meta)
	return err != nil
}

var resourcesToTest = []*schema.Resource{
	resources.ResourceAlias(),
	resources.ResourceApplication(),
	resources.ResourceApplicationMenu(),
	resources.ResourceApplicationModule(),
	resources.ResourceBasicAuthCredential(),
	resources.ResourceContentCSS(),
	resources.ResourceCSSInclude(),
	resources.ResourceCSSIncludeRelation(),
	resources.ResourceDBTable(),
	resources.ResourceExtensionPoint(),
	resources.ResourceHttpConnection(),
	resources.ResourceJsInclude(),
	resources.ResourceJsIncludeRelation(),
	resources.ResourceOAuthEntity(),
	resources.ResourceQuestionChoice(),
	resources.ResourceRole(),
	resources.ResourceRestMessage(),
	resources.ResourceRestMessageHeader(),
	resources.ResourceRestMethod(),
	resources.ResourceRestMethodHeader(),
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
	resources.ResourceUIMacro(),
	resources.ResourceUIPage(),
	resources.ResourceUIScript(),
	resources.ResourceWidget(),
	resources.ResourceWidgetDependency(),
	resources.ResourceWidgetDependencyRelation(),
}

var dataSourcesToTest = []*schema.Resource{
	resources.DataSourceACL(),
	resources.DataSourceApplication(),
	resources.DataSourceApplicationCategory(),
	resources.DataSourceDBTable(),
	resources.DataSourceRole(),
	resources.DataSourceServiceCatalog(),
	resources.DataSourceServiceCatalogCategory(),
	resources.DataSourceSystemProperty(),
	resources.DataSourceSystemPropertyCategory(),
}

func TestResourcesCanRead(t *testing.T) {
	for _, res := range resourcesToTest {
		data := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		data.SetId("hello")
		clientMock := new(ClientMock)
		clientMock.
			On("GetObject", mock.AnythingOfType("string"), "hello", mock.Anything).
			Return(nil)

		callRead(res, data, clientMock)
		clientMock.AssertExpectations(t)
	}
}

func TestResourceRestMessageHandleReadError(t *testing.T) {
	for _, res := range resourcesToTest {
		data := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		data.SetId("hello")
		clientMock := new(ClientMock)
		clientMock.
			On("GetObject", mock.AnythingOfType("string"), "hello", mock.Anything).
			Return(fmt.Errorf("nothing to see here"))

		callRead(res, data, clientMock)
		clientMock.AssertExpectations(t)
		assert.Equal(t, "", data.Id())
	}
}

func TestDataSourcesCanRead(t *testing.T) {
	for _, res := range dataSourcesToTest {
		_, hasName := res.Schema["name"]
		_, hasTitle := res.Schema["title"]

		var data *schema.ResourceData
		clientMock := new(ClientMock)

		if hasTitle && !hasName {
			data = schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"title": "oi",
			})
			clientMock.
				On("GetObjectByTitle", mock.AnythingOfType("string"), "oi", mock.Anything).
				Return(nil)
		} else {
			data = schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"name": "oi",
			})
			clientMock.
				On("GetObjectByName", mock.AnythingOfType("string"), "oi", mock.Anything).
				Return(nil)
		}

		callRead(res, data, clientMock)
		clientMock.AssertExpectations(t)
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
		clientMock.
			On("UpdateObject", mock.AnythingOfType("string"), mock.Anything).
			Return(nil)
		clientMock.
			On("GetObject", mock.AnythingOfType("string"), "fenouille", mock.Anything).
			Return(nil)

		callUpdate(res, data, clientMock)
		clientMock.AssertExpectations(t)
	}
}

func TestResourcesCanDelete(t *testing.T) {
	for _, res := range resourcesToTest {
		data := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		data.SetId("fenouille")
		clientMock := new(ClientMock)
		clientMock.
			On("DeleteObject", mock.AnythingOfType("string"), "fenouille").
			Return(nil)

		callDelete(res, data, clientMock)
		clientMock.AssertExpectations(t)
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
		clientMock.
			On("CreateObject", mock.AnythingOfType("string"), mock.Anything).
			Return(nil)
		clientMock.
			On("GetObject", mock.AnythingOfType("string"), "", mock.Anything).
			Return(nil)

		hasError := callCreate(res, data, clientMock)
		clientMock.AssertExpectations(t)
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
		clientMock.
			On("UpdateObject", mock.AnythingOfType("string"), mock.Anything).
			Return(fmt.Errorf("update failed"))

		hasError, _ := callUpdate(res, data, clientMock)
		clientMock.AssertExpectations(t)
		assert.True(t, hasError)
	}
}

func TestResourcesDeleteHandlesNotFound(t *testing.T) {
	for _, res := range resourcesToTest {
		data := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		data.SetId("fenouille")
		clientMock := new(ClientMock)
		clientMock.
			On("DeleteObject", mock.AnythingOfType("string"), "fenouille").
			Return(fmt.Errorf("record not found"))

		// When DeleteObject returns an error, the delete function propagates it.
		// We verify the mock was called correctly and the function did not panic.
		callDelete(res, data, clientMock)
		clientMock.AssertExpectations(t)
	}
}
