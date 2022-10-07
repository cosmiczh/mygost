package zbutil

import (
	"errors"
	"fmt"
)

type color bool

const (
	red        color = true
	black      color = false
	anti_clock bool  = true  // 逆时针旋转
	clockwise  bool  = false // 顺时针旋转
)

type tnode struct {
	m_parent, m_left, m_right *tnode
	m_color                   color
	m_key, m_value            interface{}
}

func (this *tnode) mostroot() (root *tnode) {
	for root = this; root.m_parent != nil; root = root.m_parent {
	}
	return
}
func (this *tnode) mostleft() (left *tnode) {
	for left = this; left.m_left != nil; left = left.m_left {
	}
	return
}
func (this *tnode) mostright() (right *tnode) {
	for right = this; right.m_right != nil; right = right.m_right {
	}
	return
}
func (this *tnode) next() *tnode {
	if this.m_right != nil {
		return this.m_right.mostleft()
	} else {
		for ; this.m_parent != nil && this == this.m_parent.m_right; this = this.m_parent {
		}
		return this.m_parent
	}
}
func (this *tnode) prev() *tnode {
	if this.m_left != nil {
		return this.m_left.mostright()
	} else {
		for ; this.m_parent != nil && this == this.m_parent.m_left; this = this.m_parent {
		}
		return this.m_parent
	}
}

type ifclear interface {
	Clear(deep int)
}

func (this *tnode) clear(deep int) {
	if this.m_left != nil {
		this.m_left.clear(deep)
		this.m_left.m_parent = nil
		this.m_left = nil
	}
	if this.m_right != nil {
		this.m_right.clear(deep)
		this.m_right.m_parent = nil
		this.m_right = nil
	}
	if deep > 1 {
		if if_clear, ok := this.m_value.(ifclear); ok {
			if_clear.Clear(deep - 1)
		}
	}
}

func (this *tnode) 旋转(旋转方向 bool) (*tnode, error) {
	if this == nil {
		return nil, nil
	} else if 旋转方向 && this.m_right == nil {
		return nil, errors.New("逆时旋右节点不能为空")
	} else if !旋转方向 && this.m_left == nil {
		return nil, errors.New("顺时旋左节点不能为空")
	}

	l_parent := this.m_parent
	if 旋转方向 { //逆时
		grandson := this.m_right.m_left

		this.m_right.m_left = this
		this.m_parent = this.m_right

		this.m_right = grandson
		if grandson != nil {
			grandson.m_parent = this
		}
	} else { //顺时
		grandson := this.m_left.m_right

		this.m_left.m_right = this
		this.m_parent = this.m_left

		this.m_left = grandson
		if grandson != nil {
			grandson.m_parent = this
		}
	}
	this.m_parent.m_parent = l_parent
	// 判断是否换了根节点
	if l_parent == nil {
		return this.m_parent, nil
	} else if l_parent.m_left == this {
		l_parent.m_left = this.m_parent
	} else {
		l_parent.m_right = this.m_parent

	}
	return nil, nil
}
func (this *RBtree) 逆时旋(node *tnode) {
	if tmproot, err := node.旋转(anti_clock); err == nil {
		if tmproot != nil {
			this.m_root = tmproot
		}
	} else {
		fmt.Printf(err.Error())
	}
}
func (this *RBtree) 顺时旋(node *tnode) {
	if tmproot, err := node.旋转(clockwise); err == nil {
		if tmproot != nil {
			this.m_root = tmproot
		}
	} else {
		fmt.Printf(err.Error())
	}
}
