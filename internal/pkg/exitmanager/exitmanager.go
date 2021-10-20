package exitmanager

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ExitMgr struct {
	_          struct{}
	exit       chan struct{}
	osSignals  chan os.Signal
	serverErrs chan error
	callbacks  []func()
}

var (
	SIGHUPSignal = 1
	SIGINTSignal = 2
)

func Init() *ExitMgr {
	// Create exit manager instance
	e := &ExitMgr{
		exit:       make(chan struct{}),
		osSignals:  make(chan os.Signal, SIGINTSignal),
		serverErrs: make(chan error, SIGHUPSignal),
	}

	// Listen for exit signals
	go e.listen()

	return e
}

func (e *ExitMgr) listen() {
	// Defer shutdown
	defer e.shutdown()

	// Listen for os signals
	signal.Notify(e.osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle exit cases
	select {
	case servererr := <-e.serverErrs:
		log.Println(servererr)
	case sig := <-e.osSignals:
		s := fmt.Sprintf("signal: %s", sig.String())
		log.Println(s)
	}
}

func (e *ExitMgr) shutdown() {
	log.Println("Shutting down uffizzi controller")

	// Iterate through tear down callbacks
	for _, f := range e.callbacks {
		f()
	}

	// Close exit channel to unblock main go routine
	close(e.exit)
}

// Wait blocks the caller until the exit channel is closed
func (e *ExitMgr) Wait() {
	<-e.exit
}

// ServerError handles server errors
func (e *ExitMgr) ServerError(err error) {
	e.serverErrs <- err
}

// AddCallback adds a teardown callback
func (e *ExitMgr) AddTeardownCallback(f func()) {
	e.callbacks = append(e.callbacks, f)
}
