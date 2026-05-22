package configs

type KeyConfig struct {
	JWTSecret string `envconfig:"JWT_SECRET"`
}
