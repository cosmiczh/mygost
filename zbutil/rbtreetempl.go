package zbutil

import (
	"bytes"
)

/*-----------------------------------------------------------
class RBtreeInt
*/
type RBtreeInt struct {
	RBtree
}

func (this *RBtreeInt) RBtreeInt(canRepeat bool) *RBtreeInt {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(int), key2.(int)
		if l_key1 > l_key2 {
			return 1
		} else if l_key1 < l_key2 {
			return -1
		} else {
			return 0
		}
	})
	return this
}
func (this *RBtreeInt) Insert(key int, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeInt) Remove(key int) int      { return this.RBtree.Remove(key) }
func (this *RBtreeInt) Find(key int) RBtreeIter { return this.RBtree.Find(key) }

/*-----------------------------------------------------------
class RBtreeInt8
*/
type RBtreeInt8 struct {
	RBtree
}

func (this *RBtreeInt8) RBtreeInt8(canRepeat bool) *RBtreeInt8 {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(int8), key2.(int8)
		if l_key1 > l_key2 {
			return 1
		} else if l_key1 < l_key2 {
			return -1
		} else {
			return 0
		}
	})
	return this
}
func (this *RBtreeInt8) Insert(key int8, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeInt8) Remove(key int8) int      { return this.RBtree.Remove(key) }
func (this *RBtreeInt8) Find(key int8) RBtreeIter { return this.RBtree.Find(key) }

/*-----------------------------------------------------------
class RBtreeInt16
*/
type RBtreeInt16 struct {
	RBtree
}

func (this *RBtreeInt16) RBtreeInt16(canRepeat bool) *RBtreeInt16 {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(int16), key2.(int16)
		if l_key1 > l_key2 {
			return 1
		} else if l_key1 < l_key2 {
			return -1
		} else {
			return 0
		}
	})
	return this
}
func (this *RBtreeInt16) Insert(key int16, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeInt16) Remove(key int16) int      { return this.RBtree.Remove(key) }
func (this *RBtreeInt16) Find(key int16) RBtreeIter { return this.RBtree.Find(key) }

/*-----------------------------------------------------------
class RBtreeInt32
*/
type RBtreeInt32 struct {
	RBtree
}

func (this *RBtreeInt32) RBtreeInt32(canRepeat bool) *RBtreeInt32 {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(int32), key2.(int32)
		if l_key1 > l_key2 {
			return 1
		} else if l_key1 < l_key2 {
			return -1
		} else {
			return 0
		}
	})
	return this
}
func (this *RBtreeInt32) Insert(key int32, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeInt32) Remove(key int32) int      { return this.RBtree.Remove(key) }
func (this *RBtreeInt32) Find(key int32) RBtreeIter { return this.RBtree.Find(key) }

/*-----------------------------------------------------------
class RBtreeInt64
*/
type RBtreeInt64 struct {
	RBtree
}

func (this *RBtreeInt64) RBtreeInt64(canRepeat bool) *RBtreeInt64 {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(int64), key2.(int64)
		if l_key1 > l_key2 {
			return 1
		} else if l_key1 < l_key2 {
			return -1
		} else {
			return 0
		}
	})
	return this
}
func (this *RBtreeInt64) Insert(key int64, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeInt64) Remove(key int64) int      { return this.RBtree.Remove(key) }
func (this *RBtreeInt64) Find(key int64) RBtreeIter { return this.RBtree.Find(key) }

//RBtreeStr Key字符串区分大小写
type RBtreeStr struct {
	RBtree
}

func (this *RBtreeStr) RBtreeStr(canRepeat bool) *RBtreeStr {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(string), key2.(string)
		return int8(bytes.Compare([]byte(l_key1), []byte(l_key2)))
	})
	return this
}
func (this *RBtreeStr) Insert(key string, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeStr) Remove(key string) int      { return this.RBtree.Remove(key) }
func (this *RBtreeStr) Find(key string) RBtreeIter { return this.RBtree.Find(key) }

/*-----------------------------------------------------------
class MAPistr
*/
func betweenuint8(my uint8, start, end uint8) bool { return my >= start && my < end }
func minint(i1, i2 int) int {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}

//RBtreeIstr Key字符串不区分大小写
type RBtreeIstr struct {
	RBtree
}

func (this *RBtreeIstr) RBtreeIstr(canRepeat bool) *RBtreeIstr {
	this.RBtree.RBtree(canRepeat, func(key1, key2 interface{}) int8 {
		l_s1, l_s2 := key1.(string), key2.(string)
		l_len := minint(len(l_s1), len(l_s2))
		for i := 0; i < l_len; i++ {
			l_c1, l_c2 := l_s1[i], l_s2[i]
			if betweenuint8(l_c1, 'a', 'z'+1) {
				l_c1 -= ('a' - 'A')
			}
			if betweenuint8(l_c2, 'a', 'z'+1) {
				l_c2 -= ('a' - 'A')
			}
			if l_c1 > l_c2 {
				return 1
			} else if l_c1 < l_c2 {
				return -1
			}
		}
		if len(l_s1) > l_len {
			return 1
		} else if len(l_s2) > l_len {
			return -1
		} else {
			return 0
		}
	})
	return this
}
func (this *RBtreeIstr) Insert(key string, value interface{}) (sucs bool, it RBtreeIter) {
	return this.RBtree.Insert(key, value)
}
func (this *RBtreeIstr) Remove(key string) int      { return this.RBtree.Remove(key) }
func (this *RBtreeIstr) Find(key string) RBtreeIter { return this.RBtree.Find(key) }
