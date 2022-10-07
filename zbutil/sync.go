package zbutil

import (
	"sync"
	"sync/atomic"
)

/*---------------type Mutex struct-------------------------------------------------------------------------
标准互斥锁封装，不支持递归(嵌套)锁
*/
type Mutex struct {
	m_mtx sync.Mutex
}

func (this *Mutex) Lock() (Unlock func()) {
	this.m_mtx.Lock()
	l_lock := true
	return func() {
		if !l_lock {
			return
		}
		l_lock = false
		this.m_mtx.Unlock()
	}
}
func (this *Mutex) LockFun(fun func()) {
	defer this.Lock()()
	fun()
}

/*---------------type RWMutex2 struct-------------------------------------------------------------------------
简单实现的读写锁2，支持递归读、不支持递归写
*/
type RWMutex2 struct {
	mtx  Mutex
	cond *sync.Cond
	init uint32

	rc, wc, pwc int
}

func (this *RWMutex2) TryLock() (Unlock func()) { return this.lock(true) }
func (this *RWMutex2) Lock() (Unlock func())    { return this.lock(false) }

func (this *RWMutex2) lock(try bool) (Unlock func()) {
	defer this.mtx.Lock()()
	if this.init == 0 {
		this.cond = sync.NewCond(&this.mtx.m_mtx)
		atomic.StoreUint32(&this.init, 1)
	}
	this.pwc++
	for this.wc != 0 || this.rc != 0 {
		if !try {
			this.cond.Wait()
		} else {
			this.pwc--
			return nil
		}
	}
	this.pwc--
	this.wc++
	l_lock := true
	return func() {
		if !l_lock {
			return
		}
		l_lock = false
		defer this.mtx.Lock()()
		this.wc--
		this.cond.Broadcast()
	}
}
func (this *RWMutex2) RLock() (RUnlock func()) {
	defer this.mtx.Lock()()
	if this.init == 0 {
		this.cond = sync.NewCond(&this.mtx.m_mtx)
		atomic.StoreUint32(&this.init, 1)
	}
	for this.wc != 0 || (this.rc == 0 && this.pwc != 0) {
		this.cond.Wait()
	}
	this.rc++
	l_lock := true
	return func() {
		if !l_lock {
			return
		}
		l_lock = false
		defer this.mtx.Lock()()
		if this.rc--; this.rc == 0 {
			this.cond.Signal()
		}
	}
}
func (this *RWMutex2) LockFun(fun func()) {
	defer this.Lock()()
	fun()
}
func (this *RWMutex2) RLockFun(fun func()) {
	defer this.RLock()()
	fun()
}
