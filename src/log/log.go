package log

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var log *logrus.Logger
var once sync.Once

func GetLogger() *logrus.Logger {
	once.Do(func() {
		log = logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{})
		switch getEnv("LOG_LEVEL", "info") {
		case "debug":
			log.SetLevel(logrus.DebugLevel)
		case "info":
			log.SetLevel(logrus.InfoLevel)
		case "warn":
			log.SetLevel(logrus.WarnLevel)
		case "fatal":
			log.SetLevel(logrus.FatalLevel)
		case "panic":
			log.SetLevel(logrus.PanicLevel)
		}
	})
	return log
}
