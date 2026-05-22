package configs

type PostgresConfig struct {
	Host         string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port         string `envconfig:"POSTGRES_PORT" default:"5432"`
	User         string `envconfig:"POSTGRES_USER" default:"postgres"`
	Password     string `envconfig:"POSTGRES_PASSWORD"`
	DBName       string `envconfig:"POSTGRES_DB" default:"superdriver"`
	SSLMode      string `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
	MaxOpenConns int    `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns int    `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"5"`
}
