package zbutil

func (this *RBtree) insertNode(pnode **tnode, parent *tnode, key interface{}) (sucs bool, newnode *tnode) {
	if node := (*pnode); node == nil {
		sucs, newnode = true, &tnode{m_parent: parent, m_key: key, m_color: red}
		*pnode = newnode
		this.insertArrange(newnode)
		return
	} else if bigsmall := this.m_compfunc(key, node.m_key); bigsmall == 0 {
		if !this.m_canRepeat {
			return false, node
		} else { //重复的插入右边
			return this.insertNode(&node.m_right, node, key)
		}
	} else if bigsmall > 0 { // 插入数据大于父节点，插入右边
		return this.insertNode(&node.m_right, node, key)
	} else { // 插入数据小于父节点，插入左节点
		return this.insertNode(&node.m_left, node, key)
	}
}
func (this *RBtree) insertArrange(node *tnode) {
	l_parent := node.m_parent
	if l_parent == nil { //根节点
		node.m_color = black
		return
	} else if l_parent.m_color == black {
		return
	}
	l_grandparent := l_parent.m_parent
	l_parentisleft := (l_parent == l_grandparent.m_left)
	var l_uncle *tnode
	if l_parentisleft {
		l_uncle = l_grandparent.m_right
	} else {
		l_uncle = l_grandparent.m_left
	}
	if l_uncle != nil && l_uncle.m_color == red {
		l_parent.m_color = black
		l_uncle.m_color = black
		l_grandparent.m_color = red
		this.insertArrange(l_grandparent)
		return
	}
	l_isleft := (node == l_parent.m_left)
	if l_isleft && l_parentisleft {
		this.顺时旋(l_grandparent)
		l_parent.m_color = black
		l_grandparent.m_color = red
		return
	}
	if !l_isleft && !l_parentisleft {
		this.逆时旋(l_grandparent)
		l_parent.m_color = black
		l_grandparent.m_color = red
		return
	}
	if !l_isleft && l_parentisleft {
		this.逆时旋(l_parent)
		this.顺时旋(l_grandparent)
		node.m_color = black
		l_grandparent.m_color = red
		return
	}
	if l_isleft && !l_parentisleft {
		this.顺时旋(l_parent)
		this.逆时旋(l_grandparent)
		node.m_color = black
		l_grandparent.m_color = red
		return
	}
}
