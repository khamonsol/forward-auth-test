package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_AzureProvider(t *testing.T) {
	name := "test-azure-provider"
	config, err := NewConfig("AZURE_AUTH_PROVIDER", name)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, name, config.GetName())
}

func TestNewConfig_UnsupportedProvider(t *testing.T) {
	config, err := NewConfig("UNSUPPORTED_PROVIDER", "test")
	assert.Error(t, err)
	assert.Nil(t, config)
}
