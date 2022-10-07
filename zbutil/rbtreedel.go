package zbutil

func (this *RBtree) deleteNode(node *tnode, confirm_isinmap bool) bool {
	if !confirm_isinmap && node.mostroot() != this.m_root {
		return false
	}
	if node.m_left != nil && node.m_right != nil {
		l_swaped := node.m_right.mostleft()

		node.m_parent, l_swaped.m_parent = l_swaped.m_parent, node.m_parent
		node.m_left, l_swaped.m_left = l_swaped.m_left, node.m_left
		node.m_right, l_swaped.m_right = l_swaped.m_right, node.m_right
		node.m_color, l_swaped.m_color = l_swaped.m_color, node.m_color

		if l_swaped.m_right == l_swaped {
			l_swaped.m_right = node
		}
		l_swaped.m_left.m_parent, l_swaped.m_right.m_parent = l_swaped, l_swaped

		if node.m_right != nil {
			node.m_right.m_parent = node
		}
		if node.m_parent.m_left == l_swaped {
			node.m_parent.m_left = node
		} else {
			node.m_parent.m_right = node
		}

		if l_swaped.m_parent == nil {
			this.m_root = l_swaped
		} else if l_swaped.m_parent.m_left == node {
			l_swaped.m_parent.m_left = l_swaped
		} else {
			l_swaped.m_parent.m_right = l_swaped
		}
	}
	this.deleteOne(node)
	return true
}
func (this *RBtree) deleteOne(node *tnode) {
	var l_child *tnode
	if node.m_left == nil {
		l_child = node.m_right
	} else {
		l_child = node.m_left
	}

	if node.m_parent != nil {
	} else if l_child == nil {
		this.m_root = nil
		return
	} else {
		this.m_root = l_child
		l_child.m_parent = nil
		l_child.m_color = black
		return
	}

	if node.m_color == red {
		if node == node.m_parent.m_left {
			node.m_parent.m_left = l_child
		} else {
			node.m_parent.m_right = l_child
		}
		if l_child != nil {
			l_child.m_parent = node.m_parent
		}
		return
	} else if l_child == nil { // 如果没有孩子节点，则添加一个临时孩子节点

		this.deleteArrange(node)

		if node.m_parent.m_left == node {
			node.m_parent.m_left = nil
		} else {
			node.m_parent.m_right = nil
		}
	} else if l_child.m_color == red {
		if node.m_parent.m_left == node {
			node.m_parent.m_left = l_child
		} else {
			node.m_parent.m_right = l_child
		}
		l_child.m_parent = node.m_parent

		l_child.m_color = black
		return
	} else {
		if node.m_parent.m_left == node {
			node.m_parent.m_left = l_child
		} else {
			node.m_parent.m_right = l_child
		}
		l_child.m_parent = node.m_parent

		this.deleteArrange(l_child)
	}
}
func getSibling(node *tnode) *tnode {
	if node == node.m_parent.m_left {
		return node.m_parent.m_right
	} else {
		return node.m_parent.m_left
	}
}
func (this *RBtree) deleteArrange(node *tnode) {
	if node.m_parent == nil {
		node.m_color = black
		return
	}
	l_sibling := getSibling(node)
	if l_sibling.m_color == red {
		if node == node.m_parent.m_left {
			this.逆时旋(node.m_parent)
		} else {
			this.顺时旋(node.m_parent)
		}
		node.m_parent.m_parent.m_color = black
		node.m_parent.m_color = red

		l_sibling = getSibling(node)
	}
	//注意：这里node的兄弟节点发生了变化，不再是原来的兄弟节点
	is_sib_left_red := black
	is_sib_right_red := black
	if l_sibling.m_left != nil {
		is_sib_left_red = l_sibling.m_left.m_color
	}
	if l_sibling.m_right != nil {
		is_sib_right_red = l_sibling.m_right.m_color
	}
	//开始处理颜色
	if is_sib_left_red && is_sib_right_red {
	} else if is_sib_left_red {
		if node == node.m_parent.m_left {
			this.顺时旋(l_sibling)
			l_sibling.m_color = red
			l_sibling = l_sibling.m_parent
			l_sibling.m_color = black
		}
	} else if is_sib_right_red {
		if node == node.m_parent.m_right {
			this.逆时旋(l_sibling)
			l_sibling.m_color = red
			l_sibling = l_sibling.m_parent
			l_sibling.m_color = black
		}
	} else {
		l_sibling.m_color = red
		if node.m_parent.m_color == red {
			node.m_parent.m_color = black
		} else {
			this.deleteArrange(node.m_parent)
		}
		return
	}
	if node == node.m_parent.m_left {
		this.逆时旋(node.m_parent)
	} else {
		this.顺时旋(node.m_parent)
	}
	node.m_parent.m_parent.m_color = node.m_parent.m_color
	node.m_parent.m_color = black
	getSibling(node.m_parent).m_color = black
}
