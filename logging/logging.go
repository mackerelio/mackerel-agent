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

// GetLogger get the logger
func GetLogger(tag string) *Logger {
	return &Logger{tag: tag}
}

// global log level
var logLv = INFO
var lgr = log.New(os.Stderr, "", log.LstdFlags)

// SetLogLevel congigure log settings
func SetLogLevel(lv level) {
	if logLv != lv {
		logLv = lv
		if logLv <= DEBUG {
			lgr.SetFlags(log.LstdFlags | log.Lshortfile)
		} else {
			lgr.SetFlags(log.LstdFlags)
		}
	}
}

func (logger *Logger) message(lv level, message string) string {
	return lv.String() + " <" + logger.tag + "> " + message
}

func (logger *Logger) log(lv level, message string, args ...interface{}) {
	if lv >= logLv {
		// caller -> Infof() -> log()
		const depth = 3
		lgr.Output(depth, fmt.Sprintf(logger.message(lv, message), args...))
	}
}

// Criticalf critical log
func (logger *Logger) Criticalf(m string, args ...interface{}) {
	logger.log(CRITICAL, m, args...)
}

// Errorf error log
func (logger *Logger) Errorf(m string, args ...interface{}) {
	logger.log(ERROR, m, args...)
}

// Warningf warning log
func (logger *Logger) Warningf(m string, args ...interface{}) {
	logger.log(WARNING, m, args...)
}

// Infof info log
func (logger *Logger) Infof(m string, args ...interface{}) {
	logger.log(INFO, m, args...)
}

// Debugf debug log
func (logger *Logger) Debugf(m string, args ...interface{}) {
	logger.log(DEBUG, m, args...)
}

// Tracef trace log for debugging details
func (logger *Logger) Tracef(m string, args ...interface{}) {
	logger.log(TRACE, m, args...)
}
