package gameserver

import (
	"cfg"
	"fmt"
	"misc"
	"strconv"
)

func init() {
	welcome()
}

var (
	GS_ID int16
)

func welcome() {
	fmt.Println("Start gameserver ...")
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

	if config["server_id"] == "" {
		fmt.Println("Server id not find.")
	}
	gsid, _ := strconv.Atoi(config["server_id"])
	GS_ID = int16(gsid)
	fmt.Println("配置的GS_ID ", GS_ID)

	if config["gate_intra_port"] == "" || config["gate_intra_ip"] == "" {
		fmt.Println("Gate addr config nil.")
	}

	fmt.Println("Starting the gameserver.")

	go SignalProc()

	gsDbInit()
	startWorker()
}

func gsDbInit() {
	//
	fmt.Println("载入gs数据")

}

// GS 逻辑启动
func gsRun() {
	fmt.Println("Game Server starting.")
}
