package util

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeAPI struct {
	ClientSet kubernetes.Interface
	Namespace string
}

func NewKubeAPI(namespace string) (*KubeAPI, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	return &KubeAPI{
		ClientSet: clientset,
		Namespace: namespace,
	}, nil
}

func (api *KubeAPI) LoadConfig(name string) (map[string]string, error) {
	configMap, err := api.ClientSet.CoreV1().ConfigMaps(api.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap: %w", err)
	}
	return configMap.Data, nil
}
