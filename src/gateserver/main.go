package gateserver

import (
	"cfg"
	"fmt"
	"misc"
	"servernet"
)

var (
	quitChan chan int
)

// 多说一点 一个package可以有多个init 但是不推荐这么做
func init() {
	quitChan = make(chan int, 1)
}

func init() {
	welcome()
}

func welcome() {
	fmt.Println("Start gateserver ...")
	misc.SayHi()
	//cfg.AllIni()
	config := cfg.Get()
	//fmt.Printf("config.ini %v \n", config)
	fmt.Printf("Author : %v \n", config["author"])
	fmt.Printf("Start Date : %v \n", config["date"])
	fmt.Printf("Mail : %v \n", config["mail"])
}

func Start() {
	config := cfg.Get()

	if config["gate_port"] == "" || config["gate_ip"] == "" ||
		config["gate_intra_port"] == "" || config["gate_intra_ip"] == "" {
		fmt.Println("Gate port/ip or Gate intra_port/intra_ip config nil.")
	}

	addr := config["gate_ip"] + ":" + config["gate_port"]
	intraAddr := config["gate_intra_ip"] + ":" + config["gate_intra_port"]

	fmt.Println("Gate start listen(intra):", intraAddr)
	err, intraServer := servernet.NewServer(intraAddr, false) // 内部服务器
	if err != nil {
		fmt.Printf("Listen(intra) fail:", err)
	}

	fmt.Println("Gate start listen(extra):", addr)
	err, extraServer := servernet.NewServer(addr, true) // 外部服务器
	if err != nil {
		fmt.Printf("Listen(internet) fail:", err)
	}

	go SignalProc()

	gsGateInit()
	gateLoop(intraServer, extraServer)
}

func gsGateInit() {
	//initAccount()
	fmt.Println("读取Gate数据")
}

func gateLoop(intraServer *servernet.ClientInfo, extraServer *servernet.ClientInfo) {
	defer func() {
		fmt.Println("Gate server quit.")
	}()

	for {
		select {
		//case id, err := <-gstimer.GS_TIMER:
		//	if !err {
		//		cfg.LogErrf("Timer %v error %v.", id, err)
		//		return
		//	}
		//	gateTimerHandler(id, err)
		case session, err := <-intraServer.ConnectChan:
			if !err {
				fmt.Printf("Intra net new chan error:", err)
				return
			}
			fmt.Println("新的gs连接进来 ", session)
			fmt.Println("Intra server new session: ", session.Id)
			//newServerSession(session)
		case session, _ := <-intraServer.CloseChan:
			fmt.Printf("Intra session %v closed.", session.Id)
			//closeServerSession(session)
		case msg, ok := <-intraServer.MsgQueue:
			if !ok {
				fmt.Printf("Intra server msg err.")
				continue
			}
			fmt.Println("intraServer msg : ", msg)
			//gsHandler(msg)

		// 外部消息 负责处理跟客户端交互
		case session, err := <-extraServer.ConnectChan:
			if !err {
				fmt.Printf("Extra net new chan error:", err)
				return
			}
			fmt.Println("Extra server new session:", session.Id)
			//newClient(session)
		case session, _ := <-extraServer.CloseChan:
			fmt.Printf("Extra server close sesion:", session.Id)
			//closeClient(session)
		case msg, ok := <-extraServer.RawMsgQueue:
			if !ok {
				fmt.Printf("Extra net queue error.")
			}
			fmt.Println("extraServer msg : ", msg)
			//userMsgHander(msg)
		// account回调
		//case info, err := <-account.VerifyChan:
		//	if !err {
		//		cfg.LogFatal("account verify chan error:", err)
		//		return
		//	}
		//	info.Callback(info)
		// 停服 通知所有GS关闭游戏服务器
		case val, _ := <-quitChan:
			fmt.Printf("Server quit message, code %v.", val)
			//if len(gsSessions) > 0 {
			//	sendGsQuitMsg()
			//}
			return
		}
	}
}
