package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Options is a function type that can be used to configure the logger
type Options func(*LogWrapper)

// WithLevel configures the log level. If level is not specified, default to InfoLevel
// If level is debug or trace, report caller is enabled
func WithLevel(level string) Options {
	return func(lw *LogWrapper) {
		l, err := logrus.ParseLevel(level)
		if err != nil {
			lw.entry.Logger.SetLevel(logrus.InfoLevel)
		} else {
			lw.entry.Logger.SetLevel(l)
			if l == logrus.DebugLevel || l == logrus.TraceLevel {
				formatter := &logrus.TextFormatter{
					TimestampFormat:        time.RFC3339,
					FullTimestamp:          true,
					DisableLevelTruncation: true,
					ForceColors:            true,
					PadLevelText:           false,
					DisableColors:          false,
					/*
						CallerPrettyfier: func(f *runtime.Frame) (string, string) {
							var files string
							var funcs string
							for i := 1; i < 12; i++ {
								if pc, file, line, ok := runtime.Caller(i); ok {
									fName := runtime.FuncForPC(pc).Name()
									if strings.Contains(fName, "logrus") || strings.Contains(fName, "log") {
										continue
									} else {
										funcs += fmt.Sprintf("func[%d]: %s - ", i, formatFilePath(fName, 1))
										files += fmt.Sprintf("file[%d]: %s:%d - ", i, formatFilePath(file, 2), line)
									}
								}
							}
							funcs += fmt.Sprintf("funcX: %s - ", formatFilePath(f.Function, 1))
							files += fmt.Sprintf("fileX: %s:%d - ", formatFilePath(f.File, 2), f.Line)
							return fmt.Sprintf("func: %s - ", funcs), fmt.Sprintf(" src: %s -", files)
						},
					*/
					CallerPrettyfier: func(f *runtime.Frame) (string, string) {
						if pc, file, line, ok := runtime.Caller(10); ok {
							fName := runtime.FuncForPC(pc).Name()
							return fmt.Sprintf("func: %s : ", formatFilePath(fName, 1)), fmt.Sprintf(" src: %s:%d -", formatFilePath(file, 2), line)
						}
						return fmt.Sprintf("func: %s : ", formatFilePath(f.Function, 1)), fmt.Sprintf(" src: %s:%d -", formatFilePath(f.File, 2), f.Line)
					},
				}

				lw.entry.Logger.SetFormatter(formatter)
				lw.entry.Logger.SetReportCaller(true)
			}
		}
	}
}

// WithOutput configures the output destination
func WithOutput(output io.Writer) Options {
	return func(lw *LogWrapper) {
		lw.entry.Logger.SetOutput(output)
	}
}

// WithFormatter configures the log formatter
func WithFormatter(formatter logrus.Formatter) Options {
	return func(lw *LogWrapper) {
		lw.entry.Logger.SetFormatter(formatter)
	}
}

// WithReportCaller configures the log to report caller
func WithReportCaller(reportCaller bool) Options {
	return func(lw *LogWrapper) {
		lw.entry.Logger.SetReportCaller(reportCaller)
	}
}

// WithField adds a field to the log entry
// func WithField(key string, value interface{}) Options {
// 	return func(lw *LogWrapper) {
// 		lw.entry.WithField(key, value)
// 	}
// }

// WithFields adds multiple fields to the log entry
// func WithFields(fields logrus.Fields) Options {
// 	return func(lw *LogWrapper) {
// 		lw.entry.WithFields(fields)
// 	}
// }

// WithNullLogger sets the logger to discard all output
func WithNullLogger() Options {
	return func(lw *LogWrapper) {
		lw.entry.Logger.SetOutput(io.Discard)
	}
}

// formatFilePath receives a string representing a path and returns the last part of it
// The 2nd argument indicates the number of parts to return
func formatFilePath(path string, parts int) string {
	arr := strings.Split(path, "/")
	return strings.Join(arr[len(arr)-parts:], "/")
}
