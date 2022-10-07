package gost

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	ipdb "github.com/ipipdotnet/ipdb-go"
)

const IPv4len = 4

type IPv4sec struct {
	start uint32
	end   uint32
}

var (
	file1 = GetExeDir() + "/chn_ip.txt"
	file2 = GetExeDir() + "/ipipfree.ipdb"
)
var (
	g_lastread               = time.Now()
	g_ipv4secs, g_ipdb, _, _ = ReadIPset(file1, file2, &g_mod1t, &g_mod2t, &g_f1size, &g_f2size)
)
var (
	g_mtx              RWMutex2
	g_mod1t, g_mod2t   time.Time
	g_f1size, g_f2size int64
	g_once_run         int32
)

func chnips_contains(ip string, in_chn bool) bool {
	l_now := time.Now()
	if l_now.Sub(g_lastread) >= 10*time.Second && g_once_run == 0 {
		func() {
			if unlock := g_mtx.TryLock(); unlock != nil {
				defer unlock()
				if atomic.LoadInt32(&g_once_run) == 0 && l_now.Sub(g_lastread) >= 10*time.Second {
					g_lastread = l_now
					atomic.StoreInt32(&g_once_run, 1)
					go func() {
						defer atomic.StoreInt32(&g_once_run, 0)
						l_ipv4s, l_ipdb, err1, err2 := ReadIPset(file1, file2, &g_mod1t, &g_mod2t, &g_f1size, &g_f2size)
						l_deleted1, l_deleted2 := false, false
						if err2 != nil {
							if _, ok := err2.(*os.PathError); ok {
								if g_ipdb != nil {
									l_deleted2 = true
									// loglv.War.Printf("IsAcceptedIP(ip:%s)1 ipipfree.ipdb文件已经被删除，关闭ipdb过滤功能", ip)
								}
							} else {
								fmt.Printf("IsAcceptedIP(ip:%s)2 错误:%v\n", ip, err1)
							}
						}
						if err1 != nil {
							if _, ok := err1.(*os.PathError); ok {
								if len(g_ipv4secs) > 0 {
									l_deleted1 = true
									// loglv.War.Printf("IsAcceptedIP(ip:%s)3 chn_ip.txt文件已经被删除，关闭txtIP过滤功能", ip)
								}
							} else {
								fmt.Printf("IsAcceptedIP(ip:%s)4 错误:%v\n", ip, err1)
							}
						}
						if len(l_ipv4s) > 0 {
							var l_ip_count int64
							TraverseFunc(len(l_ipv4s), func(i int) int {
								if l_ipv4s[i].start <= l_ipv4s[i].end && l_ipv4s[i].start>>24 != 10 && l_ipv4s[i].start>>24 != 127 &&
									l_ipv4s[i].start>>16 != 172<<8+16 && l_ipv4s[i].start>>16 != 192<<8+168 {

									l_ip_count += int64(l_ipv4s[i].end-l_ipv4s[i].start) + 1
								}
								return 0
							})
							// loglv.Dbg.Printf("IsAcceptedIP(ip:%s)5 txtIP读出了%d个区段的ip，共囊括了%d个IP地址", ip, len(l_ipv4s), l_ip_count)
						}

						if l_deleted1 || len(l_ipv4s) > 0 || l_deleted2 || l_ipdb != nil {
							defer g_mtx.Lock()()
							if l_deleted1 || len(l_ipv4s) > 0 {
								g_ipv4secs = l_ipv4s
							}
							if l_deleted2 || l_ipdb != nil {
								g_ipdb = l_ipdb
							}
						}
					}()
				}
			}
		}()
	}
	if len(g_ipv4secs) < 1 { //如果没有配置，就所有都放行
		// loglv.Dbg.Printf("IsAcceptedIP(ip:%s) 0个区段，无ip过滤功能", ip)
		if g_ipdb == nil {
			fmt.Printf("[1]IP:%s,ret:true\n", ip)
			return true //两个都不生效和后面有一个生效的状况是不一样的
		}
		if ret, err := g_ipdb.Find(ip, "CN"); err != nil {
			fmt.Printf("[2]IP:%s,ret:true\n", ip)
			return true
		} else if len(ret) < 1 {
			fmt.Printf("[3]IP:%s,ret:%v\n", ip, !in_chn)
			return !in_chn
		} else if ret[0] != "中国" {
			fmt.Printf("[4]IP:%s,ret:%v\n", ip, !in_chn)
			return !in_chn
		} else if len(ret) < 2 {
			fmt.Printf("[5]IP:%s,ret:%v\n", ip, in_chn)
			return in_chn
		} else {
			fmt.Printf("[6]IP:%s,ret[1]:%s,in_chn:%v\n", ip, ret[1], in_chn)
			l_in_twhkmk := (ret[1] == "台湾") || (ret[1] == "香港") || (ret[1] == "澳门")
			return (in_chn && !l_in_twhkmk) || (!in_chn && l_in_twhkmk)
		}
	}
	l_ip, err := SplitInt16s(ip, ".", true)
	if err != nil || len(l_ip) != IPv4len {
		fmt.Printf("[9]IPv4地址%s格式错误:%v\n", ip, err)
		return false
	}
	var l_nip uint32
	for i := 0; i < IPv4len; i++ {
		l_nip += uint32(byte(l_ip[i])) << uint((3-i)*8)
	}

	defer g_mtx.RLock()()
	l_idx := Upperbound(len(g_ipv4secs), func(i int) bool { return g_ipv4secs[i].start > l_nip })
	if l_idx > 0 && l_nip <= g_ipv4secs[l_idx-1].end {
		fmt.Printf("[10]IP:%s,ret:%v\n", ip, in_chn)
		return in_chn
	}
	// loglv.Dbg.Printf("IsAcceptedIP(ip:%s) 处于第[%d.%d.%d.%d->%d.%d.%d.%d]区段", ip,
	// 	g_ipv4secs[l_idx-1].start>>24, g_ipv4secs[l_idx-1].start<<8>>24, g_ipv4secs[l_idx-1].start<<16>>24, g_ipv4secs[l_idx-1].start<<24>>24,
	// 	g_ipv4secs[l_idx-1].end>>24, g_ipv4secs[l_idx-1].end<<8>>24, g_ipv4secs[l_idx-1].end<<16>>24, g_ipv4secs[l_idx-1].end<<24>>24,
	// )
	if g_ipdb == nil {
		fmt.Printf("[11]IP:%s,ret:false\n", ip)
		return false //有一个生效和前面的两个都不生效的状况是不一样的
	}
	if ret, err := g_ipdb.Find(ip, "CN"); err != nil {
		fmt.Printf("[12]IP:%s,ret:true\n", ip)
		return true
	} else if len(ret) < 1 {
		fmt.Printf("[13]IP:%s,ret:%v\n", ip, !in_chn)
		return !in_chn
	} else if ret[0] != "中国" {
		fmt.Printf("[14]IP:%s,ret:%v\n", ip, !in_chn)
		return !in_chn
	} else if len(ret) < 2 {
		fmt.Printf("[15]IP:%s,ret:%v\n", ip, in_chn)
		return in_chn
	} else {
		fmt.Printf("[16]IP:%s,ret[1]:%s,in_chn:%v\n", ip, ret[1], in_chn)
		l_in_twhkmk := (ret[1] == "台湾") || (ret[1] == "香港") || (ret[1] == "澳门")
		return (in_chn && !l_in_twhkmk) || (!in_chn && l_in_twhkmk)
	}
}

