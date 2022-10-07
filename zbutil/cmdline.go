package zbutil

import (
	"os"
	"sync/atomic"
)

var g_cmdln_map RBtreeIstr
var g_cmdln_once int32 = 0
var g_cmdln_load bool = false

func CmdParmLike(leftpart string) string {
	if !g_cmdln_load {
		g_cmdln_load = true
		if atomic.AddInt32(&g_cmdln_once, 1) == 1 {
			g_cmdln_map.RBtreeIstr(true)
			for i := 0; i < len(os.Args); i++ {
				g_cmdln_map.Insert(os.Args[i], nil)
			}
		}
	}
	if it := g_cmdln_map.LowerBound(leftpart); it != g_cmdln_map.End() {
		l_k := it.Key().(string)
		if len(l_k) < len(leftpart) {
			return ""
		}
		if g_cmdln_map.GetCompFunc()(leftpart, l_k[:len(leftpart)]) == 0 {
			return l_k
		} else {
			return ""
		}
	}
	return ""
}
