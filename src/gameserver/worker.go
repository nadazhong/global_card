package gameserver

import (
	"fmt"
	"misc/packet"
	//"os"
	"gsdb"
	"gstimer"
	"netcon"
	"protocol"
	"servernet"
	"time"
)

type MsgCounter struct {
	MaxCost   time.Duration
	MinCost   time.Duration
	TotalCost time.Duration
	Total     int64
	Base      int64
}

var (
	reconnGate   bool // 是否需要尝试连接Gate
	reconnCenter bool // 是否需要尝试连接Center

	gateConn   *servernet.ConnInfo = nil
	centerConn *servernet.ConnInfo = nil

	gsInited bool // 表示 gs 是否已经初始化成功过。

	DelayMsgQueue chan *servernet.NetMsg
)

func init() {
	gsInited = false
	reconnCenter = true
	reconnGate = true
	DelayMsgQueue = make(chan *servernet.NetMsg)
}

func (counter *MsgCounter) Update(start time.Time) {
	//cost := time.Now().Sub(start)
	//counter.Total++
	//counter.TotalCost += cost

	//if counter.MaxCost < cost {
	//	counter.MaxCost = cost
	//}

	//if counter.MinCost == 0 || counter.MinCost > cost {
	//	counter.MinCost = cost
	//}

	//if counter.Total-counter.Base >= 10000 {
	//	avgtime := counter.TotalCost / (time.Duration)(counter.Total)
	//	cfg.Logf("STATISTICS: %v avg: %v.\n", counter, avgtime)
	//	counter.Base = counter.Total
	//}
}

var (
	// 主 routine 中的玩家。
	Self          *gsdb.DbUser = nil
	SelfSessionId int32        = protocol.NIL_SESSION_ID
	isShutDown                 = false
	msgCounter    MsgCounter
)

// Gs 向 Gate Server 发送消息。
func SendToGate(api int16, tbl interface{}) {
	if reconnGate || gateConn == nil {
		return
	}
	gateConn.Send(packet.Pack(api, tbl, nil))
}

// Gs 向 Center Server 发送消息。
func SendToCenter(api int16, tbl interface{}) {
	//if reconnCenter || centerConn == nil {
	//	return
	//}
	//centerConn.Send(packet.Pack(api, tbl, nil))
}

// 向当前玩家发送消息。
func Send(api int16, tbl interface{}) {
	//if Self != nil {
	//	sendForwardMsg(Self.Id, SelfSessionId, api, tbl)
	//} else {
	//	sendForwardMsg(protocol.NIL_USER_ID, SelfSessionId, api, tbl)
	//}
}

// 向指定玩家发送消息。
func SendTo(userId int64, api int16, tbl interface{}) {
	//sendForwardMsg(userId, protocol.NIL_SESSION_ID, api, tbl)
}

func sendForwardMsg(userId int64, sessionId int32, api int16, tbl interface{}) {
	//payload := packet.Pack(api, tbl, nil)
	//msg := protocol.PKT_forward_msg_info{
	//	F_size:       int32(len(payload)),
	//	F_session_id: sessionId,
	//	F_user_id:    userId,
	//	F_msg:        string(payload),
	//}
	//SendToGate(protocol.FORWARD_MSG_NTF, &msg)
}

// 全服广播。
func GSBoardCast(api int16, tbl interface{}) {
	//payload := packet.Pack(api, tbl, nil)
	//msg := protocol.PKT_gs_boardcast_info{
	//	F_msg:       string(payload),
	//	F_server_id: GS_ID,
	//}
	//SendToGate(protocol.GS_BOARDCAST_NTF, &msg)
}

// 组播。
func GroupSend(userList []int64, api int16, tbl interface{}) {
	//users := make([]int32, len(userList))
	//for idx, userId := range userList {
	//	users[idx] = int32(userId)
	//}

	//payload := packet.Pack(api, tbl, nil)
	//msg := protocol.PKT_gs_groupcast_info{
	//	F_msg:       string(payload),
	//	F_server_id: GS_ID,
	//	F_user_list: users,
	//}

	//SendToGate(protocol.GS_GROUPCAST_NTF, &msg)
}

// 公会广播(具体广播在Gate上实现)。
func AllianceBoardCast(asId int32, api int16, tbl interface{}) {
	//payload := packet.Pack(api, tbl, nil)
	//msg := protocol.PKT_gs_alliance_boardcast_info{
	//	F_msg:         string(payload),
	//	F_alliance_id: asId,
	//}

	//SendToGate(protocol.GS_ALLIANCE_BOARDCAST_NTF, &msg)
}

