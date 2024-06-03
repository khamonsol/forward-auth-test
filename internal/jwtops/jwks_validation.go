package jwtops

import (
	"context"
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/provider"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
)

type JWKSValidator struct {
	*jwk.AutoRefresh
	jwksURI      string
	validAlgVals []string
}

var validators map[string]*JWKSValidator

func init() {
	validators = make(map[string]*JWKSValidator)
}

func LoadValidator(p provider.Config) (*JWKSValidator, error) {
	validator := validators[p.GetName()]
	if validator == nil {

		jwkValidator := &JWKSValidator{
			AutoRefresh: jwk.NewAutoRefresh(context.Background()),
		}
		oidcCfg, err := provider.LoadOIDCConfig(p.GetIssuerUrl())
		if err != nil {
			return nil, err
		}
		jwkValidator.jwksURI = oidcCfg.JwksURI
		jwkValidator.validAlgVals = oidcCfg.IDTokenSigningAlgValuesSupported
		jwkValidator.Configure(oidcCfg.JwksURI)
		validators[p.GetName()] = jwkValidator
	}
	return validator, nil
}

// ValidateTokenSignature validates the token's signature using JWKS.
func (j *JWKSValidator) parseAndVerifyTokenString(tokenString string) (*jwt.Token, error) {
	keyfunc, err := j.getKeyfunc()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, keyfunc, jwt.WithValidMethods(j.validAlgVals))

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Ensure the token is valid and has the correct claims type
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (j *JWKSValidator) getKeyfunc() (jwt.Keyfunc, error) {
	jwks, err := j.AutoRefresh.Fetch(context.Background(), j.jwksURI)
	if err != nil {
		return nil, err
	}
	return func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			if key, found := jwks.LookupKeyID(kid); found {
				var rawKey interface{}
				if err := key.Raw(&rawKey); err != nil {
					return nil, fmt.Errorf("failed to get raw key: %w", err)
				}
				return rawKey, nil
			}
		}
		return nil, fmt.Errorf("unable to find key")
	}, nil
}
