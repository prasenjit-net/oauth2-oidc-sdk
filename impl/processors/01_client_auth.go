package processors

import (
	"context"
	sdk "github.com/identityOrg/oidcsdk"
	"github.com/identityOrg/oidcsdk/impl/sdkerror"
)

type DefaultClientAuthenticationProcessor struct {
	ClientStore sdk.IClientStore
}

func NewDefaultClientAuthenticationProcessor(clientStore sdk.IClientStore) *DefaultClientAuthenticationProcessor {
	return &DefaultClientAuthenticationProcessor{ClientStore: clientStore}
}

func (d *DefaultClientAuthenticationProcessor) HandleRevocationEP(ctx context.Context, requestContext sdk.IRevocationRequestContext) sdk.IError {
	return d.authenticateClient(ctx, requestContext)
}

func (d *DefaultClientAuthenticationProcessor) HandleIntrospectionEP(ctx context.Context, requestContext sdk.IIntrospectionRequestContext) sdk.IError {
	return d.authenticateClient(ctx, requestContext)
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
	requestContext.GetProfile().SetClientID(clientId)
	return nil
}

func (d *DefaultClientAuthenticationProcessor) HandleTokenEP(ctx context.Context, requestContext sdk.ITokenRequestContext) sdk.IError {
	if iError := d.authenticateClient(ctx, requestContext); iError != nil {
		return iError
	} else {
		requestContext.GetProfile().SetClientID(requestContext.GetClient().GetID())
		return nil
	}
}

func (d *DefaultClientAuthenticationProcessor) authenticateClient(ctx context.Context, requestContext sdk.IClientCredentialContext) sdk.IError {
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
