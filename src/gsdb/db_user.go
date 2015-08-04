package gsdb

import (
	"fmt"
)

type DbUser struct {
	Id    int64
	Name  string
	Level int16
	Exp   int32
	Icon  int16

	RenameCnt int16 // 免费的改名次数

	Diamond     int32
	Cash        float64
	CashLimit   int32
	Energy      float64
	EnergyLimit int32
	Food        float64
	FoodLimit   int32
	Oil         float64
	OilLimit    int32
	Steel       float64
	SteelLimit  int32

	ResSyncTime   int64
	CashSpeed     int32
	BuildSpeed    int32
	DefBuildSpeed int32
	EnergySpeed   int32
	FoodSpeed     int32
	OilSpeed      int32
	SteelSpeed    int32
	FoodLoad      int32

	Alliance             int32
	JoinAllianceTime     int64
	Credit               int32 // 公会币。
	CreditLimit          int32 // 每天公会币上限
	FreshCreditLimitTime int64 // 下次每天公会币上限刷新时间

	Speed int32

	BuildingPower int32 // 建筑加成
	ResearchPower int32 // 科技加成
	TroopPower    int32 // 军队加成
	GarrisonPower int32 // 城防部队加成
	QuestPower    int32 // 任务加成
	Power         int32 // 战斗力
	KillCnt       int32 // 杀人数
	PowerRank     int32

	BuildingFreeTime int64

	LastSavingTime int64 // 最后存盘时间
	LastOpTime     int64 // 最后操作时间

	LastConLoginTime int64 // 最近一次记录连续登录的时间
	ConLoginDays     int32 // 连续登录天数
	LoginTime        int64 // 用户登录时间

	TimerId map[int16]int64 // 玩家保护盾定时器

	SaveData string

	AtkStyle [4]int16 //战斗默认阵型
	DefStyle [4]int16

	Udid string
	Ip   string
}

func (user *DbUser) String() string {
	report := "\n=====================================\n"
	report += fmt.Sprintf("user[%d]={\n", user.Id)
	//s := reflect.ValueOf(user).Elem()
	//typeOfT := s.Type()
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	report += fmt.Sprintf("\t%s=%v\n", typeOfT.Field(i).Name, f.Interface())
	//}
	report += fmt.Sprintf("}\n")
	return report
}
