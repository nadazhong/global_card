package centerserver

import (
	//"cfg"
	"fmt"
	"gstimer"
	//"misc/packet"
	"os"
	"os/signal"
	//"protocol"
	"syscall"
)

var (
	STOP_SIG = make(chan int, 1)
)

func SignalProc() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			onExit(msg)
		case syscall.SIGTERM:
			onExit(msg)
		case syscall.SIGINT:
			onExit(msg)
		case syscall.SIGQUIT:
			onExit(msg)
		default:
			fmt.Println("signal=", msg)
		}
	}
}

func setQuitTimer() {
	gstimer.CreateTimer(0, 0, gstimer.Msg{Action: gstimer.ACTION_SHUTDOWN})
}

func onExit(msg os.Signal) {
	fmt.Printf("\033[043;1m[%s, quit]\033[0m", msg)
	setQuitTimer()
	<-STOP_SIG
	os.Exit(-1)
}

func centerQuitRoutine() {
	setQuitTimer()
	<-STOP_SIG
	//if getMailSession() != nil {
	//	getMailSession().session.SyncSend(packet.Pack(protocol.SERVER_QUIT_NTF,
	//		protocol.PKT_null{}, nil))
	//}
	os.Exit(0)
}
