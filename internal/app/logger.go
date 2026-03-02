package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(cfg Config) *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	if cfg.App.Debug {
		log.SetLevel(logrus.TraceLevel)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log
}
