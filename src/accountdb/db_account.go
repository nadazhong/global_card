package accountdb

//import (
//	"cfg"
//	"errors"
//	"fmt"
//	//"labix.org/v2/mgo"
//	//"labix.org/v2/mgo/bson"
//)

//const (
//	ACCOUNT_COLLECTION = "accounts"
//	CONFIG_COLLECTION  = "gate_config"
//	SERVER_COLLECTION  = "servers"

//	CONFIG_TYPE_DEFAULT = 1
//)

//type DbAccount struct {
//	Id               int64
//	ServerId         int16 //用户当前server
//	RegisterServerId int16 //用户初始server
//	Udid             string
//	Account          string
//	AccountType      int16
//}

//type DbServer struct {
//	Id         int16
//	OpenTime   int64
//	Population int16
//}

//var (
//	accountSession *mgo.Session
//)

//func OpenAccountDB() {
//	config := cfg.Get("ACCOUNT")
//	user := config["mongo_user"]
//	pass := config["mongo_pass"]
//	addr := config["mongo_addr"]
//	dbname := config["mongo_dbname"]
//	url := fmt.Sprintf("mongodb://%s:%s@%s/%s", user, pass, addr, dbname)
//	if user == "" && pass == "" {
//		url = fmt.Sprintf("mongodb://%s/%s", addr, dbname)
//	}
//	cfg.Logf("Account mongodb connect url:%s\n", url)
//	var err error
//	accountSession, err = mgo.Dial(url)
//	if err != nil {
//		panic(err.Error())
//	}

//	accountSession.SetMode(mgo.Monotonic, true)
//	cfg.Logf("connect to account mongodb...OK")

//	accountSession.DB("").C(ACCOUNT_COLLECTION).EnsureIndexKey("id")
//}

//func m2DbServer(m bson.M) *DbServer {
//	s := new(DbServer)
//	s.Id = int16(m["id"].(float64))
//	s.Population = int16(m["population"].(float64))
//	s.OpenTime = int64(m["opentime"].(float64))

//	return s
//}

//func GetServerList() []*DbServer {
//	cond := bson.M{}

//	query, err := findAccount(SERVER_COLLECTION, cond)
//	if err != nil {
//		return nil
//	}
//	servers := make([]*DbServer, 0)
//	iter := query.Iter()
//	m := make(bson.M)
//	for iter.Next(m) == true {
//		s := m2DbServer(m)
//		servers = append(servers, s)
//	}

//	return servers
//}

//func UpdateServerPopulation(id int16, population int16) {
//	if accountSession == nil {
//		return
//	}

//	cond := bson.M{"id": id}
//	doc := bson.M{
//		"population": float64(population),
//	}

//	err := updateAccount(SERVER_COLLECTION, cond, doc)
//	if err != nil {
//		cfg.LogErrf("update server pupulation %v fail %v", id, err.Error())
//	}
//}

//func GetDefaultGs() int16 {
//	cond := bson.M{"config_type": CONFIG_TYPE_DEFAULT}

//	q, err := findAccount(CONFIG_COLLECTION, cond)
//	if err != nil {
//		return 1
//	}

//	m := make(bson.M)
//	if err := q.One(m); err != nil {
//		return 1
//	}
//	return int16(m["default_gs"].(int))
//}

//func UpdateDefaultGs(gsid int16) {
//	if accountSession == nil {
//		return
//	}

//	cond := bson.M{"config_type": CONFIG_TYPE_DEFAULT}
//	doc := bson.M{
//		"config_type": CONFIG_TYPE_DEFAULT,
//		"default_gs":  gsid,
//	}

//	err := updateAccount(CONFIG_COLLECTION, cond, doc)
//	if err != nil {
//		cfg.LogFatalf("update default gs %v fail %v", gsid, err.Error())
//	}
//}

//func UpdateAccount(userId int64, account DbAccount) {
//	if accountSession == nil {
//		return
//	}

//	cond := bson.M{"id": account.Id}
//	accountDoc := bson.M{
//		"id":                 account.Id,
//		"server_id":          account.ServerId,
//		"udid":               account.Udid,
//		"account":            account.Account,
//		"account_type":       account.AccountType,
//		"register_server_id": account.RegisterServerId,
//	}

//	err := updateAccount(ACCOUNT_COLLECTION, cond, accountDoc)
//	if err != nil {
//		cfg.LogFatalf("update account %v fail %v", account, err.Error())
//	}
//}

//func GetAccountById(userId int64) *DbAccount {
//	return getAccount(bson.M{"id": userId})
//}

//func GetAccountByUdid(udid string) *DbAccount {
//	return getAccount(bson.M{"udid": udid})
//}

////通过账号类型、账号名查找
//func GetAccountByName(name string, atype int16) *DbAccount {
//	return getAccount(bson.M{"account": name, "account_type": atype})
//}

////账号系统中获取某个server最大的userid
//func getServerMaxId(name string, serverid int16) (int64, error) {
//	if accountSession == nil {
//		return 0, errors.New("invalid accountSession")
//	}
//	c := accountSession.DB("").C(name)
//	q := c.Find(bson.M{"server_id": serverid}).Sort("-id").Limit(1)
//	m := make(bson.M)
//	if err := q.One(m); err != nil {
//		return 0, nil
//	}
//	id, present := m["id"]
//	if !present {
//		return 0, nil
//	}
//	return id.(int64), nil
//}

////账号系统中获取某个server最大的userid
//func getMaxId(name string) (int64, error) {
//	if accountSession == nil {
//		return 0, errors.New("invalid accountSession")
//	}
//	c := accountSession.DB("").C(name)
//	q := c.Find(nil).Sort("-id").Limit(1)
//	m := make(bson.M)
//	if err := q.One(m); err != nil {
//		return 0, nil
//	}
//	id, present := m["id"]
//	if !present {
//		return 0, nil
//	}
//	return id.(int64), nil
//}

//func getAccount(cond interface{}) *DbAccount {
//	q, err := findAccount(ACCOUNT_COLLECTION, cond)
//	if err != nil {
//		return nil
//	}

//	m := make(bson.M)
//	if err := q.One(m); err != nil {
//		return nil
//	}

//	return &DbAccount{
//		Id:               m["id"].(int64),
//		ServerId:         int16(m["server_id"].(int)),
//		Udid:             m["udid"].(string),
//		Account:          m["account"].(string),
//		AccountType:      int16(m["account_type"].(int)),
//		RegisterServerId: int16(m["register_server_id"].(int)),
//	}
//}

//func findAccount(name string, query interface{}) (*mgo.Query, error) {
//	if accountSession == nil {
//		return nil, errors.New("Invalid account session.")
//	}
//	c := accountSession.DB("").C(name)
//	return c.Find(query), nil
//}

//func updateAccount(name string, cond interface{}, doc interface{}) error {
//	if accountSession == nil {
//		return errors.New("Invalid account session.")
//	}

//	c := accountSession.DB("").C(name)
//	_, err := c.Upsert(cond, bson.M{"$set": doc})
//	return err
//}
