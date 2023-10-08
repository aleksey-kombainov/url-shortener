package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func GetLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
}
