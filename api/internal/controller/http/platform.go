package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softcery/shopify-app-template-go/internal/service"
	"github.com/softcery/shopify-app-template-go/pkg/errs"
)

type platformRoutes struct {
	RouterContext
}

func newPlatformRoutes(options RouterOptions) {
	r := &platformRoutes{RouterContext{
		services: options.Services,
		storages: options.Storages,
		logger:   options.Logger.Named("platformRoutes"),
		cfg:      options.Config,
	}}

	p := options.Handler.Group("")
	{
		p.GET("", wrapHandler(options, r.handler))
		p.GET("/auth/callback", wrapHandler(options, r.redirectHandler))
		p.POST("/uninstall", wrapHandler(options, r.uninstallHandler))
		p.GET("/api/products/count", wrapHandler(options, r.getProductsCount))
		p.GET("/api/products/create", wrapHandler(options, r.createProducts))
	}
}

type handlerRequestQuery struct {
	StoreName string `form:"shop" binding:"required"`
}

func (r *platformRoutes) handler(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("handler").WithContext(c)

	var requestQuery handlerRequestQuery
	err := c.ShouldBindQuery(&requestQuery)
	if err != nil {
		logger.Info("failed to parse request query", "err", err)
		return nil, &httpErr{Type: ErrorTypeClient, Message: "invalid request query", Details: err}
	}
	logger = logger.With("requestQuery", requestQuery)

	redirectURL, err := r.services.Platform.Handle(c, requestQuery.StoreName, c.Request.URL.String())
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: ErrorTypeClient, Message: err.Error()}
		}
		logger.Error("failed to handle call", "err", err)
		return nil, &httpErr{
			Type:    ErrorTypeServer,
			Message: "failed to handle call",
			Details: err,
		}
	}
	logger = logger.With("redirectURL", redirectURL)

	c.Redirect(http.StatusFound, redirectURL)

	logger.Info("successfully handled call")
	return nil, nil
}

type redirectHandlerRequestQuery struct {
	StoreName string `form:"shop" binding:"required"`
}

func (r *platformRoutes) redirectHandler(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("redirectHandler").WithContext(c)

	var requestQuery redirectHandlerRequestQuery
	err := c.ShouldBindQuery(&requestQuery)
	if err != nil {
		logger.Info("failed to parse request query", "err", err)
		return nil, &httpErr{Type: ErrorTypeClient, Message: "invalid request query", Details: err}
	}
	logger = logger.With("requestQuery", requestQuery)

	err = r.services.Platform.HandleRedirect(c, service.ServiceHandleRedirectOptions{
		StoreName:     requestQuery.StoreName,
		RedirectedURL: c.Request.URL.String(),
	})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: ErrorTypeClient, Message: err.Error()}
		}
		logger.Error("failed to handle oauth2 redirect call", "err", err)
		return nil, &httpErr{
			Type:    ErrorTypeServer,
			Message: "failed to handle oauth2 redirect call",
			Details: err,
		}
	}
	// After successful handling of redirect call, redirect user to app's UI at their platform store
	c.Redirect(http.StatusFound, fmt.Sprintf("https://%s/admin/apps/%s", requestQuery.StoreName, r.cfg.Shopify.ApiKey))

	logger.Info("successfully handled redirect call")
	return nil, nil
}

type uninstallHandlerRequestQuery struct {
	StoreName string `form:"shop" binding:"required"`
}

func (r *platformRoutes) uninstallHandler(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.
		Named("uninstallHandler").
		WithContext(c)

	var requestQuery uninstallHandlerRequestQuery
	err := c.ShouldBindQuery(&requestQuery)
	if err != nil {
		logger.Info("failed to parse request query", "err", err)
		return nil, &httpErr{Type: ErrorTypeClient, Message: "invalid request query", Details: err}
	}
	logger = logger.With("requestQuery", requestQuery)

	err = r.services.Platform.HandleUninstall(c, requestQuery.StoreName)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: ErrorTypeClient, Message: err.Error()}
		}
		logger.Error("failed to uninstall app", "err", err)
		return nil, &httpErr{
			Type:    ErrorTypeServer,
			Message: "failed to failed to uninstall app",
			Details: err,
		}
	}

	logger.Info("successfully uninstalled the application")
	return nil, nil
}

type getProductsCountResponse struct {
	Count int `json:"count"`
}

func (r *platformRoutes) getProductsCount(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("getProductsCount")

	c.Set("Authorization", c.Request.Header.Get("Authorization"))

	count, err := r.services.Platform.GetProductsCount(c)
	if err != nil {
		// TODO: return custom errors to client, instead of 500
		logger.Error("failed to create products", "err", err)
		return nil, &httpErr{
			Type:    ErrorTypeServer,
			Message: "failed to create products",
			Details: err,
		}
	}
	logger = logger.With("count", count)

	logger.Info("successfully got products count")
	return getProductsCountResponse{Count: count}, nil
}

func (r *platformRoutes) createProducts(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("createProducts")

	c.Set("Authorization", c.Request.Header.Get("Authorization"))

	err := r.services.Platform.CreateProducts(c)
	if err != nil {
		// TODO: return custom errors to client, instead of 500
		logger.Error("failed to create products", "err", err)
		return nil, &httpErr{
			Type:    ErrorTypeServer,
			Message: "failed to create products",
			Details: err,
		}
	}

	logger.Info("successfully create products")
	return "", nil
}
