package logger

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	// Logger Config
	logLevel      = "debug"
	logOutputFile = "/home/steampipe/logfile"

	timestampFormat = "02-01-2006 15:04:05.000Z"
)

func SetUpLogger() error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return errors.Wrapf(err, "Unknown log level '%s'", logLevel)
	}

	log.SetLevel(level)

	var logFile *os.File
	// If the output isn't stdout it should be a file path
	if _, err = os.Stat(logOutputFile); err == nil {
		logFile, err = os.OpenFile(logOutputFile, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return errors.Wrapf(err, "failed to open log file %s for output", logOutputFile)
		}
	} else if err != nil && errors.Is(err, os.ErrNotExist) {
		if logFile, err = os.Create(logOutputFile); err != nil {
			return errors.Wrap(err, "failed to create log file")
		}
	} else {
		return errors.Errorf("failed to find log file '%s' for output", logOutputFile)
	}

	log.SetOutput(logFile)
	log.RegisterExitHandler(func() {
		if logFile != nil {
			_ = logFile.Close() // Ignoring error because we can do nothing with it
		}
	})

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp:   true,
	})

	return nil
}

func GetLogs() string {
	logData, err := os.ReadFile(logOutputFile)
	if err != nil {
		return fmt.Sprintf("Faild to read logfile: %s", err.Error())
	}

	return string(logData)
}
