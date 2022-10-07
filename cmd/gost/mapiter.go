package main

import (
	"bytes"
	"fmt"
	"strings"
)

type MAPIter struct{ *tnode }

func (this MAPIter) Key() interface{} {
	return this.m_key
}
func (this MAPIter) Value() interface{} {
	return this.m_value
}
func (this MAPIter) SetValue(val interface{}) { //相当于c++的operator*
	this.m_value = val
}
func (this MAPIter) Next() MAPIter {
	if this.tnode == nil {
		return this
	} else {
		return MAPIter{tnode: this.next()}
	}
}
func (this MAPIter) Prev() MAPIter {
	if this.tnode == nil {
		return this
	} else {
		return MAPIter{tnode: this.prev()}
	}
}

func (this *MAP) Begin() MAPIter {
	if this.m_root != nil {
		return MAPIter{tnode: this.m_root.mostleft()}
	}
	return this.End()
}
func (this *MAP) End() MAPIter {
	return MAPIter{tnode: nil}
}

//反向遍历获取Begin迭代器
func (this *MAP) RBegin() MAPIter {
	if this.m_root != nil {
		return MAPIter{tnode: this.m_root.mostright()}
	}
	return this.REnd()
}
func (this *MAP) REnd() MAPIter {
	return MAPIter{tnode: nil}
}
func (this *MAP) MidIterate(iterfunc func(it MAPIter) (isbreak bool)) (isbreak bool) {
	if iterfunc == nil {
		return false
	}
	var lf_recurs func(tree *tnode) (isbreak bool)
	lf_recurs = func(tree *tnode) (isbreak bool) {
		if tree == nil {
			return false
		}
		if isbreak = iterfunc(MAPIter{tnode: tree}); isbreak {
			return
		}
		if isbreak = lf_recurs(tree.m_left); isbreak {
			return
		}
		return lf_recurs(tree.m_right)
	}
	return lf_recurs(this.m_root)
}
func (this *MAP) LeftIterate(iterfunc func(it MAPIter) (isbreak bool)) (isbreak bool) {
	if iterfunc == nil {
		return false
	}
	var lf_recurs func(tree *tnode) (isbreak bool)
	lf_recurs = func(tree *tnode) (isbreak bool) {
		if tree == nil {
			return false
		}
		if isbreak = lf_recurs(tree.m_left); isbreak {
			return
		}
		if isbreak = iterfunc(MAPIter{tnode: tree}); isbreak {
			return
		}
		return lf_recurs(tree.m_right)
	}
	return lf_recurs(this.m_root)
}
func (this *MAP) RightIterate(iterfunc func(it MAPIter) (isbreak bool)) (isbreak bool) {
	if iterfunc == nil {
		return false
	}
	var lf_recurs func(tree *tnode) (isbreak bool)
	lf_recurs = func(tree *tnode) (isbreak bool) {
		if tree == nil {
			return false
		}
		if isbreak = lf_recurs(tree.m_right); isbreak {
			return
		}
		if isbreak = iterfunc(MAPIter{tnode: tree}); isbreak {
			return
		}
		return lf_recurs(tree.m_left)
	}
	return lf_recurs(this.m_root)
}

type strnode struct {
	m_left, m_right                   *strnode
	mlen_left, mlen_right, mlen_space int
	m_my                              string
}

func create_showtree(node *tnode) *strnode {
	if node == nil {
		return nil
	}
	l_strnode := &strnode{m_left: create_showtree(node.m_left), m_right: create_showtree(node.m_right)}
	if l_strnode.m_left == nil {
		l_strnode.mlen_left = 0
	} else {
		l_strnode.mlen_left = l_strnode.m_left.mlen_left + l_strnode.m_left.mlen_space + l_strnode.m_left.mlen_right
		if l_strnode.mlen_left < len(l_strnode.m_left.m_my) {
			l_strnode.mlen_left = len(l_strnode.m_left.m_my)
		}
	}
	if l_strnode.m_right == nil {
		l_strnode.mlen_right = 0
	} else {
		l_strnode.mlen_right = l_strnode.m_right.mlen_left + l_strnode.m_right.mlen_space + l_strnode.m_right.mlen_right
		if l_strnode.mlen_right < len(l_strnode.m_right.m_my) {
			l_strnode.mlen_right = len(l_strnode.m_right.m_my)
		}
	}
	if node.m_color == red {
		l_strnode.m_my = "R[%v"
	} else {
		l_strnode.m_my = "B[%v"
	}
	if node.m_value == nil {
		l_strnode.m_my += "]"
		l_strnode.m_my = fmt.Sprintf(l_strnode.m_my, node.m_key)
	} else {
		l_strnode.m_my += ":%v]"
		l_strnode.m_my = fmt.Sprintf(l_strnode.m_my, node.m_key, node.m_value)
	}
	l_strnode.mlen_space = 1
	l_myhalf := len(l_strnode.m_my) / 2
	if l_myhalf > l_strnode.mlen_left {
		l_strnode.mlen_left = l_myhalf
		l_strnode.mlen_space = 0
	}
	l_myhalf = (len(l_strnode.m_my) + 1) / 2
	if l_myhalf > l_strnode.mlen_right+1 {
		l_strnode.mlen_right = l_myhalf
		l_strnode.mlen_space = 0
	}

	return l_strnode
}
func create_show_strings(node *strnode, baseoffs int, idx int, show *[]*bytes.Buffer) {
	if node == nil {
		return
	}
	l_show := *show
	for idx >= len(l_show) {
		l_show = append(l_show, &bytes.Buffer{})
		*show = l_show
	}
	l_baseoffs := baseoffs + node.mlen_left
	l_myhalf := len(node.m_my) / 2
	l_show[idx].WriteString(strings.Repeat(" ", l_baseoffs-l_myhalf-l_show[idx].Len()))
	l_show[idx].WriteString(node.m_my)
	create_show_strings(node.m_left, baseoffs, idx+1, show)
	create_show_strings(node.m_right, l_baseoffs+node.mlen_space, idx+1, show)
}
