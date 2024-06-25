package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewKubeAPI_Success(t *testing.T) {
	// Create a fake Kubernetes client
	clientset := fake.NewSimpleClientset()

	// Inject the fake client into the KubeAPI struct
	kubeAPI := &KubeAPI{Clientset: clientset}

	assert.NotNil(t, kubeAPI)
	assert.NotNil(t, kubeAPI.Clientset)
}

func TestLoadKubeConfig_Success(t *testing.T) {
	// This test assumes that the environment is set up correctly for in-cluster config
	// or that the kubeconfig file is available at the default location.
	config, err := LoadKubeConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
}
