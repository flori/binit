package binit

import (
	"context"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

type Waiter struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWaiter() *Waiter {
	waiter := Waiter{}
	waiter.init()
	return &waiter
}

func (w *Waiter) init() {
	ctx, cancel := context.WithCancel(context.Background())
	w.ctx = ctx
	w.cancel = cancel
	w.wg.Add(1)
	go w.Wait()
}

func (w *Waiter) Wait() {
	for {
		var status syscall.WaitStatus

		pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil) // wait3
		if err != nil && err != syscall.ECHILD {
			log.Printf("error while waiting: %v", err)
		}

		if pid > 0 {
			continue
		}

		time.Sleep(1 * time.Second)

		select {
		case <-w.ctx.Done():
			w.wg.Done()
			return
		default:
		}
	}
}

func (w *Waiter) Fatalf(message string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf(message, args...)
	} else {
		log.Print(message)
	}
	w.Quit(1)
}

func (w *Waiter) Quit(code int) {
	w.cancel()
	w.wg.Wait()

	os.Exit(code)
}
