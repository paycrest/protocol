package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	sentryhook "github.com/chadsr/logrus-sentry"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger = logrus.New()

func init() {
	logger.Level = logrus.InfoLevel
	logger.Formatter = &formatter{}

	sentryDSN, environment := getSentryConfig()

	if environment == "production" {
		// init sentry
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDSN,
		})
		if err != nil {
			logger.Fatalf("Sentry initialization failed: %v", err)
		}

		// Sentry hook
		hook := sentryhook.New([]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		logger.Hooks.Add(hook)
	}
	logger.SetReportCaller(true)
}

// SetLogLevel sets the log level for the logger.
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

func getSentryConfig() (string, string) {
	viper.AddConfigPath(".")
	envFilePath := os.Getenv("ENV_FILE_PATH")
	if envFilePath == "" {
		envFilePath = ".env" // Set default value to ".env"
	}
	viper.SetConfigName(envFilePath)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Print("Error to reading config file, ", err, "\n")
	}
	return viper.GetString("SENTRY_DSN"), viper.GetString("ENVIRONMENT")
}

// Fields type, used to pass to `WithFields`.
type Fields logrus.Fields

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Infof(format, args...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}

// Formatter implements logrus.Formatter interface.
type formatter struct {
	prefix string
}

// Format building log message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	sb.WriteString(strings.ToUpper(entry.Level.String()))
	sb.WriteString(" ")
	sb.WriteString(entry.Time.Format(time.RFC3339))
	sb.WriteString(" ")
	sb.WriteString(f.prefix)
	sb.WriteString(entry.Message)

	return sb.Bytes(), nil
}
