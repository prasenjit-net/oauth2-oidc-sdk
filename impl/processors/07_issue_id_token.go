package processors

import (
	"context"
	sdk "oauth2-oidc-sdk"
	"oauth2-oidc-sdk/impl/sdkerror"
	"time"
)

type DefaultIDTokenIssuer struct {
	IDTokenStrategy sdk.IIDTokenStrategy
	Lifespan        time.Duration
}

func (d *DefaultIDTokenIssuer) HandleAuthEP(_ context.Context, requestContext sdk.IAuthenticationRequestContext) (sdk.IError, sdk.Result) {
	if requestContext.GetResponseType().Has("id_token") {
		expiry := requestContext.GetRequestedAt().UTC().Add(d.Lifespan).Round(time.Second)
		profile := requestContext.GetProfile()
		client := requestContext.GetClient()
		var tClaims map[string]interface{}
		token, err := d.IDTokenStrategy.GenerateIDToken(profile, client, expiry, tClaims)
		if err != nil {
			return sdkerror.ErrInvalidGrant, sdk.ResultNoOperation //todo change
		}
		requestContext.IssueIDToken(token)
	}
	return nil, sdk.ResultNoOperation
}

func (d *DefaultIDTokenIssuer) HandleTokenEP(_ context.Context, requestContext sdk.ITokenRequestContext) sdk.IError {
	if requestContext.GetGrantedScopes().Has("openid") {
		expiry := requestContext.GetRequestedAt().UTC().Add(d.Lifespan).Round(time.Second)
		profile := requestContext.GetProfile()
		client := requestContext.GetClient()
		var tClaims map[string]interface{}
		token, err := d.IDTokenStrategy.GenerateIDToken(profile, client, expiry, tClaims)
		if err != nil {
			return sdkerror.ErrInvalidGrant //todo change
		}
		requestContext.IssueIDToken(token)
	}
	return nil
}

func (d *DefaultIDTokenIssuer) Configure(strategy interface{}, config *sdk.Config, _ ...interface{}) {
	d.IDTokenStrategy = strategy.(sdk.IIDTokenStrategy)
	d.Lifespan = config.AccessTokenLifespan
}
