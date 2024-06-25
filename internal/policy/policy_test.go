package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type MockKubernetesClient struct {
	clientset *fake.Clientset
}

func NewMockKubernetesClient() *MockKubernetesClient {
	return &MockKubernetesClient{
		clientset: fake.NewSimpleClientset(),
	}
}

func (m *MockKubernetesClient) GetConfigMap(namespace, name string) (*v1.ConfigMap, error) {
	return m.clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func TestLoadPolicies(t *testing.T) {
	mockClient := NewMockKubernetesClient()
	namespace := "default"
	host := "beyond.soleaenergy.com"

	configMapName := "access_policy_beyond_soleaenergy_com"
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"api_v1_pos_pnl_detail_post": `
provider_name: azure_beyond_prod
provider_type: azure
roles:
  - PNL_API_RW
users: []
`,
		},
	}
	_, err := mockClient.clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	assert.NoError(t, err)

	policy, err := LoadPolicies(host, mockClient, namespace)
	assert.NoError(t, err)
	assert.NotNil(t, policy)
}

func TestGetPolicy(t *testing.T) {
	mockClient := NewMockKubernetesClient()
	namespace := "default"
	host := "beyond.soleaenergy.com"

	configMapName := "access_policy_beyond_soleaenergy_com"
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"api_v1_pos_pnl_detail_post": `
provider_name: azure_beyond_prod
provider_type: azure
roles:
  - PNL_API_RW
users: []
`,
		},
	}
	_, err := mockClient.clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	assert.NoError(t, err)

	policy, err := LoadPolicies(host, mockClient, namespace)
	assert.NoError(t, err)
	assert.NotNil(t, policy)

	err = policy.GetPolicy("/api/v1/pos_pnl_detail", "POST")
	assert.NoError(t, err)
	assert.Equal(t, "azure_beyond_prod", policy.ProviderName)
	assert.Equal(t, "azure", policy.ProviderType)
	assert.Contains(t, policy.Roles, "PNL_API_RW")
}
