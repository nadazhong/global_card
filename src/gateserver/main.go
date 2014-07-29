package gateserver

import (
	"cfg"
	"fmt"
	"misc"
)

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

	//config := cfg.Get()

	//if config["cs_port"] == "" || config["cs_ip"] == "" {
	//	fmt.Println("CS config port/ip nil.")
	//}
	//addr := config["cs_ip"] + ":" + config["cs_port"]

	//if value := config["server_id"], value == ""{
	//	fmt.Println("Server id not find.")
	//}
	//gsid, _ := strconv.Atoi(value)
	//GS_ID = int16(gsid)

	//if config["intra_port"] == "" || config["intra_ip"] == "" {
	//	fmt.Println("center addr config nil.")
	//}

	//cfg.Log("Starting the game server.")

	//go SignalProc()

	//gsDbInit()
	//startWorker()
}
