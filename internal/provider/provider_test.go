package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	providerConfig := NewConfig("azure", "azure_beyond_prod", "default")
	assert.NotNil(t, providerConfig)
	assert.Equal(t, "azure_beyond_prod", providerConfig.GetName())

	azureConfig := providerConfig.(*AzureProviderConfig)
	azureConfig.TenantID = "your-tenant-id"
	azureConfig.IssuerURL = "https://login.microsoftonline.com/your-tenant-id/v2.0"

	assert.Equal(t, "https://login.microsoftonline.com/your-tenant-id/v2.0", providerConfig.GetIssuerUrl())
}
