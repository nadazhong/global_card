package centerserver

import (
	"cfg"
	"fmt"
	"misc"
	"servernet"
)

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
		//case id, err := <-gstimer.GS_TIMER:
		//	if !err {
		//		cfg.LogErrf("Timer %v error %v.", id, err)
		//		return
		//	}
		//	centerTimerHandler(id, err)
		case session, err := <-centerServer.ConnectChan:
			if !err {
				fmt.Printf("Extra net new chan error:", err)
				return
			}
			fmt.Printf("Extra server new session:", session.Id)
			//newGs(session)
		case session, _ := <-centerServer.CloseChan:
			fmt.Printf("Extra server close sesion:", session.Id)
			//closeGs(session)
		case _, ok := <-centerServer.MsgQueue:
			if !ok {
				fmt.Printf("Extra net queue error.")
			}
			//if msg.Api == protocol.MAIL_CONNECTED_NTF
			//|| msg.Api == protocol.GS_CONNECTED_NTF {
			//	//conectHandler(msg)
			//	continue
			//}
			//if getMailSession() != nil && getMailSession().session.Id == msg.Session.Id {
			//	mailHandler(msg)
			//	continue
			//}
			//msgHander(msg)
		}
	}
}
