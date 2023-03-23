package runtime

import "sync/atomic"

type rwmutex struct {
	rLock      mutex    // protects readers, readerPass, writer
	readers    muintptr // list of pending readers
	readerPass uint32   // number of pending readers to skip readers list

	wLock  mutex    // serializes writers
	writer muintptr // pending writer waiting for completing readers

	readerCount atomic.Int32 // number of pending readers
	readerWait  atomic.Int32 // number of departing readers
}

func (rw *rwmutex) rlock()   {}
func (rw *rwmutex) runlock() {}
func (rw *rwmutex) lock()    {}
func (rw *rwmutex) unlock() {
}
