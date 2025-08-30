package monitoring

import (
	"os"
	"time"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func InitLogger(level string, pretty bool) {
	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	if pretty {
		Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Caller().Logger()
	} else {
		Logger = zerolog.New(os.Stdout).With().
			Timestamp().
			Caller().
			Str("service", "crypto-server").
			Str("version", "1.0.0").
			Logger()
	}

	log.Logger = Logger
}

func LogRequest(method, path string, statusCode int, duration time.Duration, userID string) {
	Logger.Info().
		Str("method", method).
		Str("path", path).
		Int("status_code", statusCode).
		Dur("duration", duration).
		Str("user_id", userID).
		Msg("HTTP request processed")
}

func LogError(err error, context string, fields map[string]interface{}) {
	event := Logger.Error().Err(err).Str("context", context)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg("Error occurred")
}

func LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	if err != nil {
		Logger.Error().
			Err(err).
			Str("operation", operation).
			Str("table", table).
			Dur("duration", duration).
			Msg("Database operation failed")
		RecordDatabaseError(operation, table)
	} else {
		Logger.Debug().
			Str("operation", operation).
			Str("table", table).
			Dur("duration", duration).
			Msg("Database operation completed")
	}
}

func LogExternalAPICall(api, endpoint string, statusCode int, duration time.Duration, err error) {
	if err != nil {
		Logger.Error().
			Err(err).
			Str("api", api).
			Str("endpoint", endpoint).
			Int("status_code", statusCode).
			Dur("duration", duration).
			Msg("External API call failed")
	} else {
		Logger.Info().
			Str("api", api).
			Str("endpoint", endpoint).
			Int("status_code", statusCode).
			Dur("duration", duration).
			Msg("External API call completed")
	}
	RecordExternalAPICall(api, endpoint, duration)
}