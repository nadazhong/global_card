package servernet

import (
	"fmt"
	"net"
	"protocol"
	"strings"
	"misc/packet"
)

var (
	curSessionId = int32(0)
)

const (
	TCP_TIMEOUT       = 600
	NETMSG_QUEUE_SIZE = 1024
)

type Cipher interface {
	Encrypt(data []byte) []byte
	Decrypt(data []byte) []byte
}

type Session struct {
	Id          int32
	ip          net.IP
	conn        *net.TCPConn
	sender      *Sender
	kickout     bool
	closeChan   chan *Session
	msgQueue    chan *NetMsg
	rawMsgQueue chan *RawNetMsg
	rawMsg      bool
}

type ClientInfo struct {
	ConnectChan chan *Session
	CloseChan   chan *Session
	MsgQueue    chan *NetMsg
	RawMsgQueue chan *RawNetMsg
}

func NewServer(addr string, rawMsg bool) (error, *ClientInfo) {
	serverInfo := ClientInfo{
		ConnectChan: make(chan *Session, 5),
		CloseChan:   make(chan *Session, 5),
		MsgQueue:    make(chan *NetMsg, NETMSG_QUEUE_SIZE),
		RawMsgQueue: make(chan *RawNetMsg, 5),
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err, nil
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err, nil
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept failed: ", err)
				continue
			}
			news := newSession(conn, &serverInfo, rawMsg)
			fmt.Printf("New session:", news.ip)
			serverInfo.ConnectChan <- news
		}
	}()

	return nil, &serverInfo
}

func newSession(conn *net.TCPConn, serverInfo *ClientInfo,
	rawMsg bool) *Session {
	ip := net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])
	curSessionId++
	session := &Session{
		Id:          curSessionId,
		ip:          ip,
		conn:        conn,
		kickout:     false,
		sender:      newSender(conn),
		closeChan:   serverInfo.CloseChan,
		msgQueue:    serverInfo.MsgQueue,
		rawMsgQueue: serverInfo.RawMsgQueue,
		rawMsg:      rawMsg,
	}

	return session
}

func StartSession(session *Session, cipher Cipher) {
	go session.start(cipher)
	go session.sender.start(cipher)
}

func (session *Session) start(cipher Cipher) {
	defer func() {
		fmt.Printf("Session %p quit.", session)
		session.conn.Close()
		session.sender.Close()
		session.closeChan <- session
	}()

	conn := session.conn
	fmt.Printf("Session %p start.", session)
	header := make([]byte, 2)
	for {
		if session.kickout {
			return
		}

		data := readMsg(conn, header)
		if data == nil {
			return
		}

		if cipher != nil {
			cipher.Decrypt(*data)
		}

		// 如果 rawMsg 为 true，那么不将消息解码。
		if session.rawMsg {
			msg := parseRawMsg(data, session)
			if msg == nil {
				return
			}
			session.rawMsgQueue <- msg
		} else {
			msg := parseMsg(data, session)
			if msg == nil {
				return
			}
			session.msgQueue <- msg
		}
	}
}

func (session *Session) Send(data []byte) error {
	return session.sender.Send(data)
}

func (session *Session) SyncSend(data []byte) error {
	return session.sender.SyncSend(data)
}

func (session *Session) Close() {
	session.conn.Close()
	session.sender.Close()
	session.kickout = true
}

func parseMsg(data *[]byte, session *Session) *NetMsg {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("DECODE MSG ERR: %p %v %v", session, data, err)
		}
	}()

	api, payload := protocol.DecodeMsg(data)
	return &NetMsg{
		Api:     api,
		Payload: payload,
		Session: session,
	}
}

func parseRawMsg(data *[]byte, session *Session) *RawNetMsg {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("DECODE RAW MSG ERR: %p %v %v", session, data, err)
		}
	}()

	reader := packet.Reader(*data)
	api, err := reader.ReadS16()
	if err != nil {
		fmt.Printf("Get header fail: %v", data)
		return nil
	}

	return &RawNetMsg{
		Api:     api,
		Data:    *data,
		Session: session,
	}
}
