package config

type Config struct {
	PostgresUrl string `envconfig:"POSTGRES_URL" required:"true"`
	JwtSecret   string `envconfig:"JWT_SECRET" required:"true"`
	GRPCPort    int32  `envconfig:"GRPC_PORT" default:"50051"`
}
