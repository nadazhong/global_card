package gstimer

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

var (
	GS_TIMER  chan int64
	TIMER_MAP map[int64]*GSTimer
	_map      map[int64]*Timer
)

func init() {
	TIMER_MAP = make(map[int64]*GSTimer)
	_map = make(map[int64]*Timer)
	GS_TIMER = make(chan int64, 2048) // FIXME:定时器超过缓冲区会卡住。
}

func expire(timerId int64) {
	GS_TIMER <- timerId
}

func FreeTimer(timerId int64) {
	delete(_map, timerId)
	delete(TIMER_MAP, timerId)
}

// 创建定时器。delay 的单位为毫秒。
func CreateTimer(userId int64, delay int64, msg Msg) int64 {
	if delay < 0 {
		delay = 0
	}
	tr := NewTimer(delay * 1e6)
	TIMER_MAP[tr.Id] = &GSTimer{Receiver: userId, Msg: msg}
	_map[tr.Id] = tr

	return tr.Id
}

func CancelTimer(timerId int64) {
	tr, present := _map[timerId]
	if present {
		if !tr.Stop() {
			fmt.Println("Stop timer fail:", timerId)
		}
	}
	FreeTimer(timerId)
}

func ModifyTimer(timerId int64, delta int64) {
	tr, present := _map[timerId]
	if present {
		tr.Modify(delta * 1e6)
	}
}

/// marshal/unmarshal timers
type DbTimer struct {
	Receiver int64
	Data     Msg
	Expired  int64
	Id       int64
}

type TimerInfo struct {
	Timers []DbTimer
	NextId int64
}

func LoadTimers(data []byte) (error, func()) {
	info := &TimerInfo{
		Timers: make([]DbTimer, 0),
		NextId: 0,
	}
	enc := gob.NewDecoder(bytes.NewReader(data))
	if err := enc.Decode(info); err != nil {
		fmt.Println("Load GS_MAP error:", err.Error())
		return err, nil
	}

	// reset next timer id
	timerMutex.Lock()
	timerNextId = info.NextId
	timerMutex.Unlock()
	fmt.Println("timerNextId=%d\n", timerNextId)

	return nil, func() {
		rescheduleTimers(info)
	}
}

func rescheduleTimers(info *TimerInfo) error {
	/*
		timerMutex.Lock()
		timerNextId = info.NextId
		timerMutex.Unlock()

		fmt.Println("timerNextId=%d\n", timerNextId)
	*/
	for _, item := range info.Timers {
		timer := &Timer{
			C:  make(chan int64, 1),
			t:  item.Expired,
			i:  -1,
			Id: item.Id,
			f:  expire,
		}
		reschedule(timer)
		TIMER_MAP[timer.Id] = &GSTimer{Receiver: item.Receiver, Msg: item.Data}
		_map[timer.Id] = timer
		fmt.Printf("reschedule timer tid=%d, ownerId=%d\n", timer.Id, item.Receiver)
	}
	return nil
}

func DumpTimers() (error, []byte) {
	info := TimerInfo{
		Timers: make([]DbTimer, 0),
		NextId: timerNextId,
	}

	for tid, timer := range _map {
		gstimer, present := TIMER_MAP[tid]
		if present && gstimer.Receiver > 0 {
			entry := DbTimer{
				Receiver: gstimer.Receiver,
				Data:     gstimer.Msg,
				Expired:  timer.t,
				Id:       timer.Id,
			}
			info.Timers = append(info.Timers, entry)
			fmt.Printf("trying save timer id=%d\n", timer.Id)
		}
	}

	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	if err := enc.Encode(info); err != nil {
		return err, nil
	}
	return nil, buffer.Bytes()
}
