package centerserver

import (
	"cfg"
	"fmt"
	"gstimer"
	"misc"
	"protocol"
	"servernet"
)

var (
	centerQuit bool // 停服标记。
)

func init() {
	centerQuit = false
}

func init() {
	welcome()
}

func welcome() {
	fmt.Println("Start centerserver ...")
	misc.SayHi()
	//cfg.AllIni()
	config := cfg.Get()
	//fmt.Printf("config.ini %v \n", config)
	fmt.Printf("Author : %v \n", config["author"])
	fmt.Printf("Start Date : %v \n", config["date"])
	fmt.Printf("Mail : %v \n", config["mail"])
}

// 启动server入口
func Start() {
	config := cfg.Get()
	if config["cs_port"] == "" || config["cs_ip"] == "" {
		fmt.Println("CS config port/ip nil.")
	}
	addr := config["cs_ip"] + ":" + config["cs_port"]

	fmt.Println("CS listen: ", addr)
	err, centerServer := servernet.NewServer(addr, false)
	if err != nil {
		fmt.Println("Listen fail:", err)
	}

	go SignalProc()
	//	载入center的数据
	centerDbInit()

	CenterLoop(centerServer)
}

func centerDbInit() {
	fmt.Println("载入数据中心数据...")
}

func CenterLoop(centerServer *servernet.ClientInfo) {
	defer func() {
		fmt.Println("Center server quit.")
		//helper.BackTrace("worker")
	}()

	for {
		select {
		case id, err := <-gstimer.GS_TIMER:
			if !err {
				fmt.Printf("Timer %v error %v.", id, err)
				return
			}
			centerTimerHandler(id, err)
		case session, err := <-centerServer.ConnectChan:
			if !err {
				fmt.Printf("Extra net new chan error:", err)
				return
			}
			fmt.Println("新的GS连接进来 ", session)
			fmt.Printf("Extra server new session:", session.Id)
			//newGs(session)
		case session, _ := <-centerServer.CloseChan:
			fmt.Printf("Extra server close sesion:", session.Id)
			//closeGs(session)
		case msg, ok := <-centerServer.MsgQueue:
			if !ok {
				fmt.Println("Extra net queue error.")
			}
			fmt.Println("center msg queue recv:", msg.Api)
			// 如果消息是GS->Center的连接消息
			if msg.Api == protocol.GS_CONNECTED_NTF {
				conectHandler(msg)
				continue
			}
			// 处理其他消息
			msgHander(msg)
		}
	}
}

func doSave() {
	fmt.Println("center 保存数据")
}

func shutDown() {
	fmt.Println("收到关服消息.")
	doSave()
	// 关闭db连接
	//gsdb.CloseCenterDB()
}

func centerTimerHandler(id int64, ok bool) {
	if !ok {
		fmt.Println("GSTimer error: %v", ok)
	}
	timer := gstimer.TIMER_MAP[id]
	switch timer.Msg.Action {
	case gstimer.ACTION_SHUTDOWN:
		shutDown()
		STOP_SIG <- 1
		return
	}
	gstimer.FreeTimer(id)
}
