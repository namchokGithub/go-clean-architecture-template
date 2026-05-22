package configs

type OAuthConfig struct {
	GoogleClientID     string `envconfig:"OAUTH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `envconfig:"OAUTH_GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `envconfig:"OAUTH_GOOGLE_CALLBACK_URL"`
}
