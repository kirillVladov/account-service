package config

type Config struct {
	PostgresUrl string `envconfig:"POSTGRES_URL" required:"true"`
	GRPCPort    int32  `envconfig:"GRPC_PORT" default:"50051"`
	DebugPort   int32  `envconfig:"DEBUG_PORT" required:"true"`
}
