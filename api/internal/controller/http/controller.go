package http

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/DataDog/gostackparse"
	"github.com/gin-gonic/gin"
	"github.com/softcery/shopify-app-template-go/config"
	"github.com/softcery/shopify-app-template-go/internal/service"
	"github.com/softcery/shopify-app-template-go/pkg/logging"
)

// Options is used to create HTTP controller.
type Options struct {
	Handler  *gin.Engine
	Services service.Services
	Storages service.Storages
	Logger   logging.Logger
	Config   *config.Config
}

// RouterOptions provides shared options for all routers.
type RouterOptions struct {
	Handler  *gin.RouterGroup
	Services service.Services
	Storages service.Storages
	Logger   logging.Logger
	Config   *config.Config
}

// RouterContext provides a shared context for all routers.
type RouterContext struct {
	services service.Services
	storages service.Storages
	logger   logging.Logger
	cfg      *config.Config
}

func New(options *Options) {
	options.Handler.Use(
		corsMiddleware,
	)

	routerOptions := RouterOptions{
		Handler:  options.Handler.Group(""),
		Services: options.Services,
		Storages: options.Storages,
		Logger:   options.Logger.Named("HTTPController"),
		Config:   options.Config,
	}

	// K8S probe
	options.Handler.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	{
		newPlatformRoutes(routerOptions)
	}
}

// httpErr provides a base error type for all http controller errors.
type httpErr struct {
	Type             httpErrType            `json:"-"`
	Code             int                    `json:"-"`
	Message          string                 `json:"message"`
	Details          interface{}            `json:"details,omitempty"`
	ValidationErrors map[string]interface{} `json:"validationErrors,omitempty"`
}

// httpErrType is used to define error type.
type httpErrType string

const (
	// ErrorTypeServer is an "unexpected" internal server error.
	ErrorTypeServer httpErrType = "server"
	// ErrorTypeClient is an "expected" business error.
	ErrorTypeClient httpErrType = "client"
)

// Error is used to convert an error to a string.
func (e *httpErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// wrapHandler provides unified error handling for all handlers.
func wrapHandler(options RouterOptions, handler func(c *gin.Context) (interface{}, *httpErr)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := options.Logger.Named("wrapHandler")

		// handle panics
		defer func() {
			if err := recover(); err != nil {
				// get stacktrace
				stacktrace, errors := gostackparse.Parse(bytes.NewReader(debug.Stack()))
				if len(errors) > 0 || len(stacktrace) == 0 {
					logger.Error("get stacktrace errors", "stacktraceErrors", errors, "stacktrace", "unknown", "err", err)
				} else {
					logger.Error("unhandled error", "err", err, "stacktrace", stacktrace)
				}

				// return error
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
				if err != nil {
					logger.Error("failed to abort with error", "err", err)
				}
			}
		}()

		c.Set("Authorization", c.Request.Header.Get("Authorization"))

		// execute handler
		body, err := handler(c)

		// check if middleware
		if body == nil && err == nil {
			return
		}
		logger = logger.With("body", body).With("err", err)

		// check error
		if err != nil {
			if err.Type == ErrorTypeServer {
				logger.Error("internal server error")

				// whether to send error to the client
				if options.Config.HTTP.SendDetailsOnInternalError {
					// send error to the client
					c.AbortWithStatusJSON(http.StatusInternalServerError, err)
				} else {
					// don't send error to the client
					err := c.AbortWithError(http.StatusInternalServerError, err)
					if err != nil {
						logger.Error("failed to abort with error", "err", err)
					}
					logger.Info("aborted with error")
				}
			} else {
				logger.Info("client error")
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			}
			return
		}
		logger.Info("request handled")
		c.JSON(http.StatusOK, body)
	}
}

// corsMiddleware is used to allow incoming cross-origin requests.
func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
