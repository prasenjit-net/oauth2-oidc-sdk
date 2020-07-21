package manager

import (
	sdk "github.com/identityOrg/oidcsdk"
	"gopkg.in/square/go-jose.v2/json"
	"log"
	"net/http"
	"net/url"
	"path"
)

func (d *DefaultManager) ProcessDiscoveryEP(writer http.ResponseWriter, _ *http.Request) {
	issuerUrl, err := url.Parse(d.Config.Issuer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	metadata := DiscoveryMetadata{
		Issuer:                            issuerUrl.String(),
		AuthorizationEndpoint:             combinePath(*issuerUrl, "authorize"),
		TokenEndpoint:                     combinePath(*issuerUrl, "token"),
		IntrospectionEndpoint:             combinePath(*issuerUrl, "introspect"),
		RevocationEndpoint:                combinePath(*issuerUrl, "revoke"),
		UserInfoEndpoint:                  combinePath(*issuerUrl, "me"),
		JwksUri:                           combinePath(*issuerUrl, "keys"),
		SubjectTypesSupported:             []string{"public"},
		GrantTypesSupported:               []string{"authorization_code", "password", "refresh_token", "client_credentials", "implicit"},
		ResponseModesSupported:            []string{"query", "fragment"},
		ResponseTypesSupported:            []string{"code", "token", "id_token", "code id_token", "code token", "token id_token", "code token id_token"},
		TokenEndpointAuthMethodsSupported: []string{"basic", "forms"},
		IdTokenSigningAlgValuesSupported:  []string{"RS256"},
	}

	writer.Header().Add(sdk.HeaderContentType, sdk.ContentTypeJson)
	writer.WriteHeader(http.StatusOK)
	log.Print(json.NewEncoder(writer).Encode(metadata))
}

func combinePath(issuerUrl url.URL, appendPath string) string {
	issuerUrl.Path = path.Join(issuerUrl.Path, appendPath)
	return issuerUrl.String()
}

type DiscoveryMetadata struct {
	Issuer                            string   `json:"issuer,omitempty"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint,omitempty"`
	TokenEndpoint                     string   `json:"token_endpoint,omitempty"`
	IntrospectionEndpoint             string   `json:"introspection_endpoint,omitempty"`
	RevocationEndpoint                string   `json:"revocation_endpoint,omitempty"`
	UserInfoEndpoint                  string   `json:"user_info_endpoint,omitempty"`
	JwksUri                           string   `json:"jwks_uri,omitempty"`
	ScopesSupported                   []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported            []string `json:"response_types_supported,omitempty"`
	ResponseModesSupported            []string `json:"response_modes_supported,omitempty"`
	GrantTypesSupported               []string `json:"grant_types_supported,omitempty"`
	SubjectTypesSupported             []string `json:"subject_types_supported,omitempty"`
	IdTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported,omitempty"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	ClaimsSupported                   []string `json:"claims_supported,omitempty"`
}
