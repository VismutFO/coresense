package config

import (
	"time"
)

type (
	DatabaseLogger struct {
		SlowThreshold  time.Duration `json:"slow_threshold"`
		WithParameters bool          `json:"with_parameters"`
	}

	Database struct {
		URL                   string         `json:"url" sensitive:""`
		MaxOpenConnections    int            `json:"maxopenconnections"`
		MaxIdleConnections    int            `json:"maxidleconnections"`
		ConnectionMaxLifeTime time.Duration  `json:"connectionmaxlifetime"`
		Logger                DatabaseLogger `json:"logger"`
	}

	Logger struct {
		Level     int    `json:"level"`
		Formatter string `json:"formatter"`
		Indent    string `json:"indent"`
	}
)

func (d *Database) HasURL() bool {
	return d.URL != ""
}
