package binit

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

type Signaller struct {
	signals       *chan os.Signal
	skipSignalLog []string
	preStop       string
	preStopSignal []string
	cmd           *exec.Cmd
}

func NewSignaller(config Config) *Signaller {
	signaller := &Signaller{
		skipSignalLog: strings.Split(config.SKIP_SIGNAL_LOG, ","),
		preStop:       config.PRE_STOP,
		preStopSignal: strings.Split(config.PRE_STOP_SIGNAL, ","),
	}
	return signaller
}

func (s *Signaller) Start(cmd *exec.Cmd) {
	s.cmd = cmd
	signals := make(chan os.Signal, 1)
	s.signals = &signals
	signal.Notify(*s.signals)
	go s.Forward()
}

func (s *Signaller) logSignal(name string) bool {
	if s == nil {
		return true
	}
	for _, s := range s.skipSignalLog {
		if s == name {
			return false
		}
	}
	return true
}

func (s *Signaller) isPreStopSignal(name string) bool {
	if s == nil {
		return false
	}
	for _, s := range s.preStopSignal {
		if s == name {
			return true
		}
	}
	return false
}

func (s *Signaller) Forward() {
	for signal := range *s.signals {
		if signal == syscall.SIGCHLD {
			continue
		}
		process := s.cmd.Process
		pid := -1
		if process != nil {
			pid = process.Pid
		}
		sig := signal.(syscall.Signal)
		sigName := unix.SignalName(sig)

		preStopping := s.isPreStopSignal(sigName)
		if preStopping {
			log.Printf(
				"received PRE_STOP_SIGNAL %s (=%d) for supervised process group of pid %d (%s) now",
				sigName, sig, pid, s.cmd.Path,
			)
			if s.preStop != "" {
				log.Printf("Executing PRE_STOP command for pid %d (%s) now", pid, s.cmd.Path)
				waiter := NewWaiter()
				env := map[string]string{}
				if pid > 0 {
					env["BINIT_CMD_PID"] = strconv.Itoa(pid)
				}
				runCommandWithEnv(s.preStop, waiter, &env)
			} else {
				log.Printf("PRE_STOP command is empty, so nothing will be executed")
			}
		} else {
			if s.logSignal(sigName) {
				log.Printf(
					"Forwarding signal %s (=%d) to process group of pid %d (%s) now",
					sigName, sig, pid, s.cmd.Path,
				)
			}
		}

		syscall.Kill(-pid, sig)
	}
}

func (s *Signaller) Shutdown() {
	close(*s.signals)
	signal.Reset()
}
