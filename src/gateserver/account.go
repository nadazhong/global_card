package gateserver

//import (
//	. "account"
//	. "platform"

//	"accountdb"
//	"cfg"
//	"misc/packet"
//	"net/url"
//	"protocol"
//	"servernet"
//	"strconv"
//)

//func initAccount() {
//	//用户id分配器
//	accountdb.StartIdPool()
//	//初始化账号服务器
//	InitManager()
//	//账号数据库
//	accountdb.OpenAccountDB()

//	NEWUSER_GS = accountdb.GetDefaultGs()
//	cfg.LogWarn("default gs:", NEWUSER_GS)
//}

//func registerVerifyCB(ver *VerifyInfo) {
//	cfg.Log("RegisterVerify Callback")
//	msg := ver.Data.(protocol.PKT_login_info)
//	usersession := sessions[ver.Sessionid]
//	if usersession == nil {
//		cfg.LogFatal("user session id fail:", ver.Sessionid)
//		return
//	}
//	session := usersession.session
//	sys := ver.System

//	//验证结果
//	if ver.VerifyCode != 0 { //验证失败
//		cfg.LogFatal("user verify fail:", ver.VerifyCode)
//		session.Send(commonAck(protocol.REGISTER_REQ, protocol.ACCOUNT_VERIFY_FAIL))
//		return
//	}
//	cfg.Log("user verify ok")
//	//账号是否已经注册
//	if ok := sys.Exist(ver.User.UserName); ok != nil {
//		cfg.Log("register: exist account:%s type:%d", ver.User.UserName, sys.TypeId)
//		session.Send(commonAck(protocol.REGISTER_REQ, protocol.USER_NAME_EXISTED))
//		return
//	}
//	//构造游戏账号信息
//	acc := &accountdb.DbAccount{
//		Udid:        msg.F_open_udid,
//		Account:     ver.User.UserName,
//		AccountType: msg.F_account_type,
//	}
//	//选取一个server
//	acc.ServerId = NEWUSER_GS //getPriorityGsId()
//	acc.RegisterServerId = acc.ServerId
//	//分配userid
//	acc.Id = accountdb.NewUserId(acc.ServerId)
//	//创建
//	sys.Create(acc)

//	cfg.Logf("new user:%v", acc)

//	registerMsg := ver.Data.(protocol.PKT_login_info)
//	//通知gs初始化用户信息
//	gsMsg := protocol.PKT_gt_login_info{F_user_id: acc.Id,
//		F_ip:             registerMsg.F_ip,
//		F_udid:           registerMsg.F_open_udid,
//		F_client_version: registerMsg.F_client_version,
//		F_os_version:     registerMsg.F_os_version,
//		F_device_name:    registerMsg.F_device_name,
//		F_device_id:      registerMsg.F_device_id,
//		F_device_id_type: registerMsg.F_device_id_type,
//		F_locale:         registerMsg.F_locale}
//	Buf := packet.Pack(protocol.GT_REGISTER_REQ, gsMsg, nil)

//	forwardMsg := protocol.PKT_forward_msg_info{
//		F_msg:        string(Buf),
//		F_size:       int32(len(Buf)),
//		F_user_id:    acc.Id,
//		F_session_id: ver.Sessionid,
//	}
//	gsSend(acc.ServerId, protocol.FORWARD_MSG_NTF, forwardMsg)
//}

//func register(userSession *gsnet.Session, msg protocol.PKT_login_info) {
//	uSession := getUserSession(userSession)
//	if uSession.login {
//		userSession.Send(commonAck(protocol.REGISTER_REQ, protocol.USER_LOGIN_ALREADY))
//		return
//	}
//	if msg.F_account_type != 0 {
//		go registerP(userSession, msg)
//		return
//	}
//	//通过账号类型选择账号系统
//	sys := GetAccountSystem(msg.F_account_type)
//	if sys == nil {
//		cfg.LogWarnf("register: account type:%d err", msg.F_account_type)
//		userSession.Send(commonAck(protocol.REGISTER_REQ, protocol.ACCOUNT_TYPE_ERR))
//		return
//	}
//	//开启验证，验证完成后使用Callback回调继续处理
//	ver := VerifyInfo{VerifyCode: 0,
//		Token:     msg.F_token,
//		Sessionid: userSession.Id,
//		Callback:  registerVerifyCB,
//		Data:      msg,
//		System:    sys,
//	}
//	sys.StartVerify(&ver)
//}

