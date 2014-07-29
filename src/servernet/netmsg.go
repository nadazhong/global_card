package servernet

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"protocol"
	"time"
)

type NetMsg struct {
	Api     int16
	Payload interface{}
	Session *Session // 如果为作为客户端连接的受到的消息，那么Session为nil。
}

type RawNetMsg struct {
	Api     int16
	Data    []byte
	Session *Session // 如果为作为客户端连接的受到的消息，那么Session为nil。
}

func readMsg(conn *net.TCPConn, header []byte) *[]byte {
	conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
	_, err := io.ReadFull(conn, header)
	if err != nil {
		if err.Error() == "EOF" {
			fmt.Printf("Recv EOF: %p.", conn)
			return nil
		} else {
			fmt.Printf("Recv err: %p %v.", conn, err)
			return nil
		}
	}

	size := binary.BigEndian.Uint16(header)
	data := make([]byte, size)
	conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
	_, err = io.ReadFull(conn, data)
	if err != nil {
		fmt.Printf("Recv msg data: %p %v.", conn, err)
		return nil
	}
	return &data
}

func DecodeMsg(data []byte) (api int16, v interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("DECODE MSG ERR: %v %v", data, err)
			api = -1
			v = nil
		}
	}()

	//	解包协议 package protocol
	api, payload := protocol.DecodeMsg(&data)
	return api, payload
}
