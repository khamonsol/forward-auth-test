// Package policy provides structures and functions to manage and fetch access control
// policies from a Kubernetes ConfigMap. It uses the in-cluster configuration method
// to access the Kubernetes API, suitable for applications running within a Kubernetes cluster.
// https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
package policy

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesInterface defines the functions used from the Kubernetes client,
// allowing for easier testing and mocking.
type KubernetesInterface interface {
	CoreV1() v1.CoreV1Interface
}

// Policy defines a set of roles and users that have access to certain resources.
// It is structured to be stored and retrieved from a Kubernetes ConfigMap.
type Policy struct {
	Roles []string `yaml:"roles"` // Roles are Kubernetes RBAC roles applicable.
	Users []string `yaml:"users"` // Users are specific authenticated users.
}

// Config holds the internal state necessary to interact with the Kubernetes APIs
// and manage the access control policies.
type Config struct {
	Policies map[string]Policy   // Policies maps request paths and methods to access policies.
	client   KubernetesInterface // client is the interface to Kubernetes client, for real or testing use.
}

// NewConfig creates a new Config object using in-cluster configuration to communicate
// with the Kubernetes API. It returns an initialized Config object or an error if the
// in-cluster configuration setup fails.
// // NewConfig creates a new Config object using in-cluster configuration to communicate
// // with the Kubernetes API. It returns an initialized Config object or an error if the
// // in-cluster configuration setup fails.
func NewConfig() (*Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		slog.Error("Error getting kubernetes config:", err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error getting kubernetes client:", err)
		return nil, err
	}
	return &Config{
		client: clientset,
	}, nil
}

// getCurrentNamespace fetches the namespace that the current client is operating within.
// It attempts to load the client configuration from the default kubeconfig path,
// then retrieves the namespace specified for the current context.
// If no namespace is specified in the kubeconfig or an error occurs during the loading of configurations,
// it defaults to 'default'.
func getCurrentNamespace() string {
	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		slog.Error("Error loading kubeconfig: %v", err)
		return "default"
	}
	namespace := clientCfg.Contexts[clientCfg.CurrentContext].Namespace
	if namespace == "" {
		slog.Error("No namespace specified in the kubeconfig, defaulting to 'default'")
		namespace = "default"
	}
	return namespace
}

// LoadConfig loads access policies from a Kubernetes ConfigMap specified by the host.
// The host is converted to snake_case to match the naming convention of ConfigMaps.
// It returns an error if fetching the ConfigMap fails or the data is improperly formatted.
func (c *Config) LoadConfig(host string) error {
	snakeCaseHost := strings.ToLower(strings.ReplaceAll(host, ".", "_"))
	configMapName := fmt.Sprintf("access_policy_%s", snakeCaseHost)
	namespace := getCurrentNamespace()

	cm, err := c.client.CoreV1().ConfigMaps(namespace).Get(context.Background(), configMapName, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("Error getting configmap %s/%s: %v", namespace, configMapName, err)
		slog.Error(msg)
		return fmt.Errorf("failed to fetch config map: %w", err)
	}

	c.Policies = make(map[string]Policy)
	for key, yamlData := range cm.Data {
		var policy Policy
		if err := yaml.Unmarshal([]byte(yamlData), &policy); err != nil {
			return fmt.Errorf("error parsing policy data for key %s: %w", key, err)
		}
		c.Policies[key] = policy
	}
	return nil
}

// GetPolicy retrieves the policy associated with the specified path and method.
// It returns the policy and a boolean indicating if the policy was found.
func (c *Config) GetPolicy(path, method string) (*Policy, error) {
	key := fmt.Sprintf("%s_%s", strings.ReplaceAll(strings.ToLower(path), "/", "_"), strings.ToLower(method))
	policy, ok := c.Policies[key]
	if !ok {
		return nil, fmt.Errorf("policy %s not found", key)
	}
	return &policy, nil
}
