package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAzureConfig_LoadConfig_Success(t *testing.T) {
	os.Setenv("AZURE_CLIENT_ID", "test-client-id")
	os.Setenv("AZURE_TENANT_ID", "test-tenant-id")
	os.Setenv("AZURE_ISSUER_URL", "https://issuer.url")
	defer os.Unsetenv("AZURE_CLIENT_ID")
	defer os.Unsetenv("AZURE_TENANT_ID")
	defer os.Unsetenv("AZURE_ISSUER_URL")

	config := &AzureConfig{Name: "test-azure-provider"}
	err := config.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "test-client-id", config.ClientID)
	assert.Equal(t, "test-tenant-id", config.TenantID)
	assert.Equal(t, "https://issuer.url", config.IssuerURL)
}

func TestAzureConfig_LoadConfig_MissingEnvVars(t *testing.T) {
	os.Unsetenv("AZURE_CLIENT_ID")
	os.Unsetenv("AZURE_TENANT_ID")
	os.Unsetenv("AZURE_ISSUER_URL")

	config := &AzureConfig{Name: "test-azure-provider"}
	err := config.LoadConfig()
	assert.Error(t, err)
}

func TestAzureConfig_GetName(t *testing.T) {
	config := &AzureConfig{Name: "test-azure-provider"}
	assert.Equal(t, "test-azure-provider", config.GetName())
}

func TestAzureConfig_GetIssuerURL(t *testing.T) {
	config := &AzureConfig{IssuerURL: "https://issuer.url"}
	assert.Equal(t, "https://issuer.url", config.GetIssuerURL())
}
