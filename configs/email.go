package configs

type EmailConfig struct {
	Host     string `envconfig:"EMAIL_HOST"`
	Port     int    `envconfig:"EMAIL_PORT" default:"587"`
	User     string `envconfig:"EMAIL_USER"`
	Password string `envconfig:"EMAIL_PASSWORD"`
	From     string `envconfig:"EMAIL_FROM"`
}
