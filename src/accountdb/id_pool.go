package accountdb

//import (
//	"cfg"
//)

//var (
//	userIdPools map[int16]*IdPool // KEY -> ServerId.
//)

//const (
//	ID_BUF_SIZE = 10

//	ID_USER       = 1
//	ID_USER_TABLE = "accounts"

//	MAX_ID_POOL_OFFSET = 999
//	ID_MASK            = 1000
//)

//type IdPool struct {
//	idType int32
//	curId  int64
//	offset int16
//	idCh   chan int64
//}

//func NewUserId(serverid int16) int64 {
//	id := int16(0)
//	pool, ok := userIdPools[id]
//	// 如果这个server的定时器不存在，则新建一个。
//	if !ok {
//		pool = &IdPool{idCh: make(chan int64, ID_BUF_SIZE)}
//		go pool.start(ID_USER, ID_USER_TABLE, int16(id))
//		userIdPools[id] = pool
//	}
//	return <-pool.idCh
//}

//func StartIdPool() {
//	userIdPools = make(map[int16]*IdPool)
//}

//func (self *IdPool) start(idType int32, idTable string, offset int16) {
//	self.idType = idType
//	self.offset = offset

//	var err error
//	if self.curId, err = getMaxId(idTable); err != nil {
//		cfg.LogErr("Init id pool fail", idTable, idTable)
//		panic(err.Error())
//	}

//	self.curId++
//	cfg.Logf("ID pool %v start, cur %v.", self.idType, self.curId)
//	for {
//		self.idCh <- self.curId
//		cfg.Logf("ID pool %v return %v", self.idType, self.curId)
//		self.curId++
//	}
//}
