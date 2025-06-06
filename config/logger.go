package config

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339
		}).With().Timestamp().Logger()
	log.Info().Msg("Logger initialized with 'Stderr' output")
}