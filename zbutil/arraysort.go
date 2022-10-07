package zbutil

import "sort"

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
		if i += skip(i); i < -1 {
			i = -1
		}
	}
}
func TraverseFind(nlen int, eque func(i int) bool) int {
	for i := 0; i < nlen; i++ {
		if eque(i) {
			return i
		}
	}
	return -1
}

//Lowerbound 遍历二分法搜索
//返回>={搜索值}的最小索引,fmore_eq:func(i int){return {搜索值}<={val}[i]}
func Lowerbound(nlen int, fmore_eq func(i int) bool) int {
	return sort.Search(nlen, fmore_eq)
}

//Upperbound 遍历二分法搜索
//返回>{搜索值}的最小索引,fmore:func(i int){return {搜索值}<{val}[i]}
func Upperbound(neln int, fmore func(j int) bool) int {
	return sort.Search(neln, fmore)
}
