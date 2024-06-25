package provider

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
)

type AzureProviderConfig struct {
	ClientID   string `yaml:"client_id"`
	TenantID   string `yaml:"tenant_id"`
	SigningKey string `yaml:"signing_key"`
	IssuerURL  string
	JWKSURL    string
	AuthURL    string
	TokenURL   string
	Name       string
}

func (c *AzureProviderConfig) LoadProviderConfig(api util.KubernetesClient) error {
	configMap, err := api.GetConfigMap(, c.Name)
	if err != nil {
		return fmt.Errorf("failed to load config map: %v", err)
	}

	c.ClientID = configMap.Data["client_id"]
	c.TenantID = configMap.Data["tenant_id"]
	c.IssuerURL = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", c.TenantID)
	c.JWKSURL = fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", c.TenantID)
	c.AuthURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", c.TenantID)
	c.TokenURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", c.TenantID)

	return nil
}

func (c *AzureProviderConfig) GetName() string {
	return c.Name
}

func (c *AzureProviderConfig) GetIssuerUrl() string {
	return c.IssuerURL
}
