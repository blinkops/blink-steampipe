package logger

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	// Logger Config
	logLevel      = "info"
	logOutputFile = "/home/steampipe/logfile"

	timestampFormat = "02-01-2006 15:04:05.000Z"
)

func SetUpLogger() error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return errors.Wrapf(err, "Unknown log level '%s'", logLevel)
	}

	log.SetLevel(level)

	// If the output isn't stdout it should be a file path
	if _, err := os.Stat(logOutputFile); err == nil {
		logFile, err := os.OpenFile(logOutputFile, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return errors.Wrapf(err, "Failed to open log file %s for output", logOutputFile)
		}

		log.SetOutput(logFile)
		log.RegisterExitHandler(func() {
			if logFile == nil {
				return
			}
			err := logFile.Close()
			if err != nil {
				return
			}
		})
	} else if errors.Is(err, os.ErrNotExist) {
		return errors.Errorf("Failed to find log file '%s' for output", logOutputFile)
	}

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
