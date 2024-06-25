package provider

import (
	"fmt"
	"os"
)

type AzureConfig struct {
	Name      string
	ClientID  string
	TenantID  string
	IssuerURL string
}

func (c *AzureConfig) LoadConfig() error {
	c.ClientID = os.Getenv("AZURE_CLIENT_ID")
	c.TenantID = os.Getenv("AZURE_TENANT_ID")
	c.IssuerURL = os.Getenv("AZURE_ISSUER_URL")

	if c.ClientID == "" || c.TenantID == "" || c.IssuerURL == "" {
		return fmt.Errorf("missing required Azure configuration")
	}
	return nil
}

func (c *AzureConfig) GetName() string {
	return c.Name
}

func (c *AzureConfig) GetIssuerURL() string {
	return c.IssuerURL
}
