package util

import (
	"fmt"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	k8testing "k8s.io/client-go/testing"
)

// MockNamespaceResolver is a mock implementation of the NamespaceResolver interface.
type MockNamespaceResolver struct {
	mock.Mock
}

func (m *MockNamespaceResolver) GetCurrentNamespace() (*string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*string), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestServerResolver_GetCurrentNamespace_Success(t *testing.T) {
	mockNamespace := "test-namespace"
	fs := afero.NewMemMapFs()
	resolver := NewServerResolver(&fs)

	// Mock the file reading operation
	err := afero.WriteFile(fs, "/var/run/secrets/kubernetes.io/serviceaccount/namespace", []byte(mockNamespace), 0644)
	if err != nil {
		t.FailNow()
	}

	namespace, err := resolver.GetCurrentNamespace()
	assert.NoError(t, err)
	assert.NotNil(t, namespace)
	assert.Equal(t, mockNamespace, *namespace)
}

func TestServerResolver_GetCurrentNamespace_Error(t *testing.T) {
	resolver := NewServerResolver(nil)
	namespace, err := resolver.GetCurrentNamespace()

	assert.Error(t, err)
	assert.Nil(t, namespace)
}

func TestNewKubeAPI_Success(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	api, err := newApi(clientset, nil)
	assert.NoError(t, err)
	assert.NotNil(t, api)
	assert.Equal(t, &ServerNamespaceResolver{}, api.NamespaceResolver)
}

func TestNewKubeAPI_Error(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	api, err := newApi(clientset, nil)
	assert.NoError(t, err)
	assert.NotNil(t, api)
}

func TestKubeAPI_GetConfigMap_Success(t *testing.T) {
	mockNamespace := "test-namespace"
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: mockNamespace,
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	mockResolver := new(MockNamespaceResolver)
	mockResolver.On("GetCurrentNamespace").Return(&mockNamespace, nil)

	api := &KubeAPI{
		ClientSet:         clientset,
		NamespaceResolver: mockResolver,
	}

	configMapData, err := api.GetConfigMap("test-configmap")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, configMapData)
	mockResolver.AssertExpectations(t)
}

func TestKubeAPI_GetConfigMap_NotFound(t *testing.T) {
	mockNamespace := "test-namespace"
	clientset := fake.NewSimpleClientset()

	mockResolver := new(MockNamespaceResolver)
	mockResolver.On("GetCurrentNamespace").Return(&mockNamespace, nil)

	api := &KubeAPI{
		ClientSet:         clientset,
		NamespaceResolver: mockResolver,
	}

	configMapData, err := api.GetConfigMap("non-existent-configmap")
	assert.Error(t, err)
	assert.Nil(t, configMapData)
	mockResolver.AssertExpectations(t)
}

func TestKubeAPI_GetConfigMap_ClientError(t *testing.T) {
	mockNamespace := "test-namespace"
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("get", "configmaps", func(action k8testing.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("mock client error")
	})

	mockResolver := new(MockNamespaceResolver)
	mockResolver.On("GetCurrentNamespace").Return(&mockNamespace, nil)

	api := &KubeAPI{
		ClientSet:         clientset,
		NamespaceResolver: mockResolver,
	}

	configMapData, err := api.GetConfigMap("test-configmap")
	assert.Error(t, err)
	assert.Nil(t, configMapData)
	mockResolver.AssertExpectations(t)
}
