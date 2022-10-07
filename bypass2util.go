package gost

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
)

type struct4sort struct {
	m_len   int
	mf_less func(i, j int) bool
	mf_swap func(i, j int)
}

func (this *struct4sort) Len() int {
	return this.m_len
}
func (this *struct4sort) Less(i, j int) bool {
	return this.mf_less(i, j)
}
func (this *struct4sort) Swap(i, j int) {
	this.mf_swap(i, j)
}

func Sort(nlen int, fless func(i, j int) bool, fswap func(i, j int)) {
	sort.Sort(&struct4sort{m_len: nlen, mf_less: fless, mf_swap: fswap})
}
func Reverse(nlen int, fswap func(i, j int)) {
	for i := 0; i < nlen/2; i++ {
		fswap(i, nlen-1-i)
	}
}
func CountFunc(nlen int, cond func(i int) bool) int {
	l_count := 0
	for i := 0; i < nlen; i++ {
		if cond(i) {
			l_count++
		}
	}
	return l_count
}
func IndexFunc(nlen int, cond func(i int) bool) int {
	for i := 0; i < nlen; i++ {
		if cond(i) {
			return i
		}
	}
	return -1
}
func RevIndexFunc(nlen int, cond func(i int) bool) int {
	for i := nlen - 1; i >= 0; i-- {
		if cond(i) {
			return i
		}
	}
	return -1
}

//遍历
func TraverseFunc(nlen int, skip func(i int) int) {
	for i := 0; i < nlen; i++ {
		l_skipn := skip(i)
		if i += l_skipn; i < -1 {
			i = -1
		}
	}
}
func Atoi16(s string) (int16, error) {
	i64, err := strconv.ParseInt(s, 0, 16)
	if nerr, ok := err.(*strconv.NumError); ok {
		nerr.Func = "Atoi16"
	}
	return int16(i64), err
}

//返回>={搜索值}的最小索引,fmore_eq:func(i int){return {搜索值}<={val}[i]}
func Lowerbound(nlen int, fmore_eq func(i int) bool) int {
	return sort.Search(nlen, fmore_eq)
}

//返回>{搜索值}的最小索引,fmore:func(i int){return {搜索值}<{val}[i]}
func Upperbound(neln int, fmore func(j int) bool) int {
	return sort.Search(neln, fmore)
}
func SplitFields(s string, sepchas string, or_space bool) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return or_space && unicode.IsSpace(r) || len(sepchas) > 0 && strings.IndexRune(sepchas, r) >= 0
	})
}
func SplitInt16s(s string, sepchas string, or_space bool) (retval []int16, reterr error) {
	l_fields := SplitFields(s, sepchas, or_space)
	retval = make([]int16, len(l_fields))
	for i := 0; i < len(l_fields); i++ {
		if retval[i], reterr = Atoi16(l_fields[i]); reterr != nil {
			return nil, reterr
		}
	}
	return
}

func Unique2(nlen int, fless func(i, j int) bool, fswap func(i, j int), fequal func(i, j int) bool, fassign func(i, j int)) (deletedcount int) {
	if nlen < 2 {
		return 0
	}
	Sort(nlen, fless, fswap)

	deletedcount = 0
	for i := 1; i < nlen-deletedcount; i++ {
		if fequal(i-1, i+deletedcount) {
			deletedcount++
			i--
			continue
		}
		if deletedcount > 0 {
			fassign(i, i+deletedcount)
		}
	}
	return
}

func UniqueInt(arr ...int) []int {
	l_deletedcount := Unique2(len(arr),
		func(i, j int) bool { return arr[i] < arr[j] },
		func(i, j int) { arr[i], arr[j] = arr[j], arr[i] },
		func(i, j int) bool { return arr[i] == arr[j] },
		func(i, j int) { arr[i] = arr[j] })
	return arr[:len(arr)-l_deletedcount]
}

func DelMulti(nlen int, fassign func(i, j int), deletedidx ...int) (deletedcount int) {
	if len(deletedidx) < 1 {
		return 0
	} else if len(deletedidx) > 2 {
		deletedidx = UniqueInt(deletedidx...)
	}
	for i := 0; i < len(deletedidx); i++ {
		if deletedidx[i] >= 0 {
			deletedidx = deletedidx[i:]
			break
		}
	}
	if len(deletedidx) < 1 {
		return 0
	}
	if deletedidx[len(deletedidx)-1] < nlen {
		deletedidx = append(deletedidx, nlen) //末尾一定要追加一个大数
	}
	deletedcount = 0
	for i := deletedidx[0]; i < nlen-deletedcount; i++ {
		if i == deletedidx[deletedcount]-deletedcount {
			deletedcount++
			i--
			continue
		}
		fassign(i, i+deletedcount)
	}
	return
}

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
