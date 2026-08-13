package config

import "time"

type AuthToken struct {
	RefreshTokenTTL        time.Duration `envconfig:"REFRESH_TOKEN_TTL" required:"true"`
	TokenTTL               time.Duration `envconfig:"TOKEN_TTL" required:"true"`
	RefreshTokenBytesCount int64         `envconfig:"REFRESH_TOKEN_BYTES_COUNT" required:"true"`
	TokenSecret            string        `envconfig:"TOKEN_SECRET" required:"true"`
}

type Config struct {
	PostgresUrl string `envconfig:"POSTGRES_URL" required:"true"`
	GRPCPort    int32  `envconfig:"GRPC_PORT" default:"50051"`
	DebugPort   int32  `envconfig:"DEBUG_PORT" required:"true"`

	AuthToken AuthToken
}
