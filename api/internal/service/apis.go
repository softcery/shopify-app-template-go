package service

import (
	"context"

	"github.com/softcery/shopify-app-template-go/internal/entity"
	"github.com/softcery/shopify-app-template-go/pkg/errs"
)

// APIs provides a collection of API interfaces.
type APIs struct {
	Platform PlatformAPI
}

// PlatformAPI is used to communicate with shop platform.
type PlatformAPI interface {
	// HandleInstall verifies installation URL and returns url to redirect user to.
	HandleInstall(opts HandleInstallOptions) (APIHandleInstallOutput, error)
	// HandleRedirect verifies redirected URL and requests access token from shop platform
	// and then returns the access token.
	HandleRedirect(opts APIHandleRedirectOptions) (string, error)
	// SubscribeToAppUninstallWebhook subscribes application to platform's webhook.
	SubscribeToAppUninstallWebhook(opts SubscribeToAppUninstallWebhookOptions) error
	// VerifySession verifies session and returns true if session is valid.
	VerifySession(ctx context.Context) (*VerifySessionOutput, error)
	// WithConfig returns a new instance of PlatformAPI with provided store config.
	WithConfig(ctx context.Context, store *entity.Store) PlatformAPI
	// CreateProducts creates random products in shopify store.
	CreateProducts(ctx context.Context) error
	// GetProductsCount returns number of products in store.
	GetProductsCount(ctx context.Context) (int, error)
}

var (
	// ErrHandleRedirectInvalidRedirectedURL is returned when provided redirected URL is invalid.
	ErrHandleRedirectInvalidRedirectedURL = errs.New("invalid redirected url")
	// ErrHandleRedirectInvalidScopes is returned when user didn't allow all the requested scopes when installing app.
	ErrHandleRedirectInvalidScopes = errs.New("allowed access scopes are different from requested")
)

type VerifySessionOutput struct {
	StoreName  string
	IsVerified bool
}

type HandleInstallOptions struct {
	InstallationURL string
	RedirectURL     string
	StoreName       string
}

type APIHandleInstallOutput struct {
	Nonce       string
	RedirectURL string
}

type APIHandleRedirectOptions struct {
	Nonce         string
	RedirectedURL string
	StoreName     string
}

type SubscribeToAppUninstallWebhookOptions struct {
	RedirectURL string
	StoreName   string
	AccessToken string
}
