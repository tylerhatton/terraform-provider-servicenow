package client

// EndpointJdbcConnection is the endpoint to manage JDBC connection records.
const EndpointJdbcConnection = "jdbc_connection.do"

// JdbcConnection represents a JDBC connection record in ServiceNow. JDBC connections
// extend sys_connection and are used by data sources / integrations that need a
// database link, optionally through a MID server.
type JdbcConnection struct {
	BaseResult
	Name            string `json:"name"`
	Active          bool   `json:"active,string"`
	Credential      string `json:"credential"`
	ConnectionAlias string `json:"connection_alias"`
	ConnectionUrl   string `json:"connection_url"`
	DatabaseName    string `json:"database_name"`
	DatabaseType    string `json:"jdbc_driver"`
	UseMidServer    bool   `json:"use_mid,string"`
	MidSelection    string `json:"mid_selection"`
	MidServer       string `json:"mid_server"`
}
