package logging

import (
	"testing"
	"time"
)

func TestGetLogger(t *testing.T) {
	var logger = GetLogger("tag")
	if logger.tag != "tag" {
		t.Errorf("tag should be tag but %v", logger.tag)
	}
}

func TestConfigureLoggers(t *testing.T) {
	ConfigureLoggers("INFO")
	if logLevelConfigs["root"] != info {
		t.Errorf("tag should be tag but %v", logLevelConfigs["root"])
	}
}

// These tests do not check anything yet.
// You can see result by go test -v
func TestLogf(t *testing.T) {
	ConfigureLoggers("TRACE")

	var logger = GetLogger("tag")
	logger.Criticalf("This is critical log: %v", time.Now())
	logger.Errorf("This is error log: %v", time.Now())
	logger.Warningf("This is warning log: %v", time.Now())
	logger.Infof("This is info log: %v", time.Now())
	logger.Debugf("This is debug log: %v", time.Now())
	logger.Tracef("This is trace log: %v", time.Now()) // Shown until here
}

func TestInfoLogf(t *testing.T) {
	ConfigureLoggers("INFO")

	var logger = GetLogger("tag")
	logger.Criticalf("This is critical log: %v", time.Now())
	logger.Errorf("This is error log: %v", time.Now())
	logger.Warningf("This is warning log: %v", time.Now())
	logger.Infof("This is info log: %v", time.Now()) // Shown until here
	logger.Debugf("This is debug log: %v", time.Now())
	logger.Tracef("This is trace log: %v", time.Now())
}

func TestCriticalLogf(t *testing.T) {
	ConfigureLoggers("CRITICAL")

	var logger = GetLogger("tag")
	logger.Criticalf("This is critical log: %v", time.Now()) // Shown until here
	logger.Errorf("This is error log: %v", time.Now())
	logger.Warningf("This is warning log: %v", time.Now())
	logger.Infof("This is info log: %v", time.Now())
	logger.Debugf("This is debug log: %v", time.Now())
	logger.Tracef("This is trace log: %v", time.Now())
}
