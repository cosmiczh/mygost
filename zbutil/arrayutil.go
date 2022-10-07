package zbutil

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

func UniqueInt(arr ...int) []int {
	l_deletedcount := Unique2(len(arr),
		func(i, j int) bool { return arr[i] < arr[j] },
		func(i, j int) { arr[i], arr[j] = arr[j], arr[i] },
		func(i, j int) bool { return arr[i] == arr[j] },
		func(i, j int) { arr[i] = arr[j] })
	return arr[:len(arr)-l_deletedcount]
}
