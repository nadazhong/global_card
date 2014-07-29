package servernet

import (
	"fmt"
	"net"
)

type ConnInfo struct {
	conn    *net.TCPConn
	sender  *Sender
	kickout bool

	CloseChan chan int
	MsgQueue  chan *NetMsg
}

func ConnServer(addr string) (error, *ConnInfo) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err, nil
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err, nil
	}

	info := ConnInfo{
		conn:    conn,
		sender:  newSender(conn),
		kickout: false,

		CloseChan: make(chan int, 5),
		MsgQueue:  make(chan *NetMsg, NETMSG_QUEUE_SIZE),
	}

	go info.start()

	return nil, &info
}

func (serverInfo *ConnInfo) Send(data []byte) (err error) {
	return serverInfo.sender.Send(data)
}

func (serverInfo *ConnInfo) SyncSend(data []byte) (err error) {
	return serverInfo.sender.SyncSend(data)
}

func (serverInfo *ConnInfo) start() {
	defer func() {
		fmt.Printf("Connection to server %p quit.", serverInfo)
		serverInfo.conn.Close()
		serverInfo.sender.Close()
		serverInfo.CloseChan <- 1
	}()

	conn := serverInfo.conn
	fmt.Printf("Connection %p start.", serverInfo)

	go serverInfo.sender.start(nil)

	header := make([]byte, 2)
	for {
		if serverInfo.kickout {
			return
		}

		data := readMsg(conn, header)
		if data == nil {
			return
		}

		msg := parseMsg(data, nil) // 作为客户端连接的消息，Session为nil。
		if msg == nil {
			return
		}
		serverInfo.MsgQueue <- msg
	}
}
