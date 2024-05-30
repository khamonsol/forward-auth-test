package policy

import (
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
