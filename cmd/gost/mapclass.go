package main

import (
	"bytes"
)

type MAP struct {
	m_canRepeat bool
	m_compfunc  func(key1, key2 interface{}) int8
	m_root      *tnode
	m_len       int
}

func (this *MAP) MAP(canRepeat bool, compfunc func(key1, key2 interface{}) int8) *MAP {
	this.m_canRepeat = canRepeat
	this.m_compfunc = compfunc
	this.m_root = nil
	this.m_len = 0

	return this
}

func (this *MAP) Clear(deep int) {
	if this.m_root != nil {
		this.m_root.clear(deep)
		this.m_root = nil
		this.m_len = 0
	}
}
func (this *MAP) Move(_Right MAP) {
	this.Clear(1000)
	this.m_canRepeat = _Right.m_canRepeat
	this.m_compfunc = _Right.m_compfunc
	this.m_root, _Right.m_root = _Right.m_root, nil
	this.m_len, _Right.m_len = _Right.m_len, 0
}

type ifduplicate interface {
	Duplicate(deep int) interface{}
}

func (this *MAP) Duplicate(deep int) *MAP {
	l_map := new(MAP).MAP(this.m_canRepeat, this.m_compfunc)
	for it := this.Begin(); it != this.End(); it = it.Next() {
		if deep < 2 {
			l_map.Insert(it.Key(), it.Value())
		} else if map_, ok := it.Value().(*MAP); ok {
			l_map.Insert(it.Key(), map_.Duplicate(deep-1))
			// } else if map_, ok := it.Value().(*dbo.MAPVT); ok {
			// 	l_map.Insert(it.Key(), map_.Duplicate(deep-1))
		} else if dup, ok := it.Value().(ifduplicate); ok {
			l_map.Insert(it.Key(), dup.Duplicate(deep-1))
		} else {
			l_map.Insert(it.Key(), it.Value())
		}
	}
	return l_map
}
func (this *MAP) Append(begin, end MAPIter, deep int, iterfunc func(newit MAPIter)) {
	for ; begin != end; begin = begin.Next() {
		sucs, newit := false, MAPIter{nil}
		if deep < 2 {
			sucs, newit = this.Insert(begin.Key(), begin.Value())
		} else if map_, ok := begin.Value().(*MAP); ok {
			sucs, newit = this.Insert(begin.Key(), map_.Duplicate(deep-1))
			// } else if map_, ok := begin.Value().(*dbo.MAPVT); ok {
			// 	sucs, newit = this.Insert(begin.Key(), map_.Duplicate(deep-1))
		} else if dup, ok := begin.Value().(ifduplicate); ok {
			sucs, newit = this.Insert(begin.Key(), dup.Duplicate(deep-1))
		} else {
			sucs, newit = this.Insert(begin.Key(), begin.Value())
		}
		if sucs && iterfunc != nil {
			iterfunc(newit)
		}
	}
}
func (this *MAP) GetCompFunc() func(key1, key2 interface{}) int8 { return this.m_compfunc }

func (this *MAP) Empty() bool { return this.m_len < 1 }
func (this *MAP) Len() int    { return this.m_len }
func (this *MAP) Size() int   { return this.m_len }
func (this *MAP) Index(key interface{}) (val interface{}, found bool) {
	if it := this.Find(key); it != this.End() {
		return it.Value(), true
	}
	return nil, false
}
func (this *MAP) IsAmong(key interface{}) (found bool) {
	if it := this.Find(key); it != this.End() {
		return true
	}
	return false
}

func (this *MAP) Insert(key, value interface{}) (sucs bool, it MAPIter) {
	var l_it MAPIter
	if sucs, l_it.tnode = this.insertNode(&this.m_root, nil, key); sucs {
		l_it.m_value = value
		this.m_len++
	}
	it = l_it
	return
}
func (this *MAP) Erase(it MAPIter) {
	if it == this.End() || it.tnode == nil {
		return
	}
	l_it := it
	if this.deleteNode(l_it.tnode, false) {
		this.m_len--
	}
}
func (this *MAP) Remove(key interface{}) int {
	if !this.m_canRepeat {
		if it := this.Find(key); it == this.End() {
			return 0
		} else {
			l_it := it
			this.deleteNode(l_it.tnode, true)
			this.m_len--
			return 1
		}
	} else if it := this.LowerBound(key); it == this.End() {
		return 0
	} else {
		l_count := 0
		for itx := it.Next(); it != this.End(); it, itx = itx, itx.Next() {
			if this.m_compfunc(key, it.Key()) != 0 {
				break
			}
			l_it := it
			this.deleteNode(l_it.tnode, true)
			this.m_len--
			l_count++
		}
		return l_count
	}
}
func (this *MAP) Find(key interface{}) MAPIter {
	return MAPIter{tnode: this.find(this.m_root, key)}
}

func (this *MAP) LowerBound(key interface{}) MAPIter {
	return MAPIter{tnode: this.lowerbound(this.m_root, key)}
}
func (this *MAP) UpperBound(key interface{}) MAPIter {
	return MAPIter{tnode: this.uppperbound(this.m_root, key)}
}
func (this *MAP) EqualBound(key interface{}) (MAPIter, MAPIter) {
	return this.LowerBound(key), this.UpperBound(key)
}

//return x>=key1 && x<key2
func (this *MAP) Between1(key1, key2 interface{}) (MAPIter, MAPIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.LowerBound(key1), this.LowerBound(key2)
	} else {
		return this.LowerBound(key2), this.LowerBound(key1)
	}
}

//return x>=key1 && x<=key2
func (this *MAP) Between2(key1, key2 interface{}) (MAPIter, MAPIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.LowerBound(key1), this.UpperBound(key2)
	} else {
		return this.LowerBound(key2), this.UpperBound(key1)
	}
}

//return x>key1 && x<key2
func (this *MAP) Between3(key1, key2 interface{}) (MAPIter, MAPIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.UpperBound(key1), this.LowerBound(key2)
	} else {
		return this.UpperBound(key2), this.LowerBound(key1)
	}
}

//return x>key1 && x<=key2
func (this *MAP) Between4(key1, key2 interface{}) (MAPIter, MAPIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.UpperBound(key1), this.UpperBound(key2)
	} else {
		return this.UpperBound(key2), this.UpperBound(key1)
	}
}

func (this *MAP) String() string {
	l_root := create_showtree(this.m_root)
	l_lines := make([]*bytes.Buffer, 0, 8)
	create_show_strings(l_root, 0, 0, &l_lines)
	var l_strbuf bytes.Buffer
	for _, line := range l_lines {
		l_strbuf.WriteString(line.String())
		l_strbuf.WriteString("\n")
	}
	return l_strbuf.String()
}

func NewMAP(canRepeat bool, compfunc func(key1, key2 interface{}) int8) *MAP {
	return new(MAP).MAP(canRepeat, compfunc)
}

// func GetMAPIterNil() interface{} {
// 	return MAPIter{tnode: nil}
// }
