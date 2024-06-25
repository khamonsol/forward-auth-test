package provider

import (
	"fmt"
)

type Config interface {
	LoadConfig() error
	GetName() string
	GetIssuerURL() string
}

func NewConfig(providerType, name string) (Config, error) {
	switch providerType {
	case "AZURE_AUTH_PROVIDER":
		return &AzureConfig{Name: name}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
