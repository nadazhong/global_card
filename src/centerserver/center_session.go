package centerserver

import (
	"fmt"
	//"gsdb"
	"misc/packet"
	"protocol"
	"servernet"
)

var (
	gsMap         map[int32]*centerGs
	pendingGsMap  map[int32]*serverGs
	pendingSaving map[int16]bool
	savingTag     int64
)

type centerGs struct {
	session  *servernet.Session
	serverId int16
}
type serverGs struct {
	session *servernet.Session
}

func (gs *centerGs) Send(api int16, tbl interface{}) {
	gs.session.Send(packet.Pack(api, tbl, nil))
}

func (gs *centerGs) CommonAck(uid int64, api int16, code int16) {
	ack := protocol.PKT_gs_common_ack_info{
		F_user_id: uid,
		F_api:     api,
		F_code:    code,
	}
	gs.Send(protocol.GS_COMMON_ACK, &ack)
}

func (gs *centerGs) SendToUser(userId int64, api int16, tbl interface{}) {
	payload := packet.Pack(api, tbl, nil)
	msg := protocol.PKT_forward_msg_info{
		F_size:       int32(len(payload)),
		F_session_id: protocol.NIL_SESSION_ID,
		F_user_id:    userId,
		F_msg:        string(payload),
	}
	gs.session.Send(packet.Pack(protocol.FORWARD_MSG_NTF, &msg, nil))
}

// 发送公会广播信息（通过Gate广播，所以这条转发消息发给任意GS都可以）。
func AllianceBoardCast(asId int32, api int16, tbl interface{}) {
	payload := packet.Pack(api, tbl, nil)
	for _, gs := range gsMap {
		msg := protocol.PKT_gs_alliance_boardcast_info{
			F_alliance_id: asId,
			F_msg:         string(payload),
		}
		gs.session.Send(packet.Pack(protocol.GS_ALLIANCE_BOARDCAST_NTF,
			&msg, nil))
		return
	}
}

func init() {
	gsMap = make(map[int32]*centerGs)
	pendingGsMap = make(map[int32]*serverGs)
	pendingSaving = make(map[int16]bool)
	savingTag = 0
}

func newGs(session *servernet.Session) {
	gs := serverGs{
		session: session,
	}

	pendingGsMap[session.Id] = &gs
	servernet.StartSession(session, nil)
}

func closeGs(session *servernet.Session) {

	gs, present := gsMap[session.Id]
	if !present {
		delete(pendingGsMap, session.Id)
		return
	}
	fmt.Printf("Gs %v quit, session %v. \n", gs.serverId, gs.session.Id)
	delete(gsMap, session.Id)
	if centerQuit {
		fmt.Printf("Center is quiting, gs cnt: %v. \n", len(gsMap))
		if len(gsMap) == 0 {
			go centerQuitRoutine()
		}
	}
}
func conectHandler(msg *servernet.NetMsg) {

	if msg.Api == protocol.GS_CONNECTED_NTF {
		p := msg.Payload.(protocol.PKT_gs_connected_info)
		bindGs(msg.Session, p.F_server_id)
		return
	}
}

func bindGs(session *servernet.Session, serverId int16) {
	gs, present := gsMap[session.Id]
	if present {
		fmt.Printf("Rebind gs %v at %v. \n", serverId, session.Id)
		return
	}

	server, ok := pendingGsMap[session.Id]
	if !ok {
		fmt.Println("Pending err:", pendingGsMap, session.Id)
	}

	fmt.Printf("Bind gs %v at %d. \n", serverId, session.Id)
	gs = &centerGs{
		serverId: serverId,
		session:  server.session,
	}
	delete(pendingGsMap, session.Id)
	gsMap[session.Id] = gs
}

