package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

const name = "mackerel-agent"

const (
	defaultEid = iota + 1
	startEid
	stopEid
	loggerEid
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32")
	procAllocConsole             = kernel32.NewProc("AllocConsole")
	procGenerateConsoleCtrlEvent = kernel32.NewProc("GenerateConsoleCtrlEvent")
	procGetModuleFileName        = kernel32.NewProc("GetModuleFileNameW")
)

func main() {
	if len(os.Args) == 2 {
		var err error
		switch os.Args[1] {
		case "install":
			err = installService("mackerel-agent", "mackerel agent")
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

type logger interface {
	Info(eid uint32, msg string) error
	Warning(eid uint32, msg string) error
	Error(eid uint32, msg string) error
}

type handler struct {
	elog logger
	cmd  *exec.Cmd
	r    io.Reader
	w    io.WriteCloser
	wg   sync.WaitGroup
}

// ex.
// verbose log: 2017/01/21 22:21:08 command.go:434: DEBUG <command> received 'immediate' chan
// normal log:  2017/01/24 14:14:27 INFO <main> Starting mackerel-agent version:0.36.0
var logRe = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} (?:\S+\.go:\d+: )?([A-Z]+) `)

func (h *handler) start(logDone chan struct{}) error {
	procAllocConsole.Call()
	dir := execdir()
	cmd := exec.Command(filepath.Join(dir, "mackerel-agent.exe"))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	cmd.Dir = dir

	h.cmd = cmd
	h.r, h.w = io.Pipe()
	cmd.Stderr = h.w

	err := h.cmd.Start()
	if err != nil {
		return err
	}

	return h.aggregate(logDone)
}

func (h *handler) aggregate(logDone chan struct{}) error {
	br := bufio.NewReader(h.r)
	lc := make(chan string, 10)
	done := make(chan struct{})

	// It need to read data from pipe continuously. And it need to close handle
	// when process finished.
	// read data from pipe. When data arrived at EOL, send line-string to the
	// channel. Also remaining line to EOF.
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		defer h.w.Close()

		// pipe stderr to windows event log
		var body bytes.Buffer
		for {
			b, err := br.ReadByte()
			if err != nil {
				if err != io.EOF {
					h.elog.Error(loggerEid, err.Error())
				}
				break
			}
			if b == '\n' {
				if body.Len() > 0 {
					lc <- body.String()
					body.Reset()
				}
				continue
			}
			body.WriteByte(b)
		}
		if body.Len() > 0 {
			lc <- body.String()
		}
		done <- struct{}{}
	}()

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		defer close(logDone)

		linebuf := []string{}
	loop:
		for {
			select {
			case line := <-lc:
				if len(linebuf) == 0 || logRe.MatchString(line) {
					linebuf = append(linebuf, line)
				} else {
					linebuf[len(linebuf)-1] += "\n" + line
				}
			case <-time.After(10 * time.Millisecond):
				// When it take 10ms, it is located at end of paragraph. Then
				// slice appended at above should be the paragraph.
				for _, line := range linebuf {
					if match := logRe.FindStringSubmatch(line); match != nil {
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
				select {
				case <-done:
					break loop
				default:
				}
				linebuf = nil
			}
		}
		close(lc)
		close(done)
	}()

	return nil
}

func interrupt(p *os.Process) error {
	r1, _, err := procGenerateConsoleCtrlEvent.Call(syscall.CTRL_BREAK_EVENT, uintptr(p.Pid))
	if r1 == 0 {
		return err
	}
	return nil
}

func (h *handler) stop() error {
	if h.cmd != nil && h.cmd.Process != nil {
		err := interrupt(h.cmd.Process)
		if err == nil {
			end := time.Now().Add(10 * time.Second)
			for time.Now().Before(end) {
				if h.cmd.ProcessState != nil && h.cmd.ProcessState.Exited() {
					return nil
				}
				time.Sleep(1 * time.Second)
			}
		}
		return h.cmd.Process.Kill()
	}

	h.wg.Wait()
	return nil
}

// implement https://godoc.org/golang.org/x/sys/windows/svc#Handler
func (h *handler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	s <- svc.Status{State: svc.StartPending}
	defer func() {
		s <- svc.Status{State: svc.Stopped}
	}()

	logDone := make(chan struct{})
	if err := h.start(logDone); err != nil {
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
	<-logDone

	return
}

func execdir() string {
	var wpath [syscall.MAX_PATH]uint16
	r1, _, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&wpath[0])), uintptr(len(wpath)))
	if r1 == 0 {
		log.Fatal(err)
	}
	return filepath.Dir(syscall.UTF16ToString(wpath[:]))
}
