package servernet

import (
	"errors"
	"fmt"
	"misc/packet"
	"net"
)

const (
	DEFAULT_QUEUE_SIZE = 128
)

type Sender struct {
	quitctrl chan bool // Receive exit signal.
	pending  chan []byte
	max      int
	conn     *net.TCPConn
}

func (sender *Sender) Send(data []byte) (err error) {
	if len(sender.pending) < sender.max {
		sender.pending <- data
		return nil
	} else {
		fmt.Printf("Send Sender Overflow, len(pending)=", len(sender.pending))
		return errors.New(fmt.Sprintf("Send Sender Overflow, Remote: %v",
			sender.conn.RemoteAddr()))
	}
}

// 同步发送消息（可能导致消息乱序）。
func (sender *Sender) SyncSend(data []byte) (err error) {
	return sender.rawSend(data)
}

func (sender *Sender) PendingCnt() int {
	return len(sender.pending)
}

func (sender *Sender) start(cipher Cipher) {
	defer func() {
		recover()
	}()

	for {
		select {
		case data := <-sender.pending:
			if cipher != nil {
				data = cipher.Encrypt(data)
			}
			sender.rawSend(data)
		case _, ok := <-sender.quitctrl:
			if !ok {
				return
			}
		}
	}
}

func (sender *Sender) rawSend(data []byte) error {
	if len(data) > 0xffff {
		fmt.Printf("Send data overflow size:", len(data))
		return errors.New("send data overflow")
	}
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)))
	writer.WriteRawBytes(data)

	_, err := sender.conn.Write(writer.Data())
	if err != nil {
		fmt.Printf("Error send reply :", err)
		return err
	}
	return nil
}

func (sender *Sender) Close() {
	sender.quitctrl <- true
}

func newSender(conn *net.TCPConn) *Sender {
	return &Sender{
		conn:     conn,
		quitctrl: make(chan bool),
		max:      DEFAULT_QUEUE_SIZE,
		pending:  make(chan []byte, DEFAULT_QUEUE_SIZE),
	}
}