func msgHander(msg *servernet.NetMsg) {
	gs := gsMap[msg.Session.Id]
	if gs == nil && msg.Api != protocol.GS_CONNECTED_NTF {
		return
	}

	switch msg.Api {
	//case protocol.HEART_BEAT_REQ:
	//	heartBeat(gs)
	case protocol.GS_CONNECTED_NTF:
		p := msg.Payload.(protocol.PKT_gs_connected_info)
		bindGs(msg.Session, p.F_server_id)

	//case protocol.GS_JOIN_ALLIANCE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_user_alliance_info)
	//	joinAllianceReq(gs, p)

	//case protocol.GS_DISMISS_ALLIANCE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_user_alliance_info)
	//	dismissAllianceReq(gs, p)

	//case protocol.GS_QUIT_ALLIANCE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_user_alliance_info)
	//	quitAllianceReq(gs, p)

	//case protocol.GS_CREATE_ALLIANCE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_create_alliance_info)
	//	createAllianceReq(gs, p)
	//case protocol.GS_UPDATE_CU_NTF:
	//	p := msg.Payload.(protocol.PKT_gs_cu_info)
	//	updateCUNtf(gs, p)
	//case protocol.GS_ALLIANCE_LIST_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_user_alliance_info)
	//	allianceListReq(gs, p)

	//case protocol.GS_ALLIANCE_MEMBER_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_user_alliance_info)
	//	allianceMemberReq(gs, p)

	//case protocol.GS_NEWHELP_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_help_info)
	//	allianceNewHelpReq(gs, p)

	//case protocol.GS_OFFER_HELP_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_offer_help_info)
	//	allianceOfferHelpReq(gs, p)

	//case protocol.GS_ALLIANCE_LIST_HELP_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_list_help_req)
	//	allianceListHelpReq(gs, p)

	//case protocol.GS_OFFER_HELPALL_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_offer_helpall_info)
	//	allianceOfferHelpAllReq(gs, p)

	//case protocol.GS_DELETE_HELP_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_delete_help_info)
	//	allianceDeleteHelp(gs, p)

	//case protocol.GS_ALLIANCE_RANKLIST_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_ranklist_info)
	//	allianceRankListReq(gs, p)

	//case protocol.GS_DETAIL_RANKLIST_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_detail_ranklist_info)
	//	allianceDetailRanklistReq(gs, p)

	//case protocol.GS_ALLIANCE_PAYTAX_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_paytax_info)
	//	payTaxReq(gs, p)

	//case protocol.GS_ALLIANCE_MOD_RANK_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_mod_rank_info)
	//	allianceModRankReq(gs, p)

	//case protocol.GS_CHANGE_ALLIANCE_NAME_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_change_alliance_name_info)
	//	allianceChangeNameReq(gs, p)

	//case protocol.GS_CHANGE_ALLIANCE_DESC_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_change_alliance_desc_info)
	//	allianceChangeDescReq(gs, p)

	//case protocol.GS_CHANGE_ALLIANCE_LANG_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_change_alliance_lang_info)
	//	allianceChangeLangReq(gs, p)

	//case protocol.GS_CHANGE_ALLIANCE_LOGO_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_change_alliance_logo_info)
	//	allianceChangeLogoReq(gs, p)

	//case protocol.GS_CHANGE_ALLIANCE_HEADLINE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_change_alliance_headline_info)
	//	allianceChangeHeadLineReq(gs, p)

	//case protocol.GS_ALLIANCE_BLOCK_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_block_info)
	//	allianceBlockReq(gs, p)

	//case protocol.GS_ALLIANCE_UNBLOCK_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_block_info)
	//	allianceUnBlockReq(gs, p)

	//case protocol.GS_ALLIANCE_GET_BLOCKLIST_REQ:
	//	p := msg.Payload.(protocol.PKT_alliance_get_blocklist_info)
	//	allianceGetBlockListReq(gs, p)

	//case protocol.GS_CHANGE_ALLIANCE_OWNER_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_change_alliance_owner_info)
	//	changeAllianceOwnerReq(gs, p)

	//case protocol.GS_ALLIANCE_INVITE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_invite_info)
	//	allianceInviteReq(gs, p)

	//case protocol.GS_ALLIANCE_INVITE_CONFIRM_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_invite_confirm_info)
	//	allianceInviteConfirmReq(gs, p)

	//case protocol.GS_ALLIANCE_MANAGE_TORRITORY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_manage_torritory_info)
	//	allianceManageTorritoryReq(gs, p)

	//case protocol.GS_ALLIANCE_RES_LIST_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_res_list_info)
	//	allianceResListReq(gs, p)

	//case protocol.GS_ALLIANCE_RES_HISTORY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_res_history_info)
	//	allianceResHistoryReq(gs, p)

	//case protocol.GS_LIST_ALLIANCE_INVITE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_alliance_invite_info)
	//	listAllianceInviteReq(gs, p)

	//case protocol.GS_LIST_INVITE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_invite_info)
	//	listInviteReq(gs, p)

	//case protocol.GS_REVOKE_ALLIANCE_INVITE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_revoke_alliance_invite_info)
	//	revokeAllianceInviteReq(gs, p)

	//case protocol.GS_ALLIANCE_RES_APPLY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_res_apply_info)
	//	allianceResApplyReq(gs, p)

	//case protocol.GS_LIST_ALLIANCE_RES_APPLY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_alliance_res_apply_info)
	//	listAllianceResApplyReq(gs, p)

	//case protocol.GS_MANAGE_ALLIANCE_RES_APPLY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_manage_alliance_res_apply_info)
	//	manageAllianceResApplyReq(gs, p)

	//case protocol.GS_LIST_ALLIANCE_DOMINION_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_alliance_dominion_info)
	//	listAllianceDominionReq(gs, p)

	//case protocol.GS_ALLIANCE_RES_APPLY_CANCEL_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_res_apply_cancel_info)
	//	allianceResApplyCancelReq(gs, p)
	//case protocol.GS_RELOCATE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_relocate_info)
	//	gsRelocateReq(gs, p)
	//case protocol.GS_RELOCATE_ACK:
	//	p := msg.Payload.(protocol.PKT_gs_relocate_ack_info)
	//	gsRelocateAck(p)

	//case protocol.GS_HELP_SPEEDUP_ACK:
	//	p := msg.Payload.(protocol.PKT_gs_help_speedup_ack_info)
	//	gsHelpSpeedupAck(gs, p)

	//case protocol.GS_ALLIANCE_EVENT_NTF:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_event_info)
	//	gsAllianceEventNtf(gs, p)

	//case protocol.GS_ALLIANCE_EVENT_DETAIL_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_alliance_event_detail_info)
	//	gsAllianceEventDetailInfoReq(gs, p)

	//case protocol.GS_SEARCH_NAME_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_search_name_req)
	//	gsSearchUser(gs, p)
	////case protocol.SERVER_QUIT_NTF:
	////	centerQuit = true
	////	if len(gsMap) == 0 {
	////		go centerQuitRoutine()
	////	}

	//case protocol.GS_NAME_EXIST_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_name_exist_req)
	//	NameExist(gs, p)

	//case protocol.SNAPSHOT_NTF:
	//	p := msg.Payload.(protocol.PKT_snapshot_info)
	//	onSnapshotDB(gs, p)

	//case protocol.RELOAD_ALLIANCE_EVENT_NTF:
	//	reloadAllianceEvent()

	//case protocol.LIST_ALLIANCE_EVENTS_REQ:
	//	p := msg.Payload.(protocol.PKT_list_alliance_events_req_info)
	//	listAllianceEventsReq(gs, p)

	//case protocol.GS_POST_ALLIANCE_COMMENT_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_post_alliance_comment_info)
	//	postAllianceCommentsReq(gs, p)

	//case protocol.GS_LIST_ALLIANCE_COMMENT_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_alliance_comment_info)
	//	listAllianceCommentsReq(gs, p)

	//case protocol.GS_APPLY_JOIN_ALLIANCE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_apply_join_alliance_info)
	//	applyJoinAllianceReq(gs, p)

	//case protocol.GS_LIST_ALLIANCE_JOIN_APPLY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_alliance_join_apply_info)
	//	listAllianceJoinApplyReq(gs, p)

	//case protocol.GS_MANAGE_ALLIANCE_JOIN_APPLY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_manage_alliance_join_apply_info)
	//	manageAllianceJoinApplyReq(gs, p)

	//case protocol.GS_LIST_ALLIANCE_AUTHORITY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_list_alliance_authority_req)
	//	listAllianceAuthorityReq(gs, p)

	//case protocol.GS_UPDATE_ALLIANCE_AUTHORITY_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_update_alliance_authority_info)
	//	updateAllianceAuthorityReq(gs, p)

	//case protocol.GS_MANAGE_ALLIANCE_RECRUITMENT_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_manage_alliance_recruitment_info)
	//	manageAllianceRecruitmentReq(gs, p)

	//case protocol.GS_UPDATE_ALLIANCE_MEMBER_TITLE_REQ:
	//	p := msg.Payload.(protocol.PKT_gs_update_alliance_member_title_info)
	//	updateAllianceMemberTitleReq(gs, p)

	//case protocol.SWAP_USER_REQ:
	//	p := msg.Payload.(protocol.PKT_swap_user_info)
	//	swapUserReq(p)

	//case protocol.GET_SWAP_USER_REQ:
	//	p := msg.Payload.(protocol.PKT_get_swap_user_info)
	//	getSwapUserReq(gs, p)

	default:
	}
}

