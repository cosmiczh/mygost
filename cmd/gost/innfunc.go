package main

import (
	"net"
	"strconv"
)

func MaxInt(i1, i2 int) int {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
func MinInt(i1, i2 int) int {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}
func MaxInt64(i1, i2 int64) int64 {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
func MinInt64(i1, i2 int64) int64 {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}

func MaxInt32(i1, i2 int32) int32 {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
func MinInt32(i1, i2 int32) int32 {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}

func MaxUInt32(i1, i2 uint32) uint32 {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
func MinUInt32(i1, i2 uint32) uint32 {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}
func MaxInt16(i1, i2 int16) int16 {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
func MinInt16(i1, i2 int16) int16 {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}
func MaxUInt8(i1, i2 uint8) uint8 {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
func MinUInt8(i1, i2 uint8) uint8 {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}
func BetweenInt(my int, start, end int) bool          { return my >= start && my < end }
func BetweenUInt(my uint, start, end uint) bool       { return my >= start && my < end }
func BetweenInt8(my int8, start, end int8) bool       { return my >= start && my < end }
func BetweenUInt8(my uint8, start, end uint8) bool    { return my >= start && my < end }
func BetweenByte(my byte, start, end byte) bool       { return my >= start && my < end }
func BetweenInt16(my int16, start, end int16) bool    { return my >= start && my < end }
func BetweenUInt16(my uint16, start, end uint16) bool { return my >= start && my < end }
func BetweenInt32(my int32, start, end int32) bool    { return my >= start && my < end }
func BetweenUInt32(my uint32, start, end uint32) bool { return my >= start && my < end }
func BetweenInt64(my int64, start, end int64) bool    { return my >= start && my < end }
func BetweenUInt64(my uint64, start, end uint64) bool { return my >= start && my < end }
func IFbool(cond bool, true_bool bool, false_bool bool) bool {
	if cond {
		return true_bool
	} else {
		return false_bool
	}
}
func IFint(cond bool, true_int int, false_int int) int {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFint8(cond bool, true_int int8, false_int int8) int8 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFint16(cond bool, true_int int16, false_int int16) int16 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFint32(cond bool, true_int int32, false_int int32) int32 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFint64(cond bool, true_int int64, false_int int64) int64 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}

func IFunt(cond bool, true_int uint, false_int uint) uint {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFunt8(cond bool, true_int uint8, false_int uint8) uint8 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFunt16(cond bool, true_int uint16, false_int uint16) uint16 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFunt32(cond bool, true_int uint32, false_int uint32) uint32 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFunt64(cond bool, true_int uint64, false_int uint64) uint64 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func IFstr(cond bool, true_str string, false_str string) string {
	if cond {
		return true_str
	} else {
		return false_str
	}
}
func IFstrs(cond bool, true_str []string, false_str []string) []string {
	if cond {
		return true_str
	} else {
		return false_str
	}
}

type Any interface{}

func IFany(cond bool, true_any Any, false_any Any) Any {
	if cond {
		return true_any
	} else {
		return false_any
	}
}
func FFAny(any Any) func() Any { return func() Any { return any } }
func IFFany(cond bool, true_any func() Any, false_any func() Any) func() Any {
	if cond {
		return true_any
	} else {
		return false_any
	}
}
func IFFfunc(cond bool, true_func func(), false_func func()) {
	if cond {
		true_func()
	} else {
		false_func()
	}
}
func Fnbool(tf bool) func() bool { return func() bool { return tf } }
func IFFbool(cond bool, true_bool func() bool, false_bool func() bool) func() bool {
	if cond {
		return true_bool
	} else {
		return false_bool
	}
}
func Fnint(i int) func() int { return func() int { return i } }
func IFFint(cond bool, true_int func() int, false_int func() int) func() int {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnint8(i int8) func() int8 { return func() int8 { return i } }
func IFFint8(cond bool, true_int func() int8, false_int func() int8) func() int8 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnint16(i int16) func() int16 { return func() int16 { return i } }
func IFFint16(cond bool, true_int func() int16, false_int func() int16) func() int16 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnint32(i int32) func() int32 { return func() int32 { return i } }
func IFFint32(cond bool, true_int func() int32, false_int func() int32) func() int32 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnint64(i int64) func() int64 { return func() int64 { return i } }
func IFFint64(cond bool, true_int func() int64, false_int func() int64) func() int64 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}

func Fnunt(i uint) func() uint { return func() uint { return i } }
func IFFunt(cond bool, true_int func() uint, false_int func() uint) func() uint {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnunt8(i uint8) func() uint8 { return func() uint8 { return i } }
func IFFunt8(cond bool, true_int func() uint8, false_int func() uint8) func() uint8 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnunt16(i uint16) func() uint16 { return func() uint16 { return i } }
func IFFunt16(cond bool, true_int func() uint16, false_int func() uint16) func() uint16 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnunt32(i uint32) func() uint32 { return func() uint32 { return i } }
func IFFunt32(cond bool, true_int func() uint32, false_int func() uint32) func() uint32 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}
func Fnunt64(i uint64) func() uint64 { return func() uint64 { return i } }
func IFFunt64(cond bool, true_int func() uint64, false_int func() uint64) func() uint64 {
	if cond {
		return true_int
	} else {
		return false_int
	}
}

func Fnstr(s string) func() string { return func() string { return s } }
func IFFstr(cond bool, true_str func() string, false_str func() string) func() string {
	if cond {
		return true_str
	} else {
		return false_str
	}
}
func Fnstrs(s []string) func() []string { return func() []string { return s } }
func IFFstrs(cond bool, true_str func() []string, false_str func() []string) func() []string {
	if cond {
		return true_str
	} else {
		return false_str
	}
}

func FindInt(arr []int, val int) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			return i
		}
	}
	return -1
}
func FindInt32(arr []int32, val int32) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			return i
		}
	}
	return -1
}
func FindInt64(arr []int64, val int64) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			return i
		}
	}
	return -1
}
func FindString(arr []string, val string) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			return i
		}
	}
	return -1
}

