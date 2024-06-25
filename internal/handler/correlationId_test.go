package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CorrelationIdTestSuite struct {
	suite.Suite
}

func (suite *CorrelationIdTestSuite) SetupTest() {
	// Setup code, if needed
}

func (suite *CorrelationIdTestSuite) TestCorrelationId() {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rr := httptest.NewRecorder()

	handler := CorrelationId(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrId, ok := GetCorrelationId(r)
		suite.True(ok)
		suite.NotEmpty(corrId)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
}

func TestCorrelationIdTestSuite(t *testing.T) {
	suite.Run(t, new(CorrelationIdTestSuite))
}
