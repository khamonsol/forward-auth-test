package provider

import (
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"log/slog"
)

// Config defines an interface for loading provider configurations.
//
// Implementations of this interface should provide a method to load the provider
// configuration using a KubeAPI instance.
type Config interface {
	// LoadProviderConfig loads the provider configuration from a Kubernetes ConfigMap.
	//
	// Parameters:
	//   - api: An instance of the KubeAPI to use for retrieving the ConfigMap.
	//
	// Returns:
	//   - An error if there was an issue loading the ConfigMap or parsing the required keys.
	LoadProviderConfig(*util.KubeAPI) error

	GetName() string

	GetIssuerUrl() string
}

func NewConfig(providerType string, name string) Config {
	switch providerType {
	case AZURE_AUTH_PROVIDER:
		return NewAzureProviderConfig(name)
	default:
		slog.Error("Supplied provider is not yet implemented")
		return nil
	}
}
