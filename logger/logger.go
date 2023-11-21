package logger

import (
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/sirupsen/logrus"
)

func GetLogger(c *config.Config) *logrus.Logger {
	logger := logrus.StandardLogger()

	level, err := logrus.ParseLevel(c.Log.Level)

	if err != nil {
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)

	return logger
}
