package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

const name = "mackerel-agent"

const defaultEid = 1
const startEid = 2
const stopEid = 3
const loggerEid = 4

func main() {
	if len(os.Args) == 2 {
		var err error
		switch os.Args[1] {
		case "install":
			err = installService("mackerel-agent", "mackarel agent")
		case "remove":
			err = removeService("mackerel-agent")
		}
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	elog, err := eventlog.Open(name)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer elog.Close()

	// `svc.Run` blocks until windows service will stopped.
	// ref. https://msdn.microsoft.com/library/cc429362.aspx
	err = svc.Run(name, &handler{elog: elog})
	if err != nil {
		log.Fatal(err.Error())
	}
}

type handler struct {
	elog *eventlog.Log
	cmd  *exec.Cmd
}

func (h *handler) start() error {
	cmd := exec.Command(filepath.Join(filepath.Dir(execdir()), "mackerel-agent.exe"))
	h.cmd = cmd
	r, w := io.Pipe()
	cmd.Stderr = w
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines) // default
	go func() {
		// pipe stderr to windows event log
		re := regexp.MustCompile("^\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} (\\w+) ")
		for scanner.Scan() {
			line := scanner.Text()
			if match := re.FindStringSubmatch(line); match != nil {
				level := match[1]
				switch level {
				case "TRACE", "DEBUG", "INFO":
					h.elog.Info(defaultEid, line)
				case "WARNING":
					h.elog.Warning(defaultEid, line)
				case "ERROR", "CRITICAL":
					h.elog.Error(defaultEid, line)
				default:
					h.elog.Error(defaultEid, line)
				}
			} else {
				h.elog.Error(defaultEid, line)
			}
		}
		if err := scanner.Err(); err != nil {
			h.elog.Error(loggerEid, err.Error())
		} else {
			// EOF
		}
	}()
	return cmd.Start()
}

func (h *handler) stop() error {
	if h.cmd != nil && h.cmd.Process != nil {
		return h.cmd.Process.Kill()
	}
	return nil
}

// implement https://godoc.org/golang.org/x/sys/windows/svc#Handler
func (h *handler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	s <- svc.Status{State: svc.StartPending}
	defer func() {
		s <- svc.Status{State: svc.Stopped}
	}()

	if err := h.start(); err != nil {
		h.elog.Error(startEid, err.Error())
		// https://msdn.microsoft.com/library/windows/desktop/ms681383(v=vs.85).aspx
		// use ERROR_SERVICE_SPECIFIC_ERROR
		return true, 1
	}

	exit := make(chan struct{})
	go func() {
		err := h.cmd.Wait()
		// enter when the child process exited
		if err != nil {
			h.elog.Error(stopEid, err.Error())
		}
		exit <- struct{}{}
	}()

	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
L:
	for {
		select {
		case req := <-r:
			switch req.Cmd {
			case svc.Interrogate:
				s <- req.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending, Accepts: svc.AcceptStop | svc.AcceptShutdown}
				if err := h.stop(); err != nil {
					h.elog.Error(stopEid, err.Error())
					s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
				}
			}
		case <-exit:
			break L
		}
	}

	return
}

func execdir() string {
	var (
		kernel32              = syscall.NewLazyDLL("kernel32")
		procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")
	)
	var wpath [syscall.MAX_PATH]uint16
	r1, _, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&wpath[0])), uintptr(len(wpath)))
	if r1 == 0 {
		log.Fatal(err)
	}
	return syscall.UTF16ToString(wpath[:])
}
