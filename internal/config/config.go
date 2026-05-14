package config

type Config struct {
	PostgresUrl string `envconfig:"POSTGRES_URL" required:"true"`
	JwtSecret   string `envconfig:"JWT_SECRET" required:"true"`
}
