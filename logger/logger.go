package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func SetupLogger() {

	logFormat := os.Getenv("LOG_FORMAT")
	logLevelStr := os.Getenv("LOG_LEVEL")

	switch logFormat {
	case "pretty":
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}

		output.FormatLevel = func(i interface{}) string {
			if i == nil {
				return strings.ToUpper(fmt.Sprintf("| %6s|", "UNKWN"))
			}
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}

		Logger = zerolog.New(output).With().Timestamp().Logger()
		Logger.Info().Msg("Logger configured for pretty output.")
	case "json", "":
		Logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		Logger.Info().Msg("Logger configured or default for structured (JSON) output.")
	}

	var level zerolog.Level
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(logLevelStr))
	if err != nil {
		level = zerolog.TraceLevel
		Logger.Warn().Msgf("Invalid LOG_LEVEL environment variable '%s', defaulting to 'trace'", logLevelStr)
	} else {
		level = parsedLevel
		Logger.Info().Str("log_level", level.String()).Msg("Global log level set")
	}

	Logger = Logger.Level(level)
	zerolog.SetGlobalLevel(level)
	log.Logger = Logger
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
