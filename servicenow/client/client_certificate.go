package client

// EndpointCertificate is the endpoint to manage certificate records.
const EndpointCertificate = "sys_certificate.do"

// Certificate represents a certificate record in ServiceNow.
type Certificate struct {
	BaseResult
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Format           string `json:"format"`
	Type             string `json:"type"`
	PEMCertificate   string `json:"pem_certificate,omitempty"`
	KeyStorePassword string `json:"key_store_password,omitempty"`
	KeyStore         string `json:"key_store,omitempty"`
	Expiration       string `json:"expires"`
	Subject          string `json:"subject"`
	Issuer           string `json:"issuer"`
	Active           bool   `json:"active,string"`
	ValidFrom        string `json:"valid_from"`
}
