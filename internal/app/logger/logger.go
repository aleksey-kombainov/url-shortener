package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func GetLogger(logLvl zerolog.Level) zerolog.Logger {
	return zerolog.New(os.Stdout).Level(logLvl)
}
