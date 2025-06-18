package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func SetupLogger() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s |", i))
	}

	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", fmt.Sprintf("%s", i))
	}

	output.FormatFieldValue = func(i interface{}) string {

	}
}
