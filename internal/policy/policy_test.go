package policy

import (
	"testing"

	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// MockKubeAPI is a mock implementation of the KubeAPI interface.
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

func TestLoadPolicy_Success(t *testing.T) {
	mockNamespace := "test-namespace"
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "access_policy_test_host_com",
			Namespace: mockNamespace,
		},
		Data: map[string]string{
			"api_data_get": "audience: test-audience\n issuer: test-issuer\n roles:\n  - role1\n  - role2\nusers:\n  - user1\n  - user2",
		},
	})
	mockResolver := new(MockNamespaceResolver)
	mockResolver.On("GetCurrentNamespace").Return(&mockNamespace, nil)
	api := &util.KubeAPI{
		ClientSet:         clientset,
		NamespaceResolver: mockResolver,
	}

	p, err := LoadPolicies("test.host.com", *api)
	assert.NoError(t, err)
	assert.NotNil(t, p)
}

func TestLoadPolicy_Error(t *testing.T) {
	mockNamespace := "test-namespace"
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "access_policy_test_host1_com",
			Namespace: mockNamespace,
		},
		Data: map[string]string{
			"api_data_get": "audience: test-audience\nissuer: test-issuer\n roles:\n  - role1\n  - role2\nusers:\n  - user1\n  - user2",
		},
	})
	mockResolver := new(MockNamespaceResolver)
	mockResolver.On("GetCurrentNamespace").Return(&mockNamespace, nil)
	api := &util.KubeAPI{
		ClientSet:         clientset,
		NamespaceResolver: mockResolver,
	}

	p, err := LoadPolicies("test.host.com", *api)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestGetPolicy_Success(t *testing.T) {
	mockNamespace := "test-namespace"
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "access_policy_test_host_com",
			Namespace: mockNamespace,
		},
		Data: map[string]string{
			"api_data_get": "audience: test-audience\nissuer: test-issuer\nroles:\n  - role1\n  - role2\nusers:\n  - user1\n  - user2",
		},
	})
	mockResolver := new(MockNamespaceResolver)
	mockResolver.On("GetCurrentNamespace").Return(&mockNamespace, nil)
	api := &util.KubeAPI{
		ClientSet:         clientset,
		NamespaceResolver: mockResolver,
	}

	p, err := LoadPolicies("test.host.com", *api)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	err = p.GetPolicy("/api/data/", "GET")
	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"role1", "role2"}, p.Roles)
	assert.ElementsMatch(t, []string{"user1", "user2"}, p.Users)
	assert.Equal(t, "test-audience", p.Audience)
	assert.Equal(t, "test-issuer", p.Issuer)

}
