package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Key      KeyConfig
	S3       S3Config
	Email    EmailConfig
	OAuth    OAuthConfig
}

type AppConfig struct {
	ENV      string `envconfig:"APP_ENV" default:"local"`
	Prefix   string `envconfig:"APP_PREFIX" default:"SDD"`
	Port     string `envconfig:"APP_PORT" default:"8080"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	XAPIKey  string `envconfig:"X_API_KEY"`
	Version  string `envconfig:"APP_VERSION" default:"1.0.0"`
}

// Init loads .env file. Path overridden by ENV_FILE_PATH env var.
func Init() {
	path := os.Getenv("ENV_FILE_PATH")
	if path == "" {
		path = ".env"
	}
	_ = godotenv.Load(path)
}

// GetConfigs loads all config from env. Call once at startup; inject via Dependencies.
func GetConfigs() Config {
	Init()
	var cfg Config
	mustProcess(&cfg.App)
	mustProcess(&cfg.Postgres)
	mustProcess(&cfg.Key)
	mustProcess(&cfg.S3)
	mustProcess(&cfg.Email)
	mustProcess(&cfg.OAuth)
	if cfg.Key.JWTSecret == "" {
		panic("JWT_SECRET must not be empty")
	}
	return cfg
}

func mustProcess(spec interface{}) {
	if err := envconfig.Process("", spec); err != nil {
		panic(err)
	}
}
