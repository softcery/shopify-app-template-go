package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App
		Shopify  Shopify
		HTTP     HTTP
		Log      Log
		Postgres Postgres
	}

	App struct {
		Name    string `env:"APP_NAME" env-default:"app-api"`
		Version string `env:"APP_VERSION" env-default:"1.0.0"`
		BaseURL string `env:"HOST" env-default:"localhost"`
	}

	Shopify struct {
		ApiKey    string `env:"SHOPIFY_API_KEY" env-default:""`
		ApiSecret string `env:"SHOPIFY_API_SECRET" env-default:""`
		Scopes    string `env:"SCOPES" env-default:""`
	}

	HTTP struct {
		Port                       string `env:"BACKEND_PORT" env-default:"8080"`
		SendDetailsOnInternalError bool   `env:"HTTP_SEND_DETAILS_ON_INTERNAL_ERROR" env-default:"true"`
	}

	Postgres struct {
		User     string `env:"POSTGRES_USER" env-default:"postgres"`
		Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Database string `env:"POSTGRES_DATABASE" env-default:"api"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug"`
	}
)

var (
	config Config
	once   sync.Once
)

// Get returns a config.
// Get loads the config only once.
func Get() *Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			log.Fatal("failed to read env", err)
		}
	})

	return &config
}
