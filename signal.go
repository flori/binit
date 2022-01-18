package binit

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

type Signaller struct {
	signals *chan os.Signal
	cmd     *exec.Cmd
}

func NewSignaller() *Signaller {
	signaller := &Signaller{}
	return signaller
}

func (s *Signaller) Start(cmd *exec.Cmd) {
	s.cmd = cmd
	signals := make(chan os.Signal, 1)
	s.signals = &signals
	signal.Notify(*s.signals)
	go s.Forward()
}

func (s *Signaller) Forward() {
	for signal := range *s.signals {
		if signal == syscall.SIGCHLD {
			continue
		}
		pid := s.cmd.Process.Pid
		sig := signal.(syscall.Signal)
		log.Printf("forwarding signal %s (=%d) to process group of pid %d (%s) now",
			unix.SignalName(sig), unix.SignalNum(unix.SignalName(sig)), pid, s.cmd.Path)
		syscall.Kill(-pid, sig)
	}
}

func (s *Signaller) Shutdown() {
	close(*s.signals)
	signal.Reset()
}
