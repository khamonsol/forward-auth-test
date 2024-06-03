package policy

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"gopkg.in/yaml.v3"
	"log/slog"
	"strings"
)

type Policy struct {
	ProviderName string   `yaml:"provider_name"`
	ProviderType string   `yaml:"provider_type"`
	Roles        []string `yaml:"roles"` // Roles OIDC roles.
	Users        []string `yaml:"users"` // Users are specific authenticated users.
	configMap    map[string]string
}

func LoadPolicies(host string, api *util.KubeAPI) (*Policy, error) {
	snakeCaseHost := strings.ToLower(strings.ReplaceAll(host, ".", "_"))
	name := fmt.Sprintf("access_policy_%s", snakeCaseHost)

	cmap, err := api.GetConfigMap(name)
	if err != nil {
		return nil, err
	}

	p := Policy{
		configMap: cmap,
	}
	return &p, nil
}

func (p *Policy) GetPolicy(path string, method string) error {
	key := strings.ToLower(fmt.Sprintf("%s_%s", strings.ReplaceAll(strings.Trim(path, "/"), "/", "_"), method))
	slog.Info(fmt.Sprintf("Getting policy for %s", key))
	yamlData := p.configMap[key]
	err := yaml.Unmarshal([]byte(yamlData), p)

	if err != nil {
		return fmt.Errorf("error parsing policy data for key %s: %w", key, err)
	}
	return nil
}
