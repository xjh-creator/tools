package cron

import (
	"io/ioutil"
	"log"
	"os"
	"time"
	"tools/glog"
)

// DefaultLogger is used by Cron if none is specified.
var DefaultLogger Logger = PrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))

// DiscardLogger can be used by callers to discard all log messages.
var DiscardLogger Logger = PrintfLogger(log.New(ioutil.Discard, "", 0))

// Logger is the interface used in this package for logging, so that any backend
// can be plugged in. It is a subset of the github.com/go-logr/logr interface.
type Logger interface {
	// Info logs routine messages about cron's operation.
	Info(msg string, keysAndValues ...interface{})
	// Error logs an error condition.
	Error(err error, msg string, keysAndValues ...interface{})
}

// PrintfLogger wraps a Printf-based logger (such as the standard library "log")
// into an implementation of the Logger interface which logs errors only.
func PrintfLogger(l interface{ Printf(string, ...interface{}) }) Logger {
	return printfLogger{l, false}
}

// VerbosePrintfLogger wraps a Printf-based logger (such as the standard library
// "log") into an implementation of the Logger interface which logs everything.
func VerbosePrintfLogger(l interface{ Printf(string, ...interface{}) }) Logger {
	return printfLogger{l, true}
}

type printfLogger struct {
	logger  interface{ Printf(string, ...interface{}) }
	logInfo bool
}

func (pl printfLogger) Info(msg string, keysAndValues ...interface{}) {
	if glog.Glog().Level >= glog.InfoLevel {
		keysAndValues = formatTimes(keysAndValues)
		glog.Infof("[任务调度] %s"+formatString(len(keysAndValues)), append([]interface{}{msg}, keysAndValues...)...)
	}
}

func (pl printfLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	glog.Infof("[任务调度] %s"+formatString(len(keysAndValues)), append([]interface{}{msg}, keysAndValues...)...)
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	str := ""
	if n := numKeysAndValues; n > 0 {
		for i := 0; i < n/2; i++ {
			str += ", %v:%v"
		}
	}
	return str
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
