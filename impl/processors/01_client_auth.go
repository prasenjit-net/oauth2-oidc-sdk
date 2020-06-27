package processors

import (
	"context"
	sdk "oauth2-oidc-sdk"
	"oauth2-oidc-sdk/impl/sdkerror"
)

type DefaultClientAuthenticationProcessor struct {
	ClientStore sdk.IClientStore
}

func (d *DefaultClientAuthenticationProcessor) Configure(_ interface{}, _ *sdk.Config, args ...interface{}) {
	for _, arg := range args {
		if cs, ok := arg.(sdk.IClientStore); ok {
			d.ClientStore = cs
			break
		}
	}
	if d.ClientStore == nil {
		panic("failed in init of DefaultClientAuthenticationProcessor")
	}
}

func (d *DefaultClientAuthenticationProcessor) HandleAuthEP(ctx context.Context, requestContext sdk.IAuthenticationRequestContext) sdk.IError {
	clientId := requestContext.GetClientID()
	if clientId == "" {
		return sdkerror.ErrInvalidClient.WithDescription("client id not found in request")
	}
	client, err := d.ClientStore.GetClient(ctx, clientId)
	if err != nil {
		return sdkerror.ErrInvalidClient.WithDescription(err.Error())
	}
	requestContext.SetClient(client)
	return nil
}

func (d *DefaultClientAuthenticationProcessor) HandleTokenEP(ctx context.Context, requestContext sdk.ITokenRequestContext) sdk.IError {
	clientId := requestContext.GetClientID()
	if clientId == "" {
		return sdkerror.ErrInvalidClient.WithDescription("client id not found in request")
	}
	clientSecret := requestContext.GetClientSecret()
	client, err := d.ClientStore.GetClient(ctx, clientId)
	if err != nil {
		return sdkerror.ErrInvalidClient.WithDescription(err.Error())
	}
	if clientSecret == "" && client.IsPublic() {
		requestContext.SetClient(client)
		return nil
	}

	//todo handle secret encryption
	if client.GetSecret() == clientSecret {
		requestContext.SetClient(client)
		return nil
	}
	return sdkerror.ErrInvalidClient.WithDescription("could not authenticate client")
}
