package gameserver

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

var (
	SIGTERM  int32
	SIGINT   int32
	STOP_SIG = make(chan int, 1)
)

func setQuitTimer() {
	return
	//gstimer.CreateTimer(0, 0, gstimer.Msg{Action: gstimer.ACTION_SHUTDOWN})
}

// 进程信号处理，主要捕获异常，退出之前做一些持久化处理。
func SignalProc() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			//cfg.Log("\033[043;1m[SIGHUP]\033[0m")
			//cfg.Reload()
		case syscall.SIGTERM:
			atomic.StoreInt32(&SIGTERM, 1)
			sigQuit(msg)
		case syscall.SIGINT:
			atomic.StoreInt32(&SIGINT, 1)
			sigQuit(msg)
		default:
			//cfg.Log("signal=", msg)
		}
	}
}

func sigQuit(msg os.Signal) {
	fmt.Printf("\033[043;1m[%v, quit]\033[0m", msg)
	//if gsInited {
	//	setQuitTimer()
	//	<-STOP_SIG
	//}
	os.Exit(-1)
}
