package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Init(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	Logger = zerolog.New(io.MultiWriter(console, logFile)).
		With().
		Timestamp().
		Logger()

	return logFile, nil
}
