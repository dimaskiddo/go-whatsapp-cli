package hlp

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type logLevel string

const (
	LogLevelPanic logLevel = "panic"
	LogLevelFatal logLevel = "fatal"
	LogLevelError logLevel = "error"
	LogLevelWarn  logLevel = "warn"
	LogLevelDebug logLevel = "debug"
	LogLevelTrace logLevel = "trace"
	LogLevelInfo  logLevel = "info"
)

func init() {
	log = logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func LogPrintln(level logLevel, message interface{}) {
	if log != nil {
		switch level {
		case "panic":
			log.Panicln(message)
		case "fatal":
			log.Fatalln(message)
		case "error":
			log.Errorln(message)
		case "warn":
			log.Warnln(message)
		case "debug":
			log.Debugln(message)
		case "trace":
			log.Traceln(message)
		default:
			log.Infoln(message)
		}
	}
}
