package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var log = logrus.New()

func Init() {
	// log as JSON instead of the default ASCII formatter.

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: false})

	log.SetReportCaller(false)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Only log the Debug severity or above.
	// Will log anything that is Debug or above (info, warn, error, fatal, panic).
	log.SetLevel(logrus.DebugLevel)
}

func Infof(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"message": format,
	}).Infof(format, args...)
}

func Info(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{}).Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{}).Errorf(format, args...)
}

func Error(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"status":       args[0],
		"message":      args[1],
		"err":          args[2],
		"requestBody":  args[3],
		"responseBody": args[4],
	}).Error(args[1])
}

func Fatalf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{}).Fatalf(format, args...)
}
