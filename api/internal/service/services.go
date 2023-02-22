package service

import (
	"context"

	"github.com/softcery/shopify-app-template-go/config"
	"github.com/softcery/shopify-app-template-go/pkg/errs"
	"github.com/softcery/shopify-app-template-go/pkg/logging"
)

// Services contains all available services.
type Services struct {
	Platform PlatformService
}

// Options provides options for creating a new service instance.
type Options struct {
	Apis     APIs
	Storages Storages
	Config   *config.Config
	Logger   logging.Logger
}

// PlatformService provides business logic related to shop platformService.
type PlatformService interface {
	Handle(ctx context.Context, storeName, installationURL string) (string, error)
	// HandleRedirect handles an oauth2 redirect call for a platform integration.
	HandleRedirect(ctx context.Context, opts ServiceHandleRedirectOptions) error
	// HandleUninstall is called when user wants to uninstall the app from a platform.
	// In this case we need to delete all records about their store from database.
	HandleUninstall(ctx context.Context, storeName string) error
	// GetProductsCount returns number of products in store.
	GetProductsCount(ctx context.Context) (int, error)
	// CreateProducts creates random products in store.
	CreateProducts(ctx context.Context) error
}

const (
	DEFAULT_PRODUCT_COUNT = 5
)

var (
	// ErrHandleRedirectStoreNotFound is returned when store is not found.
	ErrHandleRedirectStoreNotFound = errs.New("store is not found")

	// ErrHandleUninstallStoreNotFound is returned when store is not found.
	ErrHandleUninstallStoreNotFound = errs.New("store is not found")
)

type ServiceHandlerOptions struct {
	StoreName       string
	InstallationURL string
}

type ServiceHandleInstallOptions struct {
	StoreName       string
	InstallationURL string
}

type ServiceHandleRedirectOptions struct {
	StoreName     string
	RedirectedURL string
}
