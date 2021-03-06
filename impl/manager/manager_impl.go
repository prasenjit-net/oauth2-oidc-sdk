package manager

import (
	sdk "github.com/identityOrg/oidcsdk"
	"github.com/identityOrg/oidcsdk/impl/strategies"
)

type (
	DefaultManager struct {
		Config                  *sdk.Config
		RequestContextFactory   sdk.IRequestContextFactory
		ErrorWriter             sdk.IErrorWriter
		ResponseWriter          sdk.IResponseWriter
		ErrorStrategy           sdk.ErrorStrategy
		UserSessionManager      sdk.ISessionManager
		SecretStore             sdk.ISecretStore
		PageResponseHandler     sdk.IPageResponseHandler
		RPILogoutEPHandlers     []sdk.IRPILogoutEPHandler
		AuthEPHandlers          []sdk.IAuthEPHandler
		TokenEPHandlers         []sdk.ITokenEPHandler
		IntrospectionEPHandlers []sdk.IIntrospectionEPHandler
		RevocationEPHandlers    []sdk.IRevocationEPHandler
		UserInfoEPHandlers      []sdk.IUserInfoEPHandler
	}
)

func (d *DefaultManager) SetErrorStrategy(strategy sdk.ErrorStrategy) {
	d.ErrorStrategy = strategy
}

func NewDefaultManager(config *sdk.Config, options *Options) *DefaultManager {
	manager := &DefaultManager{
		Config:                config,
		RequestContextFactory: options.RequestContextFactory,
		ErrorWriter:           options.ErrorWriter,
		ResponseWriter:        options.ResponseWriter,
		UserSessionManager:    options.UserSessionManager,
		SecretStore:           options.SecretStore,
		PageResponseHandler:   options.PageResponseHandler,
	}
	for _, arg := range options.Sequence {
		if element, ok := arg.(sdk.IAuthEPHandler); ok {
			manager.AuthEPHandlers = append(manager.AuthEPHandlers, element)
		}
		if element, ok := arg.(sdk.ITokenEPHandler); ok {
			manager.TokenEPHandlers = append(manager.TokenEPHandlers, element)
		}
		if element, ok := arg.(sdk.IIntrospectionEPHandler); ok {
			manager.IntrospectionEPHandlers = append(manager.IntrospectionEPHandlers, element)
		}
		if element, ok := arg.(sdk.IRevocationEPHandler); ok {
			manager.RevocationEPHandlers = append(manager.RevocationEPHandlers, element)
		}
		if element, ok := arg.(sdk.IUserInfoEPHandler); ok {
			manager.UserInfoEPHandlers = append(manager.UserInfoEPHandlers, element)
		}
		if element, ok := arg.(sdk.IRPILogoutEPHandler); ok {
			manager.RPILogoutEPHandlers = append(manager.RPILogoutEPHandlers, element)
		}
	}
	manager.ErrorStrategy = strategies.DefaultLoggingErrorStrategy
	return manager
}

type Options struct {
	PageResponseHandler   sdk.IPageResponseHandler
	RequestContextFactory sdk.IRequestContextFactory
	ErrorWriter           sdk.IErrorWriter
	ResponseWriter        sdk.IResponseWriter
	UserSessionManager    sdk.ISessionManager
	SecretStore           sdk.ISecretStore
	Sequence              []interface{}
}
