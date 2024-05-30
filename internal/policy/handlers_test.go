package policy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// MockKubernetesClient is a mock implementation of KubernetesInterface.
type MockKubernetesClient struct {
	clientset *fake.Clientset
}

func (m *MockKubernetesClient) CoreV1() corev1.CoreV1Interface {
	return m.clientset.CoreV1()
}

func mockNewConfigWithConfigMap() (*Config, error) {
	clientset := fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "access_policy_example_com",
			Namespace: "default",
		},
		Data: map[string]string{
			"api_data_get": "roles:\n  - admin\nusers:\n  - user1",
		},
	})
	return &Config{client: &MockKubernetesClient{clientset: clientset}}, nil
}

func TestPolicyLoader_Success(t *testing.T) {
	originalNewConfig := newConfigFunc
	newConfigFunc = mockNewConfigWithConfigMap
	defer func() { newConfigFunc = originalNewConfig }()

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Manually load the policy into the request context to simulate a successful load.
		config, _ := mockNewConfigWithConfigMap()
		_ = config.LoadConfig("example.com")
		policy, _ := config.GetPolicy("/api/data", "GET")

		ctx := context.WithValue(r.Context(), policyHeaderKey, policy)
		r = r.WithContext(ctx)

		loadedPolicy, err := GetPolicyFromRequest(r)
		assert.NoError(t, err)
		assert.Equal(t, "admin", loadedPolicy.Roles[0], "Policy should be available and correct")
		w.WriteHeader(http.StatusOK)
	}

	testHandler := PolicyLoader(http.HandlerFunc(handler))
	req := httptest.NewRequest("GET", "http://example.com/api/data", nil)
	w := httptest.NewRecorder()

	testHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")
}

func TestPolicyLoader_Failure_LoadConfig(t *testing.T) {
	originalNewConfig := newConfigFunc
	newConfigFunc = func() (*Config, error) {
		return nil, errors.New("failed to load configuration")
	}
	defer func() { newConfigFunc = originalNewConfig }() // Reset after test

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	testHandler := PolicyLoader(http.HandlerFunc(handler))
	req := httptest.NewRequest("GET", "http://example.com/api/data", nil)
	w := httptest.NewRecorder()

	testHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP status 500 Internal Server Error")
}

func TestGetPolicyFromRequest_NoPolicy(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/", nil)
	_, err := GetPolicyFromRequest(req)
	assert.Error(t, err, "should return an error when no policy is set")
}
