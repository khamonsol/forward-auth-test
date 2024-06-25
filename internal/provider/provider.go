package provider

import (
	"github.com/SoleaEnergy/forwardAuth/internal/util"
)

type Provider interface {
	LoadProviderConfig(api util.KubernetesClient) error
	GetName() string
	GetIssuerUrl() string
}

func NewConfig(providerType string, name string) Provider {
	switch providerType {
	case "azure":
		return &AzureProviderConfig{
			Name:      name
		}
	default:
		return nil
	}
}
