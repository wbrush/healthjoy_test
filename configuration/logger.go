package configuration

import (
	"github.com/sirupsen/logrus"
)

func ConfigureLogger() error {
	//configure logrus
	// will find out if running in real environment (dev/test/prod/etc) and setup stackdriver formatter
	// else, skip to the basics

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})

	//set the debug level from config, if it present
	logrus.SetLevel(logrus.DebugLevel)

	return nil
}
