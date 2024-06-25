package provider

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

func TestAzureProviderConfig_LoadProviderConfig(t *testing.T) {
	mockClient := NewMockKubernetesClient()
	namespace := "default"
	name := "azure_beyond_prod"

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"client_id": "your-client-id",
			"tenant_id": "your-tenant-id",
		},
	}
	_, err := mockClient.clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	assert.NoError(t, err)

	providerConfig := NewConfig("azure", name, namespace)
	assert.NotNil(t, providerConfig)

	err = providerConfig.LoadProviderConfig(mockClient)
	assert.NoError(t, err)
	assert.Equal(t, "your-client-id", providerConfig.(*AzureProviderConfig).ClientID)
	assert.Equal(t, "your-tenant-id", providerConfig.(*AzureProviderConfig).TenantID)
	assert.Equal(t, "https://login.microsoftonline.com/your-tenant-id/v2.0", providerConfig.GetIssuerUrl())
}

func TestAzureProviderConfig_GetName(t *testing.T) {
	providerConfig := NewConfig("azure", "azure_beyond_prod", "default")
	assert.Equal(t, "azure_beyond_prod", providerConfig.GetName())
}

func TestAzureProviderConfig_GetIssuerUrl(t *testing.T) {
	providerConfig := NewConfig("azure", "azure_beyond_prod", "default")
	providerConfig.(*AzureProviderConfig).TenantID = "your-tenant-id"
	providerConfig.(*AzureProviderConfig).IssuerURL = "https://login.microsoftonline.com/your-tenant-id/v2.0"
	assert.Equal(t, "https://login.microsoftonline.com/your-tenant-id/v2.0", providerConfig.GetIssuerUrl())
}
