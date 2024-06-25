package policy

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

const namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

type Policy struct {
	ProviderName string   `yaml:"provider_name"`
	ProviderType string   `yaml:"provider_type"`
	Roles        []string `yaml:"roles"` // Roles OIDC roles.
	Users        []string `yaml:"users"` // Users are specific authenticated users.
	configMap    map[string]string
}

func getNamespace() (string, error) {
	data, err := os.ReadFile(namespaceFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func LoadPolicies(host string, api util.KubernetesClient) (*Policy, error) {
	namespace, err := getNamespace()
	if err != nil {
		return nil, err
	}

	snakeCaseHost := strings.ToLower(strings.ReplaceAll(host, ".", "_"))
	name := fmt.Sprintf("access_policy_%s", snakeCaseHost)

	cm, err := api.GetConfigMap(namespace, name)
	if err != nil {
		return nil, err
	}

	p := Policy{
		configMap: cm.Data,
	}
	return &p, nil
}

func (p *Policy) GetPolicy(path string, method string) error {
	key := strings.ToLower(fmt.Sprintf("%s_%s", strings.ReplaceAll(strings.Trim(path, "/"), "/", "_"), method))
	yamlData, exists := p.configMap[key]
	if !exists {
		return fmt.Errorf("policy not found for key %s", key)
	}
	err := yaml.Unmarshal([]byte(yamlData), p)
	if err != nil {
		return fmt.Errorf("error parsing policy data for key %s: %w", key, err)
	}
	return nil
}
