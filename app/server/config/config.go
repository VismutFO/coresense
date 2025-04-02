package config

import (
	"time"

	"github.com/rs/zerolog"

	"coresense/pkg/common/config"
)

type (
	AppConfig struct {
		ServiceName string
		Database    config.Database
		Logger      config.Logger
		Server      Server
	}

	Server struct {
		HTTP HTTP `json:"http"`
	}

	HTTP struct {
		Port int `json:"port"`
	}
)

func New(_ string, _ zerolog.Logger) (*AppConfig, error) {
	return &AppConfig{
		ServiceName: "api",
		Database: config.Database{
			URL:                   "host=localhost port=5433 user=postgres password=postgres dbname=claimix sslmode=disable search_path=claimix",
			MaxOpenConnections:    3,
			MaxIdleConnections:    3,
			ConnectionMaxLifeTime: time.Duration(1000000),
			Logger: config.DatabaseLogger{
				SlowThreshold:  time.Duration(1000000),
				WithParameters: false,
			},
		},
		Logger: config.Logger{
			Level:     0,
			Formatter: "json",
			Indent:    "",
		},
		Server: Server{
			HTTP: HTTP{
				Port: 9998,
			},
		},
	}, nil
}