func setGateConn(conn *servernet.ConnInfo) {
	gateConn = conn
}

func setCenterConn(conn *servernet.ConnInfo) {
	centerConn = conn
}

func dummyConn() {
	dummyConn := servernet.ConnInfo{
		CloseChan: make(chan int),
		MsgQueue:  make(chan *servernet.NetMsg),
	}
	gateConn = &dummyConn
	centerConn = &dummyConn
}

// 不停的从网络/内部消息队列中处理逻辑。
func startWorker() {
	dummyConn()
	setGsHeartBeatTimer()
	defer func() {
		fmt.Println("Worker quit...")
		//helper.BackTrace("worker")
		if Self != nil {
			fmt.Println(Self.String())
		}
		shutDown()
	}()
	for {
		checkConn()
		select {
		case _, err := <-gateConn.CloseChan:
			reconnGate = true
			fmt.Println("Gate chan closed: ", err)
		case msg, err := <-gateConn.MsgQueue:
			if !err {
				fmt.Println("Net close chan error:", err)
				continue
			}
			if gsAvaliable() {
				gateMsgHanler(msg)
			} else {
				// Ignore...
			}
		case _, err := <-centerConn.CloseChan:
			reconnCenter = true
			fmt.Println("Center chan closed:", err)
		case msg, err := <-centerConn.MsgQueue:
			if !err {
				fmt.Println("Net close chan error:", err)
				continue
			}
			centerMsgHanler(msg) // Center 的消息一定处理。
		//case id, err := <-gstimer.GS_TIMER:
		//	timerChanHandler(id, err)
		case msg, err := <-DelayMsgQueue:
			if !err {
				fmt.Println("Net close chan error:", err)
				continue
			}
			centerMsgHanler(msg)
		}
	}
}

func doSave() {

}

func shutDown() {
	//if isShutDown == false {
	//	isShutDown = true
	//	cfg.Log("shutDown...")
	//	saveAllUser()               // 玩家保存退出。
	//	gsdb.SaveCU("gs_center.db") // User Center退出。
	//	gsmap.SaveMap()             // 保存大地图等。
	//	saveRanking()               // 排行榜保存。
	//	activity.Save()
	//	event.SaveEvents()
	//	gsdb.SaveTimer() // 保存 timer。
	//	cfg.Log("Save completed.")

	//	gsdb.DBClose()
	//	xmpp.Close()
	//	cfg.Log("DB closed.")
	//}
}

func setGsHeartBeatTimer() {
	gstimer.CreateTimer(0, 30000, gstimer.Msg{Action: gstimer.ACTION_HEARTBEAT})
}

func heartBeatTCB() {
	//SendToGate(protocol.HEART_BEAT_REQ, protocol.PKT_null{})
	//SendToCenter(protocol.HEART_BEAT_REQ, protocol.PKT_null{})
	//setGsHeartBeatTimer()
}

// 检查网络需要重新链接
func checkConn() {
	if reconnGate == true {
		if conn := netcon.ConnGate(); conn != nil {
			setGateConn(conn)
			reconnGate = false
			updateGateGs(GS_ID)
		}
	}

	if reconnCenter == true {
		if conn := netcon.ConnCenter(); conn != nil {
			setCenterConn(conn)
			reconnCenter = false
			updateCenterGs(GS_ID)

		}
	}
	gsStarting()

}

// 检查初始化是否完成，调用gs启动的逻辑函数。
func gsStarting() {
	if !gsInited && gsAvaliable() {
		gsRun()
		gsInited = true
	}
}

// GS 是否已经准备就绪。
func gsAvaliable() bool {
	if (reconnGate || gateConn == nil) || (reconnCenter || centerConn == nil) {
		return false
	}
	return true
}

// GS 退出，在所有工作做完之后，发送消息到 center。
func quitRoutine() {
	//setQuitTimer()
	//<-STOP_SIG
	//cfg.LogWarnf("Gs %v to be down.", GS_ID)
	//if centerConn != nil {
	//	centerConn.SyncSend(packet.Pack(protocol.SERVER_QUIT_NTF,
	//		protocol.PKT_null{}, nil))
	//}
	//os.Exit(0)
}
