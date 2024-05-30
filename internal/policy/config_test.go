package policy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLoadConfig(t *testing.T) {
	// Creating a fake client set with a predefined ConfigMap.
	clientset := fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "access_policy_api_soleaenergy_com",
			Namespace: "default",
		},
		Data: map[string]string{
			"api_data_get": "roles:\n  - admin\nusers:\n  - user1",
		},
	})

	// Inject the fake client into the configuration.
	config := &Config{
		client: clientset,
	}

	// Attempt to load the configuration.
	err := config.LoadConfig("api.soleaenergy.com")
	assert.NoError(t, err, "LoadConfig should not return an error")
	assert.NotEmpty(t, config.Policies, "Policies should not be empty after loading")
	assert.Len(t, config.Policies["api_data_get"].Roles, 1, "Should have one role.")
	assert.Len(t, config.Policies["api_data_get"].Users, 1, "Should have one user.")
}

func TestGetPolicy(t *testing.T) {
	// Set up the configuration with predefined policies.
	config := &Config{
		Policies: map[string]Policy{
			"api_data_get": Policy{
				Roles: []string{"admin"},
				Users: []string{"user1"},
			},
		},
	}

	// Retrieve a specific policy.
	policy, err := config.GetPolicy("api/data", "GET")
	assert.NoError(t, err, "GetPolicy should not return an error for existing policies")
	assert.NotNil(t, policy, "Policy should not be nil for existing keys")
	assert.Contains(t, policy.Roles, "admin", "Policy should contain role 'admin'")
}

func TestNewConfig_Success(t *testing.T) {
	// Mock the function within the test
	originalNewConfig := newConfigFunc
	newConfigFunc = func() (*Config, error) {
		// Create a fake clientset that will simulate the successful creation of a Kubernetes client
		clientset := fake.NewSimpleClientset()
		return &Config{client: clientset}, nil
	}
	defer func() { newConfigFunc = originalNewConfig }() // Restore the original function after the test

	// Call the function that has been overridden for testing
	config, err := NewConfig()

	// Assertions to ensure it behaves as expected
	assert.NoError(t, err, "NewConfig should not return an error on success")
	assert.NotNil(t, config, "Config should not be nil on successful creation")
	assert.Implements(t, (*KubernetesInterface)(nil), config.client, "The client should implement KubernetesInterface")
}

func TestNewConfig_Failure(t *testing.T) {
	// Mock the function within the test to simulate a failure in acquiring the Kubernetes configuration
	originalNewConfig := newConfigFunc
	newConfigFunc = func() (*Config, error) {
		return nil, errors.New("failed to obtain Kubernetes config")
	}
	defer func() { newConfigFunc = originalNewConfig }() // Restore the original function after the test

	// Call the function that has been overridden for testing
	config, err := NewConfig()

	// Assertions to check the handling of failures
	assert.Error(t, err, "NewConfig should return an error on failure")
	assert.Nil(t, config, "Config should be nil when there's a failure in configuration")
}
