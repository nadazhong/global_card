package netcon

import (
	"cfg"
	"fmt"
	"misc"
	"servernet"
)

const (
	ELE_TIME = 1000
)

var (
	connCtime int64
	connGtime int64
)

// 连接Gate。
func connect(intraAddr string) *servernet.ConnInfo {
	fmt.Printf("Conn server %v. \n", intraAddr)
	err, conn := servernet.ConnServer(intraAddr)
	if err == nil {
		return conn
	}
	fmt.Printf("Conn server %v fail \n", intraAddr)

	return nil
}

func ConnGate() *servernet.ConnInfo {

	now := misc.Current()
	if now < (connGtime + ELE_TIME) {
		return nil
	}
	connGtime = now
	config := cfg.Get()
	addr := config["gate_intra_ip"] + ":" + config["gate_intra_port"]
	conn := connect(addr)
	if conn == nil {
		fmt.Println("Connect gate fail:", addr)
		return nil
	}
	fmt.Printf("Connect gate %s success: \n", addr)
	return conn

}

func ConnCenter() *servernet.ConnInfo {

	now := misc.Current()
	if now < (connCtime + ELE_TIME) {
		return nil
	}
	connCtime = now
	config := cfg.Get()
	addr := config["cs_ip"] + ":" + config["cs_port"]
	conn := connect(addr)
	if conn == nil {
		fmt.Println("Connect center fail:", addr)
		return nil
	}
	fmt.Printf("Connect center %s success \n", addr)
	return conn
}
