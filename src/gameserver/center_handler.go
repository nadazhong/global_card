package gameserver

import (
	"fmt"
	"protocol"
	"servernet"
)

// 来自Center的消息。
func centerMsgHanler(msg *servernet.NetMsg) {
	fmt.Println("RECV CENTER MSG:", msg.Api)

}

func updateCenterGs(gsId int16) {
	msg := protocol.PKT_gs_connected_info{F_server_id: gsId}
	SendToCenter(protocol.GS_CONNECTED_NTF, &msg)
}
