package main

import (
	"cfg"
	"fmt"
	"misc"
)

func main() {
	fmt.Println("Start centerserver ...")
	misc.SayHi()
	//cfg.AllIni()
	config := cfg.Get()
	//fmt.Printf("config.ini %v \n", config)
	fmt.Printf("Author : %v \n", config["author"])
	fmt.Printf("Start Date : %v \n", config["date"])
	fmt.Printf("Mail : %v \n", config["mail"])

	//	读取配置

	//	启动端口

	//	启动服务

	//	CS启动完成
	fmt.Println("centerserver is ready....")
}
