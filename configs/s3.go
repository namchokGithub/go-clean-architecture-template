package configs

type S3Config struct {
	Bucket    string `envconfig:"S3_BUCKET"`
	Region    string `envconfig:"S3_REGION" default:"ap-southeast-1"`
	AccessKey string `envconfig:"S3_ACCESS_KEY"`
	SecretKey string `envconfig:"S3_SECRET_KEY"`
	Endpoint  string `envconfig:"S3_ENDPOINT"`
}
