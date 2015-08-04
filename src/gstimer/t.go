// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 修改以适应 gs 的定时器需求。

package gstimer

import (
	"container/heap"
	"sync"
	"time"
)

// Nanoseconds reports the number of nanoseconds since the Unix epoch,
// January 1, 1970 00:00:00 UTC.
func Nanoseconds() int64 {
	return time.Now().UnixNano()
}

// The Timer type represents a single event.
// When the Timer expires, the (timer id) will be sent on C
// unless the Timer represents an AfterFunc event.
type Timer struct {
	C  <-chan int64
	t  int64       // The absolute time that the event should fire.
	f  func(int64) // The function to call when the event fires.
	i  int         // The event's index inside eventHeap.
	Id int64
}

type timerHeap []*Timer

// forever is the absolute time (in ns) of an event that is forever away.
const forever = 1 << 62

// maxSleepTime is the maximum length of time that a sleeper
// sleeps for before checking if it is defunct.
const maxSleepTime = 1e9

var (
	// timerMutex guards the variables inside this var group.
	timerMutex sync.Mutex

	// timers holds a binary heap of pending events, terminated with a sentinel.
	timers timerHeap

	// currentSleeper is an ever-incrementing counter which represents
	// the current sleeper. It allows older sleepers to detect that they are
	// defunct and exit.
	currentSleeper int64

	// ID计数器。
	timerNextId int64
)

func init() {
	timerNextId = 1
	timers = timerHeap{&Timer{t: forever}}
	heap.Init(timers)
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least ns nanoseconds.
func NewTimer(ns int64) *Timer {
	c := make(chan int64, 1)
	e := after(ns, expire)
	e.C = c
	return e
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or stopped.
func (e *Timer) Stop() (ok bool) {
	timerMutex.Lock()
	// Avoid removing the first event in the queue so that
	// we don't start a new sleeper unnecessarily.
	if e.i > 0 {
		heap.Remove(timers, e.i)
	}
	ok = e.f != nil
	e.f = nil
	timerMutex.Unlock()
	return true
}

// Modify the Timer with new delay value, need go 1.2
func (e *Timer) Modify(delta int64) {

	timerMutex.Lock()
	e.t = e.t + delta
	t0 := timers[0].t
	heap.Fix(timers, e.i)
	if t0 > e.t && (t0 == forever) {
		currentSleeper++
		go sleeper(currentSleeper)
	}
	timerMutex.Unlock()
}

// after is the implementation of After and AfterFunc.
// When the current time is after ns, it calls f with the current time.
// It assumes that f will not block.
func after(ns int64, f func(int64)) (e *Timer) {
	now := Nanoseconds()
	t := now + ns
	if ns > 0 && t < now {
		panic("time: time overflow")
	}
	timerMutex.Lock()
	t0 := timers[0].t
	e = &Timer{nil, t, f, -1, timerNextId}
	timerNextId++
	heap.Push(timers, e)
	// Start a new sleeper if the new event is before
	// the first event in the queue. If the length of time
	// until the new event is at least maxSleepTime,
	// then we're guaranteed that the sleeper will wake up
	// in time to service it, so no new sleeper is needed.
	if t0 > t && (t0 == forever || ns < maxSleepTime) {
		currentSleeper++
		go sleeper(currentSleeper)
	}
	timerMutex.Unlock()
	return
}

func reschedule(t *Timer) {
	now := Nanoseconds()
	timerMutex.Lock()
	t0 := timers[0].t
	if now > t.t {
		t.t = now
	}
	ns := t.t - now
	heap.Push(timers, t)
	if t0 > t.t && (t0 == forever || ns < maxSleepTime) {
		currentSleeper++
		go sleeper(currentSleeper)
	}
	timerMutex.Unlock()
}

// sleeper continually looks at the earliest event in the queue, waits until it happens,
// then removes any events in the queue that are due. It stops when the queue
// is empty or when another sleeper has been started.
func sleeper(sleeperId int64) {
	timerMutex.Lock()
	e := timers[0]
	t := Nanoseconds()
	for e.t != forever {
		if dt := e.t - t; dt > 0 {
			if dt > maxSleepTime {
				dt = maxSleepTime
			}
			timerMutex.Unlock()
			time.Sleep(time.Duration(dt))
			timerMutex.Lock()
			if currentSleeper != sleeperId {
				// Another sleeper has been started, making this one redundant.
				break
			}
		}
		e = timers[0]
		t = Nanoseconds()
		for t >= e.t {
			if e.f != nil {
				e.f(e.Id)
				e.f = nil
			}
			heap.Pop(timers)
			e = timers[0]
		}
	}
	timerMutex.Unlock()
}

func (timerHeap) Len() int {
	return len(timers)
}

func (timerHeap) Less(i, j int) bool {
	return timers[i].t < timers[j].t
}

func (timerHeap) Swap(i, j int) {
	timers[i], timers[j] = timers[j], timers[i]
	timers[i].i = i
	timers[j].i = j
}

func (timerHeap) Push(x interface{}) {
	e := x.(*Timer)
	e.i = len(timers)
	timers = append(timers, e)
}

func (timerHeap) Pop() interface{} {
	// TODO: possibly shrink array.
	n := len(timers) - 1
	e := timers[n]
	timers[n] = nil
	timers = timers[0:n]
	e.i = -1
	return e
}