//func loginVerifyCB(ver *VerifyInfo) {
//	cfg.Log("Loginverify callback ")
//	usersession := sessions[ver.Sessionid]
//	if usersession == nil {
//		cfg.LogFatal("user session id fail:", ver.Sessionid)
//		return
//	}
//	session := usersession.session
//	sys := ver.System
//	//验证结果
//	if ver.VerifyCode != 0 { //验证失败
//		cfg.LogFatal("user verify fail:", ver.VerifyCode)
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.ACCOUNT_VERIFY_FAIL))
//		return
//	}
//	cfg.Log("user verify ok")
//	//查看账号是否存在
//	user := sys.Exist(ver.User.UserName)
//	if user == nil {
//		cfg.Logf("user not exist account:%s account:%d", ver.User.UserName, sys.TypeId)
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.USER_NOT_EXIST))
//		return
//	}
//	loginMsg := ver.Data.(protocol.PKT_login_info)
//	//通知gs初始化用户信息
//	msg := protocol.PKT_gt_login_info{F_user_id: user.Id,
//		F_ip:             loginMsg.F_ip,
//		F_udid:           loginMsg.F_open_udid,
//		F_client_version: loginMsg.F_client_version,
//		F_os_version:     loginMsg.F_os_version,
//		F_device_name:    loginMsg.F_device_name,
//		F_device_id:      loginMsg.F_device_id,
//		F_device_id_type: loginMsg.F_device_id_type,
//		F_locale:         loginMsg.F_locale}
//	Buf := packet.Pack(protocol.GT_LOGIN_REQ, msg, nil)

//	forwardMsg := protocol.PKT_forward_msg_info{
//		F_msg:        string(Buf),
//		F_size:       int32(len(Buf)),
//		F_user_id:    user.Id,
//		F_session_id: ver.Sessionid,
//	}
//	gsSend(user.ServerId, protocol.FORWARD_MSG_NTF, forwardMsg)
//}
//func login(userSession *gsnet.Session, msg protocol.PKT_login_info) {
//	uSession := getUserSession(userSession)
//	if uSession.login {
//		userSession.Send(commonAck(protocol.LOGIN_REQ, protocol.USER_LOGIN_ALREADY))
//		return
//	}
//	if msg.F_account_type != 0 {
//		go loginP(userSession, msg)
//		return
//	}

//	//通过账号类型选择账号系统
//	sys := GetAccountSystem(msg.F_account_type)
//	if sys == nil {
//		cfg.LogWarnf("login: account type:%d err", msg.F_account_type)
//		userSession.Send(commonAck(protocol.LOGIN_REQ, protocol.ACCOUNT_TYPE_ERR))
//		return
//	}

//	//开启验证，验证完成后使用Callback回调继续处理
//	ver := VerifyInfo{VerifyCode: 0,
//		Token:     msg.F_token,
//		Sessionid: userSession.Id,
//		Callback:  loginVerifyCB,
//		Data:      msg,
//		System:    sys,
//	}
//	sys.StartVerify(&ver)
//	//返回
//}

//func bindVerifyCB(ver *VerifyInfo) {
//	cfg.Log("Bindverify callback ")
//	usersession := sessions[ver.Sessionid]
//	if usersession == nil {
//		cfg.LogFatal("user session id fail:", ver.Sessionid)
//		return
//	}
//	session := usersession.session
//	sys := ver.System
//	//验证结果
//	if ver.VerifyCode != 0 { //验证失败
//		cfg.LogFatal("user verify fail:", ver.VerifyCode)
//		session.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.ACCOUNT_VERIFY_FAIL))
//		return
//	}
//	cfg.Log("user verify ok")

