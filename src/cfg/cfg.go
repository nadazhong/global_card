package cfg

import (
	"fmt"
	"github.com/cfg"
	"log"
)

//	查看所有配置文件
func AllIni() {
	myIni := make(map[string]string)
	err := cfg.Load("./config.ini", myIni)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", myIni)
}

// 获取AllIini配置
func Get() map[string]string {
	myIni := make(map[string]string)
	err := cfg.Load("./config.ini", myIni)
	if err != nil {
		log.Fatal(err)
	}
	return myIni
}
