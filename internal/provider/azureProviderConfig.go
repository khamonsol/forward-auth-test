package provider

import (
	"errors"
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
)

const AZURE_AUTH_PROVIDER = "azureAuthProvider"

// AzureConfig holds the configuration necessary for Azure JWT provider.
// It includes TenantID, ClientID, ClientSecret, and the name of the ConfigMap to load these values from.
type AzureConfig struct {
	TenantID      string `yaml:"tenantId"`
	ClientID      string `yaml:"clientId"`
	ClientSecret  string `yaml:"clientSecret"`
	configMapName string
}

// NewAzureProviderConfig creates a new AzureConfig with the given ConfigMap name.
//
// Parameters:
//   - configMapName: The name of the ConfigMap to load the Azure provider configuration from.
//
// Returns:
//   - A pointer to the newly created AzureConfig.
func NewAzureProviderConfig(configMapName string) *AzureConfig {
	return &AzureConfig{
		configMapName: configMapName,
	}
}

// requiredKey retrieves a required key from the provided data map.
//
// Parameters:
//   - key: The key to retrieve from the data map.
//   - data: The map containing the key-value pairs.
//
// Returns:
//   - The value associated with the key.
//   - An error if the key is not found or the value is empty.
func requiredKey(key string, data map[string]string) (string, error) {
	val := data[key]
	if val == "" {
		return "", errors.New(fmt.Sprintf("Required key %s missing", key))
	}
	return val, nil
}

// LoadProviderConfig loads the Azure provider configuration from the specified ConfigMap.
//
// Parameters:
//   - api: An instance of the KubeAPI to use for retrieving the ConfigMap.
//
// Returns:
//   - An error if there was an issue loading the ConfigMap or parsing the required keys.
func (provider AzureConfig) LoadProviderConfig(api *util.KubeAPI) error {

	cmap, err := api.GetConfigMap(provider.configMapName)

	if err != nil {
		return err
	}
	tenantId, err := requiredKey("tenantId", cmap)
	if err != nil {
		return err
	}
	clientId, err := requiredKey("clientId", cmap)
	if err != nil {
		return err
	}
	clientSecret := cmap["clientSecret"]

	provider.TenantID = tenantId
	provider.ClientID = clientId
	provider.ClientSecret = clientSecret
	return nil
}

func (provider AzureConfig) GetName() string {
	return AZURE_AUTH_PROVIDER
}

func (provider AzureConfig) GetIssuerUrl() string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/")
}
