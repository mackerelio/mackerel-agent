package logging

import (
	"log"
)

type Logger struct {
	tag string
}

type LogLevel struct {
	name  string
	level uint
}

var trace = &LogLevel{name: "TRACE", level: 1}
var debug = &LogLevel{name: "DEBUG", level: 2}
var info = &LogLevel{name: "INFO", level: 3}
var warning = &LogLevel{name: "WARNING", level: 4}
var error = &LogLevel{name: "ERROR", level: 5}
var critical = &LogLevel{name: "CRITICAL", level: 6}

func stringToLogLevel(name string) *LogLevel {
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
	return &LogLevel{name: "unknown", level: 0}
}

var logLevelConfigs = map[string]*LogLevel{
	"root": info,
}

func GetLogger(tag string) *Logger {
	return &Logger{tag: tag}
}

func ConfigureLoggers(rootLogLevel string) {
	logLevelConfigs["root"] = stringToLogLevel(rootLogLevel)
}

func (logger *Logger) currentLogLevel() *LogLevel {
	return logLevelConfigs["root"]
}

func (logger *Logger) message(logLevel *LogLevel, message string) string {
	return logLevel.name + " " + logger.tag + " " + message
}

func (logger *Logger) log(logLevel *LogLevel, message string, args ...interface{}) {
	if logLevel.level >= logger.currentLogLevel().level {
		log.Printf(logger.message(logLevel, message), args...)
	}
}

func (logger *Logger) Criticalf(m string, args ...interface{}) {
	logger.log(critical, m, args...)
}

func (logger *Logger) Errorf(m string, args ...interface{}) {
	logger.log(error, m, args...)
}

func (logger *Logger) Warningf(m string, args ...interface{}) {
	logger.log(warning, m, args...)
}

func (logger *Logger) Infof(m string, args ...interface{}) {
	logger.log(info, m, args...)
}

func (logger *Logger) Debugf(m string, args ...interface{}) {
	logger.log(debug, m, args...)
}

func (logger *Logger) Tracef(m string, args ...interface{}) {
	logger.log(trace, m, args...)
}
