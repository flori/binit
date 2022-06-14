package binit

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kballard/go-shellquote"
)

func runCommandWithEnv(command string, waiter *Waiter, env *map[string]string) {
	if command == "" {
		return
	}
	args, err := shellquote.Split(command)
	if err != nil {
		waiter.Fatalf("error when splitting %s: %v", command, err)
	}
	cmd := exec.Command(args[0], args[1:]...)

	if env != nil {
		cmd.Env = os.Environ()
		for name, value := range *env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", name, value))
		}
	}

	log.Printf("Running command { %v } now.", shellquote.Join(cmd.Args...))

	stat, err := os.Stdin.Stat()
	if err != nil {
		waiter.Fatalf("error while stating stdin: %v", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		waiter.Fatalf("error while running %s: %v", command, err)
	}
}