func InSet(a int, set ...int) bool {
	for i := 0; i < len(set); i++ {
		if a == set[i] {
			return true
		}
	}
	return false
}
func IPv4itoa(IPi uint32) string {
	l_c2sip := make(net.IP, net.IPv4len)
	for i := len(l_c2sip) - 1; i >= 0; i-- {
		l_c2sip[i] = byte(IPi % 256)
		IPi /= 256
	}
	return l_c2sip.String()
}
func I16toa(i int16) string { return strconv.FormatInt(int64(i), 10) }
func I64toa(i int64) string { return strconv.FormatInt(i, 10) }
func Itoa(i int32) string   { return strconv.FormatInt(int64(i), 10) }
func Iutoa(i uint32) string { return strconv.FormatInt(int64(i), 10) }

func Atoiu(s string) (uint32, error) {
	i64, err := strconv.ParseInt(s, 0, 64)
	if nerr, ok := err.(*strconv.NumError); ok {
		nerr.Func = "Atoi32"
	}
	return uint32(uint64(i64)), err
}
func Atoi(s string) (int32, error) {
	i64, err := strconv.ParseInt(s, 0, 32)
	if nerr, ok := err.(*strconv.NumError); ok {
		nerr.Func = "Atoi32"
	}
	return int32(i64), err
}

func UnsafeAtoi(s string) int32 {
	i64, _ := strconv.ParseInt(s, 0, 32)
	return int32(i64)
}

func Atoi64(s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 0, 64)
	if nerr, ok := err.(*strconv.NumError); ok {
		nerr.Func = "Atoi64"
	}
	return i64, err
}
func Atobool(s string) (bool, error) {
	switch s {
	case "true", "TRUE", "True":
		return true, nil
	case "false", "FALSE", "False":
		return false, nil
	}
	return false, &strconv.NumError{"ParseBool", s, strconv.ErrSyntax}
}

func UnsafeAtoi64(s string) int64 {
	i64, _ := strconv.ParseInt(s, 0, 64)
	return i64
}

func Atoi16(s string) (int16, error) {
	i64, err := strconv.ParseInt(s, 0, 16)
	if nerr, ok := err.(*strconv.NumError); ok {
		nerr.Func = "Atoi16"
	}
	return int16(i64), err
}

func AnyToString(ipara interface{}) string {
	l_sparam := ""
	switch v := ipara.(type) {
	case string:
		l_sparam = v
	case int:
		l_sparam = strconv.FormatInt(int64(v), 10)
	case int8:
		l_sparam = strconv.FormatInt(int64(v), 10)
	case int32:
		l_sparam = strconv.FormatInt(int64(v), 10)
	case uint32:
		l_sparam = strconv.FormatUint(uint64(v), 10)
	case int64:
		l_sparam = strconv.FormatInt(int64(v), 10)
	case uint64:
		l_sparam = strconv.FormatUint(uint64(v), 10)
	case float32:
		l_sparam = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		l_sparam = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return l_sparam
}

func AnyToInt32(ipara interface{}) int32 {
	l_i32val := int32(0)
	switch v := ipara.(type) {
	case string:
		l_i32val = UnsafeAtoi(string(v))
	case int:
		l_i32val = int32(v)
	case int8:
		l_i32val = int32(v)
	case int32:
		l_i32val = int32(v)
	case int64:
		l_i32val = int32(v)
	case float32:
		l_i32val = int32(v)
	case float64:
		l_i32val = int32(v)
	}
	return l_i32val
}
