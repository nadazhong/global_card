package account

/*
* 对账号信息进行处理
 */
//import (
//	"accountdb"
//)

//const (
//	ANONYMOUS_TYPE = 0
//	VERIFY_BUF     = 1000
//)

//var (
//	accounts   map[int16]*AccountSystem       // 各个账号系统下的账号缓存
//	VerifyChan chan *VerifyInfo               // 保存完成的操作
//	id2Account map[int64]*accountdb.DbAccount // UserId 到账号信息的缓存。
//)

//type VerifyCallBack func(acc *VerifyInfo)

////账号验证成功后，返回的用户基础信息
//type UserInfo struct {
//	UserName string
//}

////验证完成回调信息
//type VerifyInfo struct {
//	/*input*/
//	Callback  VerifyCallBack //回调函数
//	Token     string         //用于验证的token
//	Sessionid int32          //当前会话
//	System    *AccountSystem //当前账号系统
//	Data      interface{}    //用于callback保存自己的定制信息
//	/*output*/
//	VerifyCode int16    //验证返回码
//	User       UserInfo //账号基本信息
//}

////一个账号系统应该提供的接口
//type SystemInterface interface {
//	init(sys *AccountSystem) bool
//	startVerify(msg *VerifyInfo, sys *AccountSystem) bool
//	unInit(sys *AccountSystem) bool
//}

////一个账号系统下的信息，如：本地账号、tap4fun、facebook
//type AccountSystem struct {
//	TypeId int16
//	//暂时不缓存users   map[string]*accountdb.DbAccount
//	handler SystemInterface
//}

////验证账号类型是否合法
//func GetAccountSystem(atype int16) *AccountSystem {
//	sys, exist := accounts[atype]
//	//账号系统是否支持
//	if !exist {
//		return nil
//	}
//	return sys
//}

////判断一个账号是否存在
//func (sys *AccountSystem) Exist(aname string) *accountdb.DbAccount {
//	//在数据库中查看是否存在
//	if user := accountdb.GetAccountByName(aname, sys.TypeId); user != nil {
//		//	sys.users[aname] = user
//		id2Account[user.Id] = user
//		return user
//	}
//	return nil
//}

////开始账号验证，有可能会有异步操作
//func (sys *AccountSystem) StartVerify(ver *VerifyInfo) bool {
//	return sys.handler.startVerify(ver, sys)
//}

////更新账号
//func (sys *AccountSystem) Update(oldAcc *accountdb.DbAccount, newAcc *accountdb.DbAccount) {
//	//创建新的账号信息
//	id2Account[newAcc.Id] = newAcc
//	//更新数据库，id不变
//	accountdb.UpdateAccount(newAcc.Id, *newAcc)
//}

////创建新账号
//func (sys *AccountSystem) Create(acc *accountdb.DbAccount) {
//	//sys.users[acc.Account] = acc
//	id2Account[acc.Id] = acc
//	accountdb.UpdateAccount(acc.Id, *acc)
//}

//// 创建新账号
//func CreateAccount(acc *accountdb.DbAccount) {
//	//sys.users[acc.Account] = acc
//	id2Account[acc.Id] = acc
//	accountdb.UpdateAccount(acc.Id, *acc)
//}

//// 通过 id 拉去用户信息
//func GetAccountById(id int64) *accountdb.DbAccount {
//	acc, exist := id2Account[id]
//	if exist {
//		return acc
//	}
//	acc = accountdb.GetAccountById(id)
//	if acc != nil {
//		//accounts[acc.AccountType].users[acc.Account] = acc
//		id2Account[acc.Id] = acc
//		return acc
//	}
//	return nil
//}

//// 账号管理系统初始化
//func InitManager() {
//	VerifyChan = make(chan *VerifyInfo, VERIFY_BUF)
//	id2Account = make(map[int64]*accountdb.DbAccount)

//	accounts = map[int16]*AccountSystem{
//		//本地临时账号系统
//		ANONYMOUS_TYPE: {TypeId: ANONYMOUS_TYPE,
//			handler: &localAccount{},
//		},
//		1: {TypeId: 1,
//			handler: &localAccount{},
//		},
//	}
//	//初始化各个系统
//	for _, acc := range accounts {
//		acc.handler.init(acc)
//	}
//}
