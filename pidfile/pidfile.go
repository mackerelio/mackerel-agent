package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("pidfile")

// Create pidfile
func Create(pidfile string) error {
	if pidfile == "" {
		return nil
	}
	if pidString, err := ioutil.ReadFile(pidfile); err == nil {
		if pid, err := strconv.Atoi(string(pidString)); err == nil {
			if pid == os.Getpid() {
				return nil
			}
			if GetCmdName(pid) == filepath.Base(os.Args[0]) {
				return fmt.Errorf("pidfile found, try stopping another running mackerel-agent or delete %s", pidfile)
			}
			// Note mackerel-agent in windows can't remove pidfile during stoping the service
			logger.Warningf("Pidfile found, but there seems no another process of mackerel-agent. Ignoring %s", pidfile)
		} else {
			logger.Warningf("Malformed pidfile found. Ignoring %s", pidfile)
		}
	}

	err := os.MkdirAll(filepath.Dir(pidfile), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(pidfile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", os.Getpid())
	return err
}

// Remove pidfile
func Remove(pidfile string) error {
	if pidfile == "" {
		return nil
	}
	err := os.Remove(pidfile)
	if err != nil {
		logger.Errorf("Failed to remove the pidfile: %s: %s", pidfile, err)
	}
	return err
}
