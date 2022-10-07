package main

func (this *MAP) lowerbound(node *tnode, key interface{}) *tnode {
	if node == nil {
		return nil
	} else if bigsmall := this.m_compfunc(node.m_key, key); bigsmall >= 0 {
		if node := this.lowerbound(node.m_left, key); node != nil {
			return node
		}
		return node
	} else {
		return this.lowerbound(node.m_right, key)
	}
}
func (this *MAP) uppperbound(node *tnode, key interface{}) *tnode {
	if node == nil {
		return nil
	} else if bigsmall := this.m_compfunc(node.m_key, key); bigsmall > 0 {
		if node := this.uppperbound(node.m_left, key); node != nil {
			return node
		}
		return node
	} else {
		return this.uppperbound(node.m_right, key)
	}
}
func (this *MAP) find(node *tnode, key interface{}) *tnode {
	if node == nil {
		return nil
	} else if bigsmall := this.m_compfunc(key, node.m_key); bigsmall == 0 {
		return node
	} else if bigsmall < 0 {
		return this.find(node.m_left, key)
	} else {
		return this.find(node.m_right, key)
	}
}
