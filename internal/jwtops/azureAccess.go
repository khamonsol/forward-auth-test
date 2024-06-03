package jwtops

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/SoleaEnergy/forwardAuth/internal/provider"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

// AzureMSALAccessToken represents a Microsoft Azure MSAL access token.
type AzureMSALAccessToken struct {
	Aud               string   `json:"aud"`
	Iss               string   `json:"iss"`
	Iat               int64    `json:"iat"`
	Nbf               int64    `json:"nbf"`
	Exp               int64    `json:"exp"`
	Aio               string   `json:"aio"`
	Groups            []string `json:"groups"`
	Name              string   `json:"name"`
	Nonce             string   `json:"nonce"`
	Oid               string   `json:"oid"`
	PreferredUsername string   `json:"preferred_username"`
	Rh                string   `json:"rh"`
	Roles             []string `json:"roles"`
	Sub               string   `json:"sub"`
	Tid               string   `json:"tid"`
	Uti               string   `json:"uti"`
	Ver               string   `json:"ver"`
	Alg               string   `json:"alg"`
	tokenString       string
}

// LoadToken Function performs the following:
// - Inspect http.Request for header configured on server.Config
// - Strip any prefixes on the token like 'Bearer', 'bearer' etc.
// - Serialize the JWT using JWKS validation, if JWKS validation fails,an error will be returned
// - Validate that the Token is allowed to access the destination endpoint based on policy
// - Return AzureMSALAccessToken if the token is valid or an error if it fails at any step.
func LoadToken(r *http.Request, p provider.Config) (*AzureMSALAccessToken, error) {
	var tokenString string
	for _, val := range SOLEA_AUTH_HEADERS {
		if tokenString == "" {
			tokenString = r.Header.Get(val)
		}
	}
	if tokenString == "" {
		return nil, fmt.Errorf("missing access token")
	}
	if strings.Contains(tokenString, "Bearer") {
		tokenString = strings.Replace(tokenString, "Bearer", "", 1)
	}
	if strings.Contains(tokenString, "bearer") {
		tokenString = strings.Replace(tokenString, "bearer", "", 1)
	}
	tokenString = strings.TrimSpace(tokenString)
	v, err := LoadValidator(p)
	if err != nil {
		return nil, err
	}
	token, parseError := v.parseAndVerifyTokenString(tokenString)
	if parseError != nil {
		return nil, parseError
	}
	resp, err := newAzureMSALAccessToken(token)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func newAzureMSALAccessToken(token *jwt.Token) (*AzureMSALAccessToken, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to parse claims")
	}

	alg, ok := token.Header["alg"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse alg from token header")
	}

	return &AzureMSALAccessToken{
		Aud:               util.GetStringClaim(claims, "aud"),
		Iss:               util.GetStringClaim(claims, "iss"),
		Iat:               util.GetInt64Claim(claims, "iat"),
		Nbf:               util.GetInt64Claim(claims, "nbf"),
		Exp:               util.GetInt64Claim(claims, "exp"),
		Aio:               util.GetStringClaim(claims, "aio"),
		Groups:            util.GetStringSliceClaim(claims, "groups"),
		Name:              util.GetStringClaim(claims, "name"),
		Nonce:             util.GetStringClaim(claims, "nonce"),
		Oid:               util.GetStringClaim(claims, "oid"),
		PreferredUsername: util.GetStringClaim(claims, "preferred_username"),
		Rh:                util.GetStringClaim(claims, "rh"),
		Roles:             util.GetStringSliceClaim(claims, "roles"),
		Sub:               util.GetStringClaim(claims, "sub"),
		Tid:               util.GetStringClaim(claims, "tid"),
		Uti:               util.GetStringClaim(claims, "uti"),
		Ver:               util.GetStringClaim(claims, "ver"),
		Alg:               alg,
	}, nil
}

func (t *AzureMSALAccessToken) ValidateTokenForPolicy(p policy.Policy) error {
	matchingRoles := util.Intersection(t.Roles, p.Roles)
	matchingUsers := util.Intersection([]string{t.GetShortName()}, p.Users)
	numMatching := len(matchingRoles) + len(matchingUsers)
	if numMatching <= 0 {
		return fmt.Errorf("no users or roles matching")
	}
	return nil
}

func (t *AzureMSALAccessToken) GetShortName() string {
	return strings.Split(t.PreferredUsername, "@")[0]
}
