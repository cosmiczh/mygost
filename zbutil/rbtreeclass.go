package zbutil

import (
	"bytes"
)

type RBtree struct {
	m_canRepeat bool
	m_compfunc  func(key1, key2 interface{}) int8
	m_root      *tnode
	m_len       int
}

func (this *RBtree) RBtree(canRepeat bool, compfunc func(key1, key2 interface{}) int8) *RBtree {
	this.m_canRepeat = canRepeat
	this.m_compfunc = compfunc
	this.m_root = nil
	this.m_len = 0

	return this
}

func (this *RBtree) Clear(deep int) {
	if this.m_root != nil {
		this.m_root.clear(deep)
		this.m_root = nil
		this.m_len = 0
	}
}
func (this *RBtree) Move(_Right RBtree) {
	this.Clear(1000)
	this.m_canRepeat = _Right.m_canRepeat
	this.m_compfunc = _Right.m_compfunc
	this.m_root, _Right.m_root = _Right.m_root, nil
	this.m_len, _Right.m_len = _Right.m_len, 0
}

type ifduplicate interface {
	Duplicate(deep int) interface{}
}

func (this *RBtree) Duplicate(deep int) *RBtree {
	l_map := new(RBtree).RBtree(this.m_canRepeat, this.m_compfunc)
	for it := this.Begin(); it != this.End(); it = it.Next() {
		if deep < 2 {
			l_map.Insert(it.Key(), it.Value())
		} else if map_, ok := it.Value().(*RBtree); ok {
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
func (this *RBtree) Append(begin, end RBtreeIter, deep int, iterfunc func(newit RBtreeIter)) {
	for ; begin != end; begin = begin.Next() {
		sucs, newit := false, RBtreeIter{nil}
		if deep < 2 {
			sucs, newit = this.Insert(begin.Key(), begin.Value())
		} else if map_, ok := begin.Value().(*RBtree); ok {
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
func (this *RBtree) GetCompFunc() func(key1, key2 interface{}) int8 { return this.m_compfunc }

func (this *RBtree) Empty() bool { return this.m_len < 1 }
func (this *RBtree) Len() int    { return this.m_len }
func (this *RBtree) Size() int   { return this.m_len }
func (this *RBtree) Index(key interface{}) (val interface{}, found bool) {
	if it := this.Find(key); it != this.End() {
		return it.Value(), true
	}
	return nil, false
}
func (this *RBtree) IsAmong(key interface{}) (found bool) {
	if it := this.Find(key); it != this.End() {
		return true
	}
	return false
}

func (this *RBtree) Insert(key, value interface{}) (sucs bool, it RBtreeIter) {
	var l_it RBtreeIter
	if sucs, l_it.tnode = this.insertNode(&this.m_root, nil, key); sucs {
		l_it.m_value = value
		this.m_len++
	}
	it = l_it
	return
}
func (this *RBtree) Erase(it RBtreeIter) {
	if it == this.End() || it.tnode == nil {
		return
	}
	l_it := it
	if this.deleteNode(l_it.tnode, false) {
		this.m_len--
	}
}
func (this *RBtree) Remove(key interface{}) int {
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
func (this *RBtree) Find(key interface{}) RBtreeIter {
	return RBtreeIter{tnode: this.find(this.m_root, key)}
}

func (this *RBtree) LowerBound(key interface{}) RBtreeIter {
	return RBtreeIter{tnode: this.lowerbound(this.m_root, key)}
}
func (this *RBtree) UpperBound(key interface{}) RBtreeIter {
	return RBtreeIter{tnode: this.uppperbound(this.m_root, key)}
}
func (this *RBtree) EqualBound(key interface{}) (RBtreeIter, RBtreeIter) {
	return this.LowerBound(key), this.UpperBound(key)
}

//return x>=key1 && x<key2
func (this *RBtree) Between1(key1, key2 interface{}) (RBtreeIter, RBtreeIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.LowerBound(key1), this.LowerBound(key2)
	} else {
		return this.LowerBound(key2), this.LowerBound(key1)
	}
}

//return x>=key1 && x<=key2
func (this *RBtree) Between2(key1, key2 interface{}) (RBtreeIter, RBtreeIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.LowerBound(key1), this.UpperBound(key2)
	} else {
		return this.LowerBound(key2), this.UpperBound(key1)
	}
}

//return x>key1 && x<key2
func (this *RBtree) Between3(key1, key2 interface{}) (RBtreeIter, RBtreeIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.UpperBound(key1), this.LowerBound(key2)
	} else {
		return this.UpperBound(key2), this.LowerBound(key1)
	}
}

//return x>key1 && x<=key2
func (this *RBtree) Between4(key1, key2 interface{}) (RBtreeIter, RBtreeIter) {
	if this.m_compfunc(key1, key2) <= 0 {
		return this.UpperBound(key1), this.UpperBound(key2)
	} else {
		return this.UpperBound(key2), this.UpperBound(key1)
	}
}

func (this *RBtree) String() string {
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

func NewRBtree(canRepeat bool, compfunc func(key1, key2 interface{}) int8) *RBtree {
	return new(RBtree).RBtree(canRepeat, compfunc)
}

// func GetMAPIterNil() interface{} {
// 	return MAPIter{tnode: nil}
// }
