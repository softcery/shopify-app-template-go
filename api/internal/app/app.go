package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/softcery/shopify-app-template-go/config"
	"github.com/softcery/shopify-app-template-go/internal/api/shopify"
	httpcontroller "github.com/softcery/shopify-app-template-go/internal/controller/http"
	"github.com/softcery/shopify-app-template-go/internal/entity"
	"github.com/softcery/shopify-app-template-go/internal/service"
	"github.com/softcery/shopify-app-template-go/internal/storage"
	"github.com/softcery/shopify-app-template-go/pkg/database"
	"github.com/softcery/shopify-app-template-go/pkg/httpserver"
	"github.com/softcery/shopify-app-template-go/pkg/logging"
)

func Run(cfg *config.Config) {
	logger := logging.NewZap(cfg.Log.Level)

	// Init db
	sql, err := database.NewPostgreSQL(&database.PostgreSQLConfig{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Database: cfg.Postgres.Database,
	})
	if err != nil {
		logger.Fatal("failed to connect to PostgreSQL", "err", err)
	}

	err = sql.DB.AutoMigrate(
		&entity.Store{},
	)
	if err != nil {
		logger.Fatal("automigration failed", "err", err)
	}

	storages := service.Storages{
		Store: storage.NewStoreStorage(sql),
	}

	apis := service.APIs{
		Platform: shopify.NewAPI(shopify.Options{
			Config: cfg,
			Logger: logger,
		}),
	}

	serviceOptions := &service.Options{
		Apis:     apis,
		Storages: storages,
		Config:   cfg,
		Logger:   logger,
	}

	services := service.Services{
		Platform: service.NewPlatformService(serviceOptions),
	}

	// Init HTTP framework of choice
	httpHandler := gin.New()

	httpcontroller.New(&httpcontroller.Options{
		Handler:  httpHandler,
		Services: services,
		Storages: storages,
		Logger:   logger,
		Config:   cfg,
	})

	httpServer := httpserver.New(
		httpHandler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(120*time.Second),
		httpserver.WriteTimeout(120*time.Second),
		httpserver.ShutdownTimeout(30*time.Second),
	)

	// Waiting for a signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify", "err", err)
	}

	// Shutdown HTTP server
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
