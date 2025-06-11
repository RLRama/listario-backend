package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func SetupLogger() *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().
		Timestamp().
		Logger()
	return &logger
}
