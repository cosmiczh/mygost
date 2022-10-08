package zbutil

import "sync"

type ChkOnce struct {
	once_armor map[string]bool
	once_lock  sync.RWMutex
}

func NewChkOnce() *ChkOnce {
	return &ChkOnce{
		once_armor: make(map[string]bool, 10),
	}
}
func (this *ChkOnce) CheckFirst(key string) bool {
	this.once_lock.RLock()
	if this.once_armor[key] {
		this.once_lock.RUnlock()
		return false
	}
	this.once_lock.RUnlock()

	this.once_lock.Lock()
	defer this.once_lock.Unlock()
	if this.once_armor[key] { //再次检查是否已存在
		return false
	}
	this.once_armor[key] = true
	return true
}
func (this *ChkOnce) Clear() {
	this.once_lock.Lock()
	defer this.once_lock.Unlock()
	this.once_armor = make(map[string]bool, len(this.once_armor)+10)
}
