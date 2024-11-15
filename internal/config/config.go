package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"order_service/pkg/db/cache"
	"order_service/pkg/db/postgres"
)

type Config struct {
	cache.RedisConfig
	postgres.Config

	GRPCServerPort int `env:"GRPC_SERVER_PORT" env-default:"50051"`
	RestServerPort int `env:"REST_SERVER_PORT" env-default:"8080"`
}

func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &cfg
}
