package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

type Files struct {
	App          *os.File
	Access       *os.File
	AccessWriter io.Writer
}

func Init(appPath, accessPath string) (*Files, error) {
	appFile, err := openLogFile(appPath)
	if err != nil {
		return nil, err
	}

	accessFile, err := openLogFile(accessPath)
	if err != nil {
		appFile.Close()
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	Logger = zerolog.New(io.MultiWriter(console, appFile)).
		With().
		Timestamp().
		Logger()

	return &Files{
		App:          appFile,
		Access:       accessFile,
		AccessWriter: io.MultiWriter(os.Stdout, accessFile),
	}, nil
}

func (f *Files) Close() {
	_ = f.Access.Close()
	_ = f.App.Close()
}

func openLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}
