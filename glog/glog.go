package glog

import (
	"github.com/kataras/golog"
	"path/filepath"
	"runtime"
)

var(
	/*
	5 disable
	4 fatal
	3 error
	2 warn
	1 info
	0 debug
	*/
	log *golog.Logger
)

const(
	DebugLevel golog.Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	DisableLevel
)

func Init(g *golog.Logger)  {
	log = g
}

func Glog() *golog.Logger {
	return log
}

func init()  {
	if log == nil{
		log = golog.New()
		log.SetTimeFormat("2006/01/02 15:04:05.000")
	}
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Log(level golog.Level, v ...interface{}) {
	log.Log(level, v...)
}

func Logf(level golog.Level, format string, args ...interface{}) {
	log.Logf(level, format, args...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Error(v ...interface{}) {
	log.Error(v...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Warn(v ...interface{}) {
	log.Warn(v...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func Info(v ...interface{}) {
	log.Info(v...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func GetTest() {
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	Println()
	Println()
	for {
		f, more := frames.Next()
		Println(filepath.Base(f.File), f.Line)
		if !more || true {
			break
		}
	}
	Println()
	Println()
}


