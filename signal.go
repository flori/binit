package binit

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

type Signaller struct {
	signals       *chan os.Signal
	skipSignalLog []string
	cmd           *exec.Cmd
}

func NewSignaller(config Config) *Signaller {
	signaller := &Signaller{
		skipSignalLog: strings.Split(config.SKIP_SIGNAL_LOG, ","),
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
	for _, s := range s.skipSignalLog {
		if s == name {
			return false
		}
	}
	return true
}

func (s *Signaller) Forward() {
	for signal := range *s.signals {
		if signal == syscall.SIGCHLD {
			continue
		}
		pid := s.cmd.Process.Pid
		sig := signal.(syscall.Signal)
		sigName := unix.SignalName(sig)
		if s.logSignal(sigName) {
			log.Printf("forwarding signal %s (=%d) to process group of pid %d (%s) now",
				sigName, sig, pid, s.cmd.Path)
		}
		syscall.Kill(-pid, sig)
	}
}

func (s *Signaller) Shutdown() {
	close(*s.signals)
	signal.Reset()
}
