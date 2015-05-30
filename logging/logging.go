package logging

import (
	"fmt"
	"log"
	"os"
)

// Logger struct for logging
type Logger struct {
	tag string
}

func stringToLevel(name string) level {
	switch name {
	case "TRACE":
		return TRACE
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "CRITICAL":
		return CRITICAL
	}
	return UNKNOWN
}

// GetLogger get the logger
func GetLogger(tag string) *Logger {
	return &Logger{tag: tag}
}

// global log level
var logLv = INFO
var lgr = log.New(os.Stderr, "", log.LstdFlags)

// ConfigureLoggers congigure log settings
func ConfigureLoggers(logLevel string) {
	logLv = stringToLevel(logLevel)
	if logLv <= DEBUG {
		lgr.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		lgr.SetFlags(log.LstdFlags)
	}
}

func (logger *Logger) message(lv level, message string) string {
   return lv.String() + " <" + logger.tag + "> " + message
}

func (logger *Logger) log(lv level, message string, args ...interface{}) {
	if lv >= logLv {
		const depth = 3
		lgr.Output(depth, fmt.Sprintf(logger.message(lv, message), args...))
	}
}

// Criticalf XXX
func (logger *Logger) Criticalf(m string, args ...interface{}) {
	logger.log(CRITICAL, m, args...)
}

// Errorf XXX
func (logger *Logger) Errorf(m string, args ...interface{}) {
	logger.log(ERROR, m, args...)
}

// Warningf XXX
func (logger *Logger) Warningf(m string, args ...interface{}) {
	logger.log(WARNING, m, args...)
}

// Infof XXX
func (logger *Logger) Infof(m string, args ...interface{}) {
	logger.log(INFO, m, args...)
}

// Debugf XXX
func (logger *Logger) Debugf(m string, args ...interface{}) {
	logger.log(DEBUG, m, args...)
}

// Tracef XXX
func (logger *Logger) Tracef(m string, args ...interface{}) {
	logger.log(TRACE, m, args...)
}
