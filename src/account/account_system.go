package account

/*
*已经支持的各账号系统
 */

////本地账号
//type localAccount struct{}

//func (a *localAccount) init(sys *AccountSystem) bool {
//	return true
//}

//func (a *localAccount) startVerify(msg *VerifyInfo, sys *AccountSystem) bool {
//	msg.VerifyCode = 0            //成功
//	msg.User.UserName = msg.Token //本地账号直接用token作为id
//	msg.Callback(msg)             //因为不需要验证，直接调用回调即可，不使用channel，因为在同一个线程中，可能导致死锁
//	return true
//}

//func (a *localAccount) unInit(sys *AccountSystem) bool {
//	return true
//}
