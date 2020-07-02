package processors

import (
	"context"
	sdk "oauth2-oidc-sdk"
	"time"
)

type DefaultAccessTokenIssuer struct {
	AccessTokenStrategy sdk.IAccessTokenStrategy
	Lifespan            time.Duration
}

func (d *DefaultAccessTokenIssuer) HandleAuthEP(_ context.Context, requestContext sdk.IAuthenticationRequestContext) sdk.IError {
	if requestContext.GetResponseType().Has(sdk.ResponseTypeToken) {
		token, signature := d.AccessTokenStrategy.GenerateAccessToken()
		expiry := requestContext.GetRequestedAt().UTC().Add(d.Lifespan).Round(time.Second)
		requestContext.IssueAccessToken(token, signature, expiry)
	}
	return nil
}

func (d *DefaultAccessTokenIssuer) HandleTokenEP(_ context.Context, requestContext sdk.ITokenRequestContext) sdk.IError {
	token, signature := d.AccessTokenStrategy.GenerateAccessToken()
	expiry := requestContext.GetRequestedAt().UTC().Add(d.Lifespan).Round(time.Second)
	requestContext.IssueAccessToken(token, signature, expiry)
	return nil
}

func (d *DefaultAccessTokenIssuer) Configure(strategy interface{}, config *sdk.Config, _ ...interface{}) {
	d.AccessTokenStrategy = strategy.(sdk.IAccessTokenStrategy)
	d.Lifespan = config.AccessTokenLifespan
}
