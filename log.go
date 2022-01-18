package binit

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	// Switch logging to stdout as well
	return fmt.Print(time.Now().Format(time.RFC3339) + " " + string(bytes))
}

func ConfigureLogging() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}
