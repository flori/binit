package binit

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/kballard/go-shellquote"
)

func Run(config Config, args []string) {
	waiter := NewWaiter()

	if len(args) == 0 {
		waiter.Fatalf("error, need a program to execute")
	}

	cmd := exec.Command(args[0], args[1:]...)
	stat, err := os.Stdin.Stat()
	if err != nil {
		waiter.Fatalf("error while stating stdin: %v", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if config.WORKDIR != "" {
		if err := os.Chdir(config.WORKDIR); err != nil {
			waiter.Fatalf("error changing to workdir %s for %s: %v", config.WORKDIR, cmd.Path, err)
		}
	}

	signaller := NewSignaller()
	signaller.Start(cmd)
	defer signaller.Shutdown()

	log.Printf("Running command { %v } now.", shellquote.Join(cmd.Args...))
	if err := cmd.Start(); err != nil {
		waiter.Fatalf("error while running %s: %v", cmd.Path, err)
	}

	if err := cmd.Wait(); err != nil {
		waiter.Fatalf("error while waiting for %s: %v", cmd.Path, err)
	}

	waiter.Quit(0)
}
