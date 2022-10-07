package main

import (
	"bytes"
)

/*-----------------------------------------------------------
class MAPInt
*/
type MAPInt struct {
	MAP
}

func (this *MAPInt) MAPInt(canRepeat bool) *MAPInt {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
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
func (this *MAPInt) Insert(key int, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPInt) Remove(key int) int   { return this.MAP.Remove(key) }
func (this *MAPInt) Find(key int) MAPIter { return this.MAP.Find(key) }

/*-----------------------------------------------------------
class MAPInt8
*/
type MAPInt8 struct {
	MAP
}

func (this *MAPInt8) MAPInt8(canRepeat bool) *MAPInt8 {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
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
func (this *MAPInt8) Insert(key int8, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPInt8) Remove(key int8) int   { return this.MAP.Remove(key) }
func (this *MAPInt8) Find(key int8) MAPIter { return this.MAP.Find(key) }

/*-----------------------------------------------------------
class MAPInt16
*/
type MAPInt16 struct {
	MAP
}

func (this *MAPInt16) MAPInt16(canRepeat bool) *MAPInt16 {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
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
func (this *MAPInt16) Insert(key int16, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPInt16) Remove(key int16) int   { return this.MAP.Remove(key) }
func (this *MAPInt16) Find(key int16) MAPIter { return this.MAP.Find(key) }

/*-----------------------------------------------------------
class MAPInt32
*/
type MAPInt32 struct {
	MAP
}

func (this *MAPInt32) MAPInt32(canRepeat bool) *MAPInt32 {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
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
func (this *MAPInt32) Insert(key int32, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPInt32) Remove(key int32) int   { return this.MAP.Remove(key) }
func (this *MAPInt32) Find(key int32) MAPIter { return this.MAP.Find(key) }

/*-----------------------------------------------------------
class MAPInt64
*/
type MAPInt64 struct {
	MAP
}

func (this *MAPInt64) MAPInt64(canRepeat bool) *MAPInt64 {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
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
func (this *MAPInt64) Insert(key int64, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPInt64) Remove(key int64) int   { return this.MAP.Remove(key) }
func (this *MAPInt64) Find(key int64) MAPIter { return this.MAP.Find(key) }

/*-----------------------------------------------------------
class MAPstr
*/
type MAPstr struct {
	MAP
}

func (this *MAPstr) MAPstr(canRepeat bool) *MAPstr {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
		l_key1, l_key2 := key1.(string), key2.(string)
		return int8(bytes.Compare([]byte(l_key1), []byte(l_key2)))
	})
	return this
}
func (this *MAPstr) Insert(key string, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPstr) Remove(key string) int   { return this.MAP.Remove(key) }
func (this *MAPstr) Find(key string) MAPIter { return this.MAP.Find(key) }

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

type MAPistr struct {
	MAP
}

func (this *MAPistr) MAPistr(canRepeat bool) *MAPistr {
	this.MAP.MAP(canRepeat, func(key1, key2 interface{}) int8 {
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
func (this *MAPistr) Insert(key string, value interface{}) (sucs bool, it MAPIter) {
	return this.MAP.Insert(key, value)
}
func (this *MAPistr) Remove(key string) int   { return this.MAP.Remove(key) }
func (this *MAPistr) Find(key string) MAPIter { return this.MAP.Find(key) }
