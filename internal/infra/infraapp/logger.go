package infraapp

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type AppLogger struct {
	*logrus.Logger
}

func NewLogger() *AppLogger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	})
	l.SetLevel(logrus.InfoLevel)

	return &AppLogger{Logger: l}
}