//func heartBeat(gs *centerGs) {
//	fmt.Printf("Gs %v heart beat msg. \n", gs.serverId)
//	gs.Send(protocol.HEART_BEAT_REQ, protocol.PKT_null{})
//}

func getGSById(serverId int16) *centerGs {
	for _, gs := range gsMap {
		if gs.serverId == serverId {
			return gs
		}
	}
	return nil
}

func gsRelocateReq(gs *centerGs, info protocol.PKT_gs_relocate_info) {
	targetGs := getGSById(info.F_dst.F_server_id)
	if targetGs == nil {
		ack := protocol.PKT_gs_relocate_ack_info{
			F_code: protocol.RELOCATE_SERVER_NOT_EXIST,
			F_info: info,
		}
		gs.Send(protocol.GS_RELOCATE_ACK, ack)
		return
	}
	targetGs.Send(protocol.GS_RELOCATE_REQ, &info)
}

//func swapUserReq(msg protocol.PKT_swap_user_info) {
//	cfg.LogErrf("swap user req, len %v.", len(msg.F_data))
//	err := gsdb.PushSwapUser(msg)
//	if err != nil {
//		cfg.LogErrf("Push swap user fail: %v", err)
//	}
//}

//func getSwapUserReq(gs *centerGs, msg protocol.PKT_get_swap_user_info) {
//	swapNtf := gsdb.PullSwapUser(msg.F_user_id)
//	if swapNtf != nil {
//		gs.Send(protocol.SWAP_USER_NTF, swapNtf)
//	}
//}
