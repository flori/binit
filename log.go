package binit

import (
	"fmt"
	"log"
	"time"
)

const PROGRAM_NAME = "binit"

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	// Switch logging to stdout as well
	return fmt.Printf("%s %s: %s", time.Now().Format(time.RFC3339), PROGRAM_NAME, string(bytes))
}

func ConfigureLogging() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}