//	//判断目标账号是否已经绑定过，不能重复绑定
//	if user := sys.Exist(ver.User.UserName); user != nil {
//		cfg.Logf("Bind:cant rebind account:%s account:%d", ver.User.UserName, sys.TypeId)
//		session.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.ACCOUNT_REBIND))
//		return
//	}
//	//原用户信息
//	oldUser := GetAccountById(usersession.id)
//	if oldUser == nil {
//		cfg.LogErrf("get account by id:%d err", usersession.id)
//		session.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.USER_NOT_EXIST))
//		return
//	}
//	//构造新账号信息,只是改一下账号名和账号类型
//	user := &accountdb.DbAccount{
//		Udid:        oldUser.Udid,
//		ServerId:    oldUser.ServerId,
//		Account:     ver.User.UserName,
//		AccountType: sys.TypeId,
//		Id:          oldUser.Id,
//	}

//	//验证成功后，更新db
//	sys.Update(oldUser, user)

//	session.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.OK))
//}
//func bind(userSession *gsnet.Session, msg protocol.PKT_account_bind_info) {
//	//用户必须登陆
//	uSession := getUserSession(userSession)
//	if !uSession.login {
//		userSession.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.USER_NOT_LOGIN))
//		return
//	}
//	//获取用户资料
//	user := GetAccountById(uSession.id)
//	if user == nil {
//		cfg.LogErrf("get account by id:%d err", uSession.id)
//		userSession.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.USER_NOT_EXIST))
//		return
//	}
//	//只有临时账号可以被绑定
//	if user.AccountType != 0 {
//		cfg.Logf("bind old type err oldtype:%d", user.AccountType)
//		userSession.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.ACCOUNT_BIND_OLDTYPE_ERR))
//		return
//	}
//	//是否是支持的绑定目标类型
//	if msg.F_account_type == 0 {
//		cfg.Logf("bind new type err newtype:%d", msg.F_account_type)
//		userSession.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.ACCOUNT_BIND_NEWTYPE_ERR))
//		return
//	}
//	//通过账号类型选择账号系统
//	sys := GetAccountSystem(msg.F_account_type)
//	if sys == nil {
//		cfg.LogWarnf("register: account type:%d err", msg.F_account_type)
//		userSession.Send(commonAck(protocol.ACCOUNT_BIND_REQ, protocol.ACCOUNT_TYPE_ERR))
//		return
//	}
//	//开启验证，验证完成后使用Callback回调继续处理
//	ver := VerifyInfo{VerifyCode: 0,
//		Token:     msg.F_token,
//		Sessionid: userSession.Id,
//		Callback:  bindVerifyCB,
//		Data:      msg,
//		System:    sys,
//	}

//	sys.StartVerify(&ver)
//	//返回
//}

//func registerP(session *gsnet.Session, msg protocol.PKT_login_info) {
//	uSession := getUserSession(session)
//	if uSession.login {
//		session.Send(commonAck(protocol.REGISTER_REQ, protocol.USER_LOGIN_ALREADY))
//		return
//	}
//	if AuthUrl == "" {
//		session.Send(commonAck(protocol.REGISTER_REQ, protocol.AUTH_INIT_ERR))
//		return
//	}
//	//创建角色
//	value := url.Values{}
//	value.Add("access_token", msg.F_token)
//	value.Add("for_user", "me")
//	value.Add("server_id", "0")
//	cfg.Log("user token:\n", msg.F_token)
//	var ac AuthCreate
//	errCode := HttpSend(POST, HTTPS, AuthUrl, RULE_PATH_CREATE, value, &ac)
//	//回包错误
//	if errCode != protocol.OK {
//		cfg.Log("auth create err:", errCode)
//		session.Send(commonAck(protocol.REGISTER_REQ, errCode))
//		return
//	}
//	userId := ac.Character_id
//	if userId <= 0 {
//		cfg.Log("auth create ok, userId err:", userId)
//		session.Send(commonAck(protocol.REGISTER_REQ, errCode))
//		return
//	}
//	cfg.Log("auth create ok, userId:", userId)

//	//登录
//	loginSimple(session, userId, msg)
//	return
//}

