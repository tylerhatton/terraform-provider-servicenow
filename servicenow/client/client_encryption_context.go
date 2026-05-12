package client

// EndpointEncryptionContext is the endpoint to manage encryption context records.
const EndpointEncryptionContext = "sys_encryption_context.do"

// EncryptionContext represents an encryption context record in ServiceNow Edge
// Encryption. This table is only present on ServiceNow instances that have the
// Edge Encryption plugin installed and activated. On instances without the
// plugin, the underlying endpoint will return an "Invalid table" error.
type EncryptionContext struct {
	BaseResult
	Name          string `json:"name"`
	Type          string `json:"type"`
	EncryptionKey string `json:"encryption_key,omitempty"`
	Description   string `json:"description"`
	Active        bool   `json:"active,string"`
	Algorithm     string `json:"algorithm"`
}
