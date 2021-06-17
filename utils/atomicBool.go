package utils

import (
	"sync/atomic"
)

type AtomicBool struct {
	flag int32
}

func (a *AtomicBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(a.flag), i)
}

func (a *AtomicBool) Get() bool {
	if atomic.LoadInt32(&(a.flag)) != 0 {
		return true
	}
	return false
}
