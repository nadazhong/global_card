package gameserver

import (
	"fmt"
	"protocol"
	"servernet"
	//"time"
)

// 来自Gate的消息。
func gateMsgHanler(msg *servernet.NetMsg) {
	fmt.Println("RECV GATE MSG:", msg.Api)

	//switch msg.Api {
	//case protocol.FORWARD_MSG_NTF:
	//	info := msg.Payload.(protocol.PKT_forward_msg_info)
	//	api, payload := gsnet.DecodeMsg([]byte(info.F_msg))
	//	if payload == nil {
	//		fmt.Printf("Decode gate msg %v fail.", msg.Api)
	//		//用户数据包错误，需要通知gate断掉他的连接
	//		kickMsg := protocol.PKT_user_session_info{
	//			F_session_id: info.F_session_id,
	//		}
	//		SendToGate(protocol.GS_KICKOUT_NTF, &kickMsg)
	//		return
	//	}
	//	netReqHandler(info.F_session_id, int64(info.F_user_id), api, payload)
	//case protocol.GATE_LOGOUT_NTF:
	//	logoutMsg := msg.Payload.(protocol.PKT_user_session_info)
	//	Logout(logoutMsg.F_user_id)
	//case protocol.GUEST_LOGIN_REQ:
	//	guestLogin(msg.Payload.(protocol.PKT_guest_info))
	//case protocol.GUEST_LOGOUT_REQ:
	//	guestLogout(msg.Payload.(protocol.PKT_guest_info))
	//case protocol.GUEST_GRID_VIEW_REQ:
	//	guestGridViewReq(msg.Payload.(protocol.PKT_server_grid_pos))
	//case protocol.SERVER_QUIT_NTF:
	//	go quitRoutine()
	//case protocol.SNAPSHOT_NTF:
	//	beginSnapshotDB(msg.Payload.(protocol.PKT_snapshot_info))
	//default:
	//}
}

// 来自玩家网络的消息处理。
func netReqHandler(sessionId int32, userId int64,
	api int16, payload interface{}) {
	//t := time.Now()
	//SelfSessionId = sessionId

	//if userId == protocol.NIL_USER_ID {
	//	fmt.Println("userId nil")
	//	return
	//}
	//if api == protocol.GT_LOGIN_REQ || api == protocol.GT_REGISTER_REQ {
	//	GT_Handle(api, payload)
	//	SelfSessionId = protocol.NIL_SESSION_ID
	//	Self = nil
	//	fmt.Printf("API %v cost %v.\n", api, time.Now().Sub(t))
	//	return
	//}

	//Self = GetUser(userId)
	//if Self == nil {
	//	fmt.Println("User %v nil.", userId)
	//}

	//Handle(api, payload)
	//Self = nil
	//SelfSessionId = protocol.NIL_SESSION_ID
}

// 告诉gate我(gs)的ID
func updateGateGs(gsId int16) {
	msg := protocol.PKT_gs_connected_info{F_server_id: gsId}
	SendToGate(protocol.GS_CONNECTED_NTF, &msg)
}

func beginSnapshotDB(msg protocol.PKT_snapshot_info) {

	//// save all data && snapshot database
	//dbSnapshot(GetAllUsers(), msg.F_timestamp)

	//// notify center server
	//SendToCenter(protocol.SNAPSHOT_NTF, &msg)

	//onMsg := func(p interface{}) {
	//	msg := p.(protocol.PKT_snapshot_res_info)
	//}
	//onTimeout := func() {
	//}
	//// wait forever until got target message or timeout
	//wait4(protocol.SNAPSHOT_FINISH_NTF, onMsg, onTimeout)
}

func wait4(api int16, onMsg func(p interface{}), onTimeout func()) {
	//loop := true
	//st := time.Now()
	//for loop {
	//	select {
	//	case msg, err := <-centerConn.MsgQueue:
	//		if !err {
	//			continue
	//		}
	//		if msg.Api == api {
	//			onMsg(msg.Payload)
	//			loop = false
	//			break
	//		} else {
	//			DelayMsgQueue <- msg
	//		}
	//	case <-time.After(protocol.GS_WAITTING_TIMEOUT * time.Millisecond):
	//		if onTimeout != nil {
	//			onTimeout()
	//		}
	//		loop = false
	//		break
	//	}
	//}
}
