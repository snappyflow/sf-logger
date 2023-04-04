package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_LOG_TIMESTAMP_FORMAT = "2006/02/01 15:04:05"
	DEFAULT_LOG_FORMAT           = "%time% [%level%] %filename%:%line% %message% {%fields%}"
	COLON_SEPARATOR              = ":"
)

type Formatter struct {
	TimestampFormat        string
	LogFormat              string
	DisableLevelTruncation bool
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = DEFAULT_LOG_FORMAT
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DEFAULT_LOG_TIMESTAMP_FORMAT
	}
	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%message%", entry.Message, 1)
	if entry.HasCaller() {
		output = strings.Replace(output, "%filename%", entry.Caller.File, 1)
		output = strings.Replace(output, "%line%", strconv.Itoa(entry.Caller.Line), 1)
	} else {
		output = strings.Replace(output, "%filename%", "", 1)
		output = strings.Replace(output, "%line%", "", 1)
	}
	level := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation {
		level = level[0:4]
	}
	output = strings.Replace(output, "%level%", level, 1)

	var fieldFormatter string

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + v + ", "
		case int:
			s := strconv.Itoa(v)
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + s + ", "
		case int32:
			s := strconv.Itoa(int(v))
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + s + ", "
		case int64:
			s := strconv.Itoa(int(v))
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + s + ", "
		case bool:
			s := strconv.FormatBool(v)
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + s + ", "
		case float32:
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + fmt.Sprintf("%.2f", v) + ", "
		case float64:
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + fmt.Sprintf("%.2f", v) + ", "
		case error:
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + v.Error() + ", "
		case time.Time:
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + v.String() + ", "
		default:
			fieldFormatter = fieldFormatter + k + COLON_SEPARATOR + fmt.Sprint(v) + ", "
		}
	}

	output = strings.ReplaceAll(output, "%fields%", strings.TrimSuffix(fieldFormatter, ", "))

	output = output + "\n"

	return []byte(output), nil
}
func GetLogger(debug bool, logFormat string) (initializedLogger *logrus.Logger) {
	initializedLogger = &logrus.Logger{
		Hooks: make(logrus.LevelHooks),
		Out:   os.Stderr,
		Formatter: &Formatter{
			TimestampFormat:        logFormat,
			DisableLevelTruncation: true,
		},
	}
	if debug {
		initializedLogger.Level = logrus.DebugLevel
	} else {
		initializedLogger.Level = logrus.InfoLevel
	}
	initializedLogger.SetReportCaller(true)
	return
}
