package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func Init() zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	return logger
}