//"chn_ip.txt"
func ReadIPset(file1, file2 string, f1mod, f2mod *time.Time, f1size, f2size *int64) (
	ret1ips []IPv4sec, ret2ips *ipdb.City, ret1err, ret2err error) {

	l_fi2, ret2err := os.Stat(file2)
	if ret2err != nil {
		fmt.Printf("文件操作os.Stat(%s)失败:%v\n", file2, ret2err)
	} else {
		l_modified := false
		if modtime := l_fi2.ModTime(); !f2mod.Equal(modtime) {
			*f2mod = modtime
			l_modified = true
		}
		if size := l_fi2.Size(); *f2size != size {
			*f2size = size
			l_modified = true
		}

		if l_modified {
			ret2ips, ret2err = ipdb.NewCity(file2)
			if ret2ips == nil || ret2err != nil {
				fmt.Printf("文件操作ipdb.NewCity(%s)失败:%v\n", file2, ret2err)
			}
		}
	}

	l_fd1, ret1err := os.Open(file1)
	if l_fd1 == nil || ret1err != nil {
		fmt.Printf("文件操作os.Open(%s)失败:%v\n", file1, ret1err)
		return
	}
	defer l_fd1.Close()
	if fi1, err := l_fd1.Stat(); err != nil {
		fmt.Printf("文件操作os.Stat(%s)失败:%v\n", file1, ret1err)
		ret1err = err
		return
	} else {
		l_modified := false
		if modtime := fi1.ModTime(); !f1mod.Equal(modtime) {
			*f1mod = modtime
			l_modified = true
		}
		if size := fi1.Size(); *f1size != size {
			*f1size = size
			l_modified = true
		}
		if !l_modified {
			return
		}
	}

	l_freader := bufio.NewReader(l_fd1)
	lf_readln := func() (string, error) {
		l_line, err := l_freader.ReadString('\n')
		if idx := strings.Index(l_line, "//"); idx >= 0 {
			l_line = l_line[0:idx]
		}
		return strings.TrimSpace(l_line), err
	}
	var l_lineno int
	for linetxt, err0 := lf_readln(); err0 == nil || err0 == io.EOF; linetxt, err0 = lf_readln() {
		l_lineno++
		l_ips := SplitFields(linetxt, ",;-", true)
		if len(l_ips) != 2 {
			if err0 == io.EOF {
				break
			}
			continue
		}
		var l_ipsec IPv4sec
		if true { //start
			l_ip, err := SplitInt16s(l_ips[0], ".", true)
			if err != nil || len(l_ip) != IPv4len {
				if ret1err == nil {
					ret1err = fmt.Errorf("第%d行格式错误", l_lineno)
				}
				if err0 == io.EOF {
					break
				}
				continue
			}
			for i := 0; i < IPv4len; i++ {
				l_ipsec.start += uint32(byte(l_ip[i])) << uint((3-i)*8)
			}
		}
		if true { //end
			l_ip, err := SplitInt16s(l_ips[1], ".", true)
			if err != nil || len(l_ip) != IPv4len {
				if ret1err == nil {
					ret1err = fmt.Errorf("第%d行格式错误", l_lineno)
				}
				if err0 == io.EOF {
					break
				}
				continue
			}
			for i := 0; i < IPv4len; i++ {
				l_ipsec.end += uint32(byte(l_ip[i])) << uint((3-i)*8)
			}
		}
		ret1ips = append(ret1ips, l_ipsec)

		if err0 == io.EOF {
			break
		}
	}
	if len(ret1ips) > 0 {
		//添加私有区段
		// ret1ips = append(ret1ips, IPv4sec{0x0a000000, 0x0aFFFFFF}, //10.x.x.x
		// 	IPv4sec{127 * (1 << 24), 127*(1<<24) + 0xFFFFFF},                       //127.x.x.x
		// 	IPv4sec{172*(1<<24) + 16*(1<<16), 172*(1<<24) + 31*(1<<16) + 0xFFFF},   //172.16.x.x-172.31.x.x
		// 	IPv4sec{192*(1<<24) + 168*(1<<16), 192*(1<<24) + 168*(1<<16) + 0xFFFF}, //192.168.x.x
		// )
		Sort(len(ret1ips),
			func(i, j int) bool {
				return ret1ips[i].start < ret1ips[j].start
			}, func(i, j int) {
				ret1ips[i], ret1ips[j] = ret1ips[j], ret1ips[i]
			})
		var l_deleted_idx []int
		TraverseFunc(len(ret1ips)-1,
			func(i int) int {
				if ret1ips[i].end+1 >= ret1ips[i+1].start {
					idx := RevIndexFunc(i+1, func(i int) bool { return ret1ips[i].start != 0xFFFFFFFF })
					if ret1ips[idx].end < ret1ips[i+1].end {
						ret1ips[idx].end = ret1ips[i+1].end
					}
					l_deleted_idx = append(l_deleted_idx, i+1)
					ret1ips[i+1].start = 0xFFFFFFFF
				}
				return 0
			})
		l_deleted_count := DelMulti(len(ret1ips), func(i, j int) { ret1ips[i] = ret1ips[j] }, l_deleted_idx...)
		ret1ips = ret1ips[:len(ret1ips)-l_deleted_count]
	}
	return
}
func Stackf(format string, v ...interface{}) {
	l_buf := [2048]byte{'\n'}
	l_stack := l_buf[:runtime.Stack(l_buf[1:len(l_buf)-2], false)+2]
	l_stack[len(l_stack)-1] = '\n'
	fmt.Printf("\n" + fmt.Sprintf(format, v...) + string(l_stack))
}
