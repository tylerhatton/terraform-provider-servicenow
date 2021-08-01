package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const serverName = "name"
const serverCompany = "company"
const serverAssetTag = "asset_tag"
const serverSerialNumber = "serial_number"
const serverManufacturer = "manufacturer"
const serverModelId = "model_id"
const serverAssignedTo = "assigned_to"
const serverOsDomain = "os_domain"
const serverRam = "ram"
const serverOperatingSystem = "operating_system"
const serverCpuManufacturer = "cpu_manufacturer"
const serverOsVersion = "os_version"
const serverCpuType = "cpu_type"
const serverOsServicePack = "os_service_pack"
const serverCpuSpeed = "cpu_speed"
const serverDnsDomain = "dns_domain"
const serverCpuCount = "cpu_count"
const serverDiskSpace = "disk_space"
const serverCpuCoreCount = "cpu_core_count"
const serverDescription = "description"
const serverIpAddress = "ip_address"

// ResourceServer manages a server cmdb entry in ServiceNow.
func ResourceServer() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_server` manages a server entry within ServiceNow.",

		Create: createResourceServer,
		Read:   readResourceServer,
		Update: updateResourceServer,
		Delete: deleteResourceServer,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			serverName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the server.",
			},
			serverCompany: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of company associated with server.",
			},
			serverAssetTag: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Asset tag associated with server.",
			},
			serverSerialNumber: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Serial number associated with server.",
			},
			serverManufacturer: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of manufacturer associated with server.",
			},
			serverModelId: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of model ID associated with server.",
			},
			serverAssignedTo: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of user associated with server.",
			},
			serverOsDomain: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of user associated with server.",
			},
			serverRam: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Amount of ram in MB associated with server.",
			},
			serverOperatingSystem: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Operating system associated with server.",
			},
			serverCpuManufacturer: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of CPU manufacturer associated with server.",
			},
			serverOsVersion: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Version of operating system associated with server.",
			},
			serverCpuType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "CPU type associated with server.",
			},
			serverOsServicePack: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Operating system service pack associated with server.",
			},
			serverCpuSpeed: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "CPU speed in MHz associated with server.",
			},
			serverDnsDomain: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "DNS Domain associated with server.",
			},
			serverCpuCount: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Number of CPUs associated with server.",
			},
			serverDiskSpace: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Disk space in GB associated with server.",
			},
			serverCpuCoreCount: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Number of CPU cores associated with server.",
			},
			serverDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Short description associated with server.",
			},
			serverIpAddress: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "IP Address associated with server.",
			},
		},
	}
}

func readResourceServer(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	server := &client.Server{}
	if err := snowClient.GetObject(client.EndpointServer, data.Id(), server); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServer(data, server)

	return nil
}

func createResourceServer(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	server := resourceToServer(data)
	if err := snowClient.CreateObject(client.EndpointServer, server); err != nil {
		return err
	}

	resourceFromServer(data, server)

	return readResourceServer(data, serviceNowClient)
}

func updateResourceServer(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointServer, resourceToServer(data)); err != nil {
		return err
	}

	return readResourceServer(data, serviceNowClient)
}

func deleteResourceServer(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointServer, data.Id())
}

func resourceFromServer(data *schema.ResourceData, server *client.Server) {
	data.SetId(server.ID)
	data.Set(serverName, server.Name)
	data.Set(serverCompany, server.Company)
	data.Set(serverAssetTag, server.AssetTag)
	data.Set(serverSerialNumber, server.SerialNumber)
	data.Set(serverManufacturer, server.Manufacturer)
	data.Set(serverModelId, server.ModelId)
	data.Set(serverAssignedTo, server.AssignedTo)
	data.Set(serverOsDomain, server.OsDomain)
	data.Set(serverRam, server.Ram)
	data.Set(serverOperatingSystem, server.OperatingSystem)
	data.Set(serverCpuManufacturer, server.CpuManufacturer)
	data.Set(serverOsVersion, server.OsVersion)
	data.Set(serverCpuType, server.CpuType)
	data.Set(serverOsServicePack, server.OsServicePack)
	data.Set(serverCpuSpeed, server.CpuSpeed)
	data.Set(serverDnsDomain, server.DnsDomain)
	data.Set(serverCpuCount, server.CpuCount)
	data.Set(serverDiskSpace, server.DiskSpace)
	data.Set(serverCpuCoreCount, server.CpuCoreCount)
	data.Set(serverDescription, server.Description)
	data.Set(serverIpAddress, server.IpAddress)
}

func resourceToServer(data *schema.ResourceData) *client.Server {
	server := client.Server{
		Name:            data.Get(serverName).(string),
		Company:         data.Get(serverCompany).(string),
		AssetTag:        data.Get(serverAssetTag).(string),
		SerialNumber:    data.Get(serverSerialNumber).(string),
		Manufacturer:    data.Get(serverManufacturer).(string),
		ModelId:         data.Get(serverModelId).(string),
		AssignedTo:      data.Get(serverAssignedTo).(string),
		OsDomain:        data.Get(serverOsDomain).(string),
		Ram:             data.Get(serverRam).(string),
		OperatingSystem: data.Get(serverOperatingSystem).(string),
		CpuManufacturer: data.Get(serverCpuManufacturer).(string),
		OsVersion:       data.Get(serverOsVersion).(string),
		CpuType:         data.Get(serverCpuType).(string),
		OsServicePack:   data.Get(serverOsServicePack).(string),
		CpuSpeed:        data.Get(serverCpuSpeed).(string),
		DnsDomain:       data.Get(serverDnsDomain).(string),
		CpuCount:        data.Get(serverCpuCount).(string),
		DiskSpace:       data.Get(serverDiskSpace).(string),
		CpuCoreCount:    data.Get(serverCpuCoreCount).(string),
		Description:     data.Get(serverDescription).(string),
		IpAddress:       data.Get(serverIpAddress).(string),
	}
	server.ID = data.Id()
	return &server
}
