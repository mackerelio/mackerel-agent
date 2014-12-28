package logging

import (
	"log"
)

// Logger XXX
type Logger struct {
	tag string
}

type logLevel struct {
	name  string
	level uint
}

var trace = &logLevel{name: "TRACE", level: 1}
var debug = &logLevel{name: "DEBUG", level: 2}
var info = &logLevel{name: "INFO", level: 3}
var warning = &logLevel{name: "WARNING", level: 4}
var error = &logLevel{name: "ERROR", level: 5}
var critical = &logLevel{name: "CRITICAL", level: 6}

func stringTologLevel(name string) *logLevel {
	switch name {
	case "TRACE":
		return trace
	case "DEBUG":
		return debug
	case "INFO":
		return info
	case "WARNING":
		return warning
	case "CRITICAL":
		return critical
	}
	return &logLevel{name: "unknown", level: 0}
}

var logLevelConfigs = map[string]*logLevel{
	"root": info,
}

// GetLogger XXX
func GetLogger(tag string) *Logger {
	return &Logger{tag: tag}
}

// ConfigureLoggers XXX
func ConfigureLoggers(rootlogLevel string) {
	logLevelConfigs["root"] = stringTologLevel(rootlogLevel)
}

func (logger *Logger) currentlogLevel() *logLevel {
	return logLevelConfigs["root"]
}

func (logger *Logger) message(logLevel *logLevel, message string) string {
	return logLevel.name + " " + logger.tag + " " + message
}

func (logger *Logger) log(logLevel *logLevel, message string, args ...interface{}) {
	if logLevel.level >= logger.currentlogLevel().level {
		log.Printf(logger.message(logLevel, message), args...)
	}
}

// Criticalf XXX
func (logger *Logger) Criticalf(m string, args ...interface{}) {
	logger.log(critical, m, args...)
}

// Errorf XXX
func (logger *Logger) Errorf(m string, args ...interface{}) {
	logger.log(error, m, args...)
}

// Warningf XXX
func (logger *Logger) Warningf(m string, args ...interface{}) {
	logger.log(warning, m, args...)
}

// Infof XXX
func (logger *Logger) Infof(m string, args ...interface{}) {
	logger.log(info, m, args...)
}

// Debugf XXX
func (logger *Logger) Debugf(m string, args ...interface{}) {
	logger.log(debug, m, args...)
}

// Tracef XXX
func (logger *Logger) Tracef(m string, args ...interface{}) {
	logger.log(trace, m, args...)
}
