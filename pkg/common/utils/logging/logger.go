package logging

import (
	"os"

	"github.com/rs/zerolog"
)

func NewDefault() zerolog.Logger {
	return zerolog.New(os.Stdout).With().
		Int("ProcessID", os.Getpid()).
		Timestamp().
		Logger()
}
