package main

import (
	// external
	"github.com/softcery/shopify-app-template-go/config"
	"github.com/softcery/shopify-app-template-go/pkg/logging"

	// internal
	"github.com/softcery/shopify-app-template-go/internal/app"
)

func main() {
	logger := logging.NewZap("main")

	cfg := config.Get()
	logger.Info("read config", "config", cfg)

	app.Run(cfg)
}
