package util

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	k8testing "k8s.io/client-go/testing"
)

func TestNewKubeAPI_Success(t *testing.T) {
	namespace := "test-namespace"
	clientset := fake.NewSimpleClientset()
	api, err := NewKubeAPI(namespace)
	assert.NoError(t, err)
	assert.NotNil(t, api)
	assert.Equal(t, namespace, api.Namespace)
}

func TestKubeAPI_LoadConfig_Success(t *testing.T) {
	namespace := "test-namespace"
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: namespace,
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	api := &KubeAPI{
		ClientSet: clientset,
		Namespace: namespace,
	}

	configMapData, err := api.LoadConfig("test-configmap")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, configMapData)
}

func TestKubeAPI_LoadConfig_NotFound(t *testing.T) {
	namespace := "test-namespace"
	clientset := fake.NewSimpleClientset()

	api := &KubeAPI{
		ClientSet: clientset,
		Namespace: namespace,
	}

	configMapData, err := api.LoadConfig("non-existent-configmap")
	assert.Error(t, err)
	assert.Nil(t, configMapData)
}

func TestKubeAPI_LoadConfig_ClientError(t *testing.T) {
	namespace := "test-namespace"
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("get", "configmaps", func(action k8testing.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("mock client error")
	})

	api := &KubeAPI{
		ClientSet: clientset,
		Namespace: namespace,
	}

	configMapData, err := api.LoadConfig("test-configmap")
	assert.Error(t, err)
	assert.Nil(t, configMapData)
}
