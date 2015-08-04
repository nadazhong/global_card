package gstimer

const (
	ACTION_DUMMY                   = 0
	ACTION_FORWARD                 = 1
	ACTION_BUILDING_COMPLETE       = 2
	ACTION_TROOP_TRAIN_COMPLETE    = 3
	ACTION_MARCH_COMPLETE          = 4
	ACTION_GATHER_COMPLETE         = 5
	ACTION_RALLY_COMPLETE          = 6 // Rally集结结束。
	ACTION_RESEARCH_COMPLETE       = 7
	ACTION_BUFF_COMPLETE           = 8
	ACTION_MAP_SAVE                = 9
	ACTION_POWER_UPDATE            = 10
	ACTION_REPAIR_COMPLETE         = 11 // 修理结束
	ACTION_MAP_POS_FIGHTING        = 12 // 地图战斗状态播放
	ACTION_CRAFT_COMPLETE          = 13 // 锻造结束
	ACTION_USER_MAP_INFO_UPDATE    = 14 // user map 更新
	ACTION_USER_ALERT_UPDATE       = 15 // user alert 更新
	ACTION_GARRISON_TRAIN_COMPLETE = 16
	ACTION_POWER_RANK_UPDATE       = 17

	ACTION_FORT_TO_LAND              = 18 // 据点转化为领土更新
	ACTION_LAND_CHANGE_ALLIANCE      = 19 // 领土变更新的主人
	ACTION_DISPATCH_RES              = 20 // 公会分发资源
	ACTION_USER_SHIELD               = 21 // 防护盾
	ACTION_USER_ANTISCOUT            = 22 // 反侦察
	ACTION_ALLIANCE_RES_APPLY_UPDATE = 23 //公会资源申请超时
	ACTION_DOMINION_MAP_INFO_UPDATE  = 24 //领土地图更新
	ACTION_EVENT_UPDATE              = 25 // 活动事件更新

	// 全局消息
	ACTION_SHUTDOWN        = 10000
	ACTION_MAP_UPDATE      = 10001
	ACTION_ACTIVITY_UPDATE = 10002 // 活动更新
	ACTION_USER_AUTO_SAVE  = 10003 // 玩家定时存盘
	ACTION_HEARTBEAT       = 10004 //心跳包定时器
	ACTION_SNAPSHOT_DB     = 10007 // 数据库快照
)

type Msg struct {
	Action int16
	Data   interface{}
}

type GSTimer struct {
	Receiver int64
	Msg      Msg
}
