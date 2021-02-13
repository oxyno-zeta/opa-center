package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	logrus "github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

type loggerIns struct {
	logrus.FieldLogger
}

// This is dirty pkg/errors.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (ll *loggerIns) GetCronLogger() cron.Logger {
	return &cronLogger{
		logger: ll.WithField("logger", "cron"),
	}
}

func (ll *loggerIns) GetGormLogger() gormlogger.Interface {
	return &gormLogger{
		logger: ll,
	}
}

func (ll *loggerIns) GetTracingLogger() TracingLogger {
	return &tracingLogger{
		logger: ll,
	}
}

func (ll *loggerIns) GetLockDistributorLogger() LockDistributorLogger {
	return &lockDistributorLogger{
		logger: ll,
	}
}

func (ll *loggerIns) Configure(level string, format string, filePath string) error {
	// Parse log level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Get logrus logger
	lll := ll.FieldLogger.(*logrus.Logger)

	// Set log level
	lll.SetLevel(lvl)

	// Set format
	if format == "json" {
		lll.SetFormatter(&logrus.JSONFormatter{})
	} else {
		lll.SetFormatter(&logrus.TextFormatter{})
	}

	if filePath != "" {
		// Create directory if necessary
		err2 := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err2 != nil {
			return err2
		}

		// Open file
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		// Set output file
		lll.SetOutput(f)
	}

	return nil
}

func (ll *loggerIns) WithField(key string, value interface{}) Logger {
	// Create new field logger
	fieldL := ll.FieldLogger.WithField(key, value)

	return &loggerIns{
		FieldLogger: fieldL,
	}
}

func (ll *loggerIns) WithFields(fields map[string]interface{}) Logger {
	// Transform fields
	var ff logrus.Fields = fields
	// Create new field logger
	fieldL := ll.FieldLogger.WithFields(ff)

	return &loggerIns{
		FieldLogger: fieldL,
	}
}

func (ll *loggerIns) WithError(err error) Logger {
	// Create new field logger
	fieldL := ll.FieldLogger.WithError(err)

	addStackTrace := func(pError stackTracer) {
		// Get stack trace from error
		st := pError.StackTrace()
		// Stringify stack trace
		valued := fmt.Sprintf("%+v", st)
		// Remove all tabs
		valued = strings.ReplaceAll(valued, "\t", "")
		// Split on new line
		stack := strings.Split(valued, "\n")
		// Remove first empty string
		stack = stack[1:]
		// Add stack trace to field logger
		fieldL = fieldL.WithField("stack", strings.Join(stack, ","))
	}

	// Check if error is matching stack trace interface
	// nolint: errorlint // Ignore this because the aim is to catch stack trace error at first level
	if err2, ok := err.(stackTracer); ok {
		addStackTrace(err2)
	}

	// Check if error cause is matching stack trace interface
	// nolint: errorlint // Ignore this because the aim is to catch stack trace error at first level
	if err2, ok := errors.Cause(err).(stackTracer); ok {
		addStackTrace(err2)
	}

	// Return new logger
	return &loggerIns{
		FieldLogger: fieldL,
	}
}
