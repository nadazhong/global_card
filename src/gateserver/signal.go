package gateserver

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SignalProc() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			fmt.Println("Gate server quit.")
			quitChan <- 1
		}
	}
}
