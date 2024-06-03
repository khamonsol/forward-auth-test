package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OIDCConfig represents the structure of the OIDC well-known configuration response.
type OIDCConfig struct {
	Issuer                             string   `json:"issuer"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint"`
	TokenEndpoint                      string   `json:"token_endpoint"`
	UserinfoEndpoint                   string   `json:"userinfo_endpoint"`
	JwksURI                            string   `json:"jwks_uri"`
	RegistrationEndpoint               string   `json:"registration_endpoint"`
	ScopesSupported                    []string `json:"scopes_supported"`
	ResponseTypesSupported             []string `json:"response_types_supported"`
	ResponseModesSupported             []string `json:"response_modes_supported"`
	GrantTypesSupported                []string `json:"grant_types_supported"`
	SubjectTypesSupported              []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported   []string `json:"id_token_signing_alg_values_supported"`
	ClaimsSupported                    []string `json:"claims_supported"`
	RequestURIParameterSupported       bool     `json:"request_uri_parameter_supported"`
	RequireRequestURIRegistration      bool     `json:"require_request_uri_registration"`
	OpPolicyURI                        string   `json:"op_policy_uri"`
	OpTosURI                           string   `json:"op_tos_uri"`
	CheckSessionIframe                 string   `json:"check_session_iframe"`
	EndSessionEndpoint                 string   `json:"end_session_endpoint"`
	RevocationEndpoint                 string   `json:"revocation_endpoint"`
	IntrospectionEndpoint              string   `json:"introspection_endpoint"`
	CodeChallengeMethodsSupported      []string `json:"code_challenge_methods_supported"`
	TokenEndpointAuthMethodsSupported  []string `json:"token_endpoint_auth_methods_supported"`
	TokenEndpointAuthSigningAlgValues  []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	DisplayValuesSupported             []string `json:"display_values_supported"`
	ClaimTypesSupported                []string `json:"claim_types_supported"`
	ClaimsLocalesSupported             []string `json:"claims_locales_supported"`
	UILocalesSupported                 []string `json:"ui_locales_supported"`
	ClaimsParameterSupported           bool     `json:"claims_parameter_supported"`
	RequestParameterSupported          bool     `json:"request_parameter_supported"`
	RequestObjectSigningAlgValues      []string `json:"request_object_signing_alg_values_supported"`
	BackchannelLogoutSupported         bool     `json:"backchannel_logout_supported"`
	BackchannelLogoutSessionSupported  bool     `json:"backchannel_logout_session_supported"`
	FrontchannelLogoutSupported        bool     `json:"frontchannel_logout_supported"`
	FrontchannelLogoutSessionSupported bool     `json:"frontchannel_logout_session_supported"`
}

// LoadOIDCConfig fetches the OIDC well-known configuration from the provided issuer URL.
//
// Parameters:
//   - issuerURL: The base URL of the OIDC issuer (e.g., "https://accounts.google.com").
//
// Returns:
//   - An OIDCConfig struct containing the OIDC configuration details.
//   - An error if there was an issue fetching or parsing the configuration.
func LoadOIDCConfig(issuerURL string) (*OIDCConfig, error) {
	wellKnownURL := fmt.Sprintf("%s/.well-known/openid-configuration", issuerURL)

	resp, err := http.Get(wellKnownURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OIDC configuration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch OIDC configuration: received status code %d", resp.StatusCode)
	}

	var config OIDCConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode OIDC configuration: %w", err)
	}

	return &config, nil
}
