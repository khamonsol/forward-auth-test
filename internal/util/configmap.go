package util

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"time"
)

type NamespaceResolver interface {
	GetCurrentNamespace() (*string, error)
}

type ServerNamespaceResolver struct {
	FileSystem afero.Fs
}

func NewServerResolver(fs *afero.Fs) *ServerNamespaceResolver {
	var useFs afero.Fs
	if fs == nil {
		realFs := afero.NewOsFs()
		cacheFs := afero.NewMemMapFs()
		//The namespace can really be cached for all time as this will always be the same until the server restarts.
		useFs = afero.NewCacheOnReadFs(realFs, cacheFs, 99999*time.Hour)
	} else {
		useFs = *fs
	}
	return &ServerNamespaceResolver{
		FileSystem: useFs,
	}
}
func (r *ServerNamespaceResolver) GetCurrentNamespace() (*string, error) {
	const namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	ns, nsErr := afero.ReadFile(r.FileSystem, namespacePath)

	if nsErr != nil {
		return nil, fmt.Errorf("failed to read namespace: %w", nsErr)
	}
	namespace := string(ns)
	return &namespace, nil
}

type KubeAPI struct {
	ClientSet         kubernetes.Interface
	NamespaceResolver NamespaceResolver
}

func NewKubeAPI(fs *afero.Fs) (*KubeAPI, error) {
	config, cfgErr := rest.InClusterConfig()
	if cfgErr != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", cfgErr)
	}
	clientset, cliErr := kubernetes.NewForConfig(config)
	if cliErr != nil {
		return nil, fmt.Errorf("failed to get client: %w", cliErr)
	}
	return newApi(clientset, fs)
}

func newApi(clientSet kubernetes.Interface, fs *afero.Fs) (*KubeAPI, error) {
	return &KubeAPI{
		ClientSet:         clientSet,
		NamespaceResolver: NewServerResolver(fs),
	}, nil
}

func (api *KubeAPI) GetConfigMap(name string) (map[string]string, error) {
	cns, err := api.NamespaceResolver.GetCurrentNamespace()
	if err != nil {
		slog.Error("failed to get current namespace")
		return nil, err
	}

	configMap, err := api.ClientSet.CoreV1().ConfigMaps(*cns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap: %w", err)
	}
	return configMap.Data, nil
}