//func loginSimple(session *gsnet.Session, userId int64, msg protocol.PKT_login_info) {
//	uSession := getUserSession(session)
//	value := url.Values{}
//	value.Add("access_token", msg.F_token)
//	value.Add("for_user", "me")
//	value.Add("server_id", "0")
//	value.Add("character_id", strconv.Itoa(int(userId)))
//	var al AuthLogin
//	errCode := HttpSend(POST, HTTPS, AuthUrl, RULE_PATH_LOGIN, value, &al)
//	//回包错误
//	if errCode != protocol.OK {
//		cfg.Log("auth login err:", errCode)
//		session.Send(commonAck(protocol.LOGIN_REQ, errCode))
//		return
//	}
//	if al.Success != true {
//		cfg.Log("auth login fail")
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.AUTH_VERIFY_ERR))
//		return
//	}
//	cfg.Log("auth login ok:")

//	var serverId int16
//	if acc := accountdb.GetAccountById(userId); acc != nil {
//		serverId = acc.ServerId
//	} else { //找不到就新建
//		//构造游戏账号信息
//		acc = &accountdb.DbAccount{
//			Udid:        msg.F_open_udid,
//			Account:     "",
//			AccountType: 0,
//		}
//		//选取一个server
//		serverId = NEWUSER_GS //getPriorityGsId()
//		acc.ServerId = serverId
//		acc.RegisterServerId = acc.ServerId
//		//分配userid
//		acc.Id = userId
//		//创建
//		CreateAccount(acc)

//		cfg.Logf("new user:%v", acc)
//	}
//	//通知gs初始化用户信息
//	fMsg := protocol.PKT_gt_login_info{F_user_id: userId,
//		F_ip:             msg.F_ip,
//		F_udid:           msg.F_open_udid,
//		F_client_version: msg.F_client_version,
//		F_os_version:     msg.F_os_version,
//		F_device_name:    msg.F_device_name,
//		F_device_id:      msg.F_device_id,
//		F_device_id_type: msg.F_device_id_type,
//		F_locale:         msg.F_locale}
//	Buf := packet.Pack(protocol.GT_LOGIN_REQ, fMsg, nil)

//	forwardMsg := protocol.PKT_forward_msg_info{
//		F_msg:        string(Buf),
//		F_size:       int32(len(Buf)),
//		F_user_id:    userId,
//		F_session_id: uSession.session.Id,
//	}
//	gsSend(serverId, protocol.FORWARD_MSG_NTF, forwardMsg)
//}

//func loginP(session *gsnet.Session, msg protocol.PKT_login_info) {
//	uSession := getUserSession(session)
//	if uSession.login {
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.USER_LOGIN_ALREADY))
//		return
//	}
//	if AuthUrl == "" {
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.AUTH_INIT_ERR))
//		return
//	}
//	//拉取角色
//	value := url.Values{}
//	value.Add("access_token", msg.F_token)
//	value.Add("for_user", "me")
//	cfg.Log("values:", value.Encode())
//	var agr AuthGetRules
//	errCode := HttpSend(POST, HTTPS, AuthUrl, RULE_PATH_GETLIST, value, &agr)
//	//回包错误
//	if errCode != protocol.OK {
//		cfg.Log("auth getrules err:", errCode)
//		session.Send(commonAck(protocol.LOGIN_REQ, errCode))
//		return
//	}
//	if agr.Characters == nil {
//		cfg.Log("character_id nil")
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.AUTH_DATA_ERR))
//		return
//	}
//	if len(agr.Characters) == 0 { //没创建角色
//		cfg.Log("character_id == 0")
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.USER_NOT_EXIST))
//		return
//	}
//	if len(agr.Characters) != 1 {
//		cfg.Log("len(characters) > 1:", len(agr.Characters))
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.AUTH_DATA_ERR))
//		return
//	}
//	userId := agr.Characters[0].Character_id
//	//si, _ := strconv.Atoi(agr.Characters[0].Server_id)
//	//serverId := int16(si)
//	if userId <= 0 {
//		cfg.LogWarnf("AuthGetRules data err, userid:%d ", userId)
//		session.Send(commonAck(protocol.LOGIN_REQ, protocol.AUTH_DATA_ERR))
//		return
//	}
//	cfg.Logf("auth getrules userid:%d\n", userId)

//	//登录
//	loginSimple(session, userId, msg)
//	return
//}

//func setDefaultGs(msg protocol.PKT_gs_set_default_gs_info) {
//	NEWUSER_GS = msg.F_gsid
//	accountdb.UpdateDefaultGs(msg.F_gsid)
//	cfg.LogWarn("default gs:", msg.F_gsid)
//}
