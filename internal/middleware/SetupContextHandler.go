package middleware

import (
	"context"
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"log/slog"
	"net/http"
)

const kubeApiKey = "KUBE_API"
const policyKey = "REQUEST_POLICY"

func SetupContextHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrId := GetAuthRequestTxId(r)

		ctx, kApi, err := addKubeApiToContext(corrId, w, r.Context())
		if err != nil {
			return
		}
		ctx, err = addPolicyToContext(corrId, kApi, r.Host, r.Method, r.URL.Path, w, r.Context())
		if err != nil {
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addKubeApiToContext(corrId string, w http.ResponseWriter, c context.Context) (context.Context, *util.KubeAPI, error) {
	kApi, err := util.NewKubeAPI(nil)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusInternalServerError, corrId)
		return nil, nil, err
	}
	return context.WithValue(c, kubeApiKey, kApi), kApi, nil
}

func addPolicyToContext(corrId string, k *util.KubeAPI, host string, method string, path string,
	w http.ResponseWriter, c context.Context) (context.Context, error) {
	pols, err := policy.LoadPolicies(host, k)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusInternalServerError, corrId)
		return c, err
	}
	err = pols.GetPolicy(path, method)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusInternalServerError, corrId)
		return c, err
	}
	return context.WithValue(c, policyKey, pols), nil
}

func GetKubeApi(r *http.Request) *util.KubeAPI {
	ctx := r.Context()
	val, ok := ctx.Value(kubeApiKey).(*util.KubeAPI)
	if !ok {
		slog.Info("Unable to get Kube API from context")
	}
	return val
}

func GetPolicy(r *http.Request) *policy.Policy {
	ctx := r.Context()
	val, ok := ctx.Value(policyKey).(*policy.Policy)
	if !ok {
		slog.Info("Unable to get policy from context")
	}
	return val
}
