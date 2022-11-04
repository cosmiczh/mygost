package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ginuerzh/gost/zbutil"
	"github.com/ginuerzh/gost/zbutil/loglv"
)

func bufreset(buf *bytes.Buffer) *bytes.Buffer { buf.Reset(); return buf }
func correct_router() {
	var l_once_err1, l_once_err2 bool
	var l_out, l_err bytes.Buffer
	for time.Sleep(time.Second * 2); true; time.Sleep(time.Second * 5) {
		l_once_err1 = l_once_err2

		l_cmd := exec.Command("route", "-n")
		l_cmd.Stdout = bufreset(&l_out)
		l_cmd.Stderr = bufreset(&l_err)

		if err := l_cmd.Run(); err != nil {
			loglv.Err.Printf("CMD(route -n) error:%v:\t%s\n---将关闭矫正路由表功能-----", err, l_err.String())
			return
		}
		var l_lines = strings.FieldsFunc(l_out.String(), func(r rune) bool { return r == '\n' })

		var l_headln = map[string]int{}
		var l_headi = 0
		var l_defa_gw, l_defa_iface = "", ""
		for lineno, line := range l_lines {
			l_cols := strings.FieldsFunc(line, func(r rune) bool {
				return r == ' ' || r == '\t'
			})
			if len(l_cols) != 8 {
				continue
			} else if len(l_headln) == 0 { //准备head行
				for colno, txt := range l_cols {
					l_headln[txt] = colno
				}
				l_headi = lineno
				continue
			}
			l_dest := l_cols[l_headln["Destination"]+l_headln["目标"]]
			l_mask := l_cols[l_headln["Genmask"]+l_headln["子网掩码"]]
			l_iface := l_cols[l_headln["Iface"]+l_headln["接口"]]
			if !strings.HasPrefix(l_iface, "ppp") && l_dest == "0.0.0.0" && l_mask == "0.0.0.0" {
				l_defa_iface, l_defa_gw = l_iface,
					l_cols[l_headln["Gateway"]+l_headln["网关"]]
			}
		}
		for _, line := range l_lines[l_headi+1:] {
			l_cols := strings.FieldsFunc(line, func(r rune) bool {
				return r == ' ' || r == '\t'
			})
			if len(l_cols) != 8 {
				continue
			}
			l_dest := l_cols[l_headln["Destination"]+l_headln["目标"]]
			l_mask := l_cols[l_headln["Genmask"]+l_headln["子网掩码"]]
			l_iface := l_cols[l_headln["Iface"]+l_headln["接口"]]
			l_gateway := l_cols[l_headln["Gateway"]+l_headln["网关"]]
			//删除目标主机路由与缺省全网路由一样的多余路由
			if l_mask == "255.255.255.255" && l_iface == l_defa_iface && l_gateway == l_defa_gw {
				l_cmd := exec.Command("route", "del", "-host", l_dest)
				l_cmd.Stdout = bufreset(&l_out)
				l_cmd.Stderr = bufreset(&l_err)

				if err := l_cmd.Run(); err == nil {
					loglv.Inf.Printf("删除 [route del -host %s],output(%s)", l_dest, l_out.String())
				} else if errtxt := strings.TrimSpace(l_err.String()); !l_once_err1 {
					l_once_err2 = true
					loglv.Err.Printf("[ONCE]CMD(route del -host %s) error:%v\toutput(%s)", l_dest, err, errtxt)
				}
				continue
			} else if !strings.HasPrefix(l_iface, "ppp") || l_gateway == "0.0.0.0" {
				continue
			} else { //if strings.HasPrefix(l_iface, "ppp") && l_gateway != "0.0.0.0"
				//ppp接口上还有网关的路由为错误的ppp链路转发路由，删除之
				l_net := fmt.Sprintf("%s/%d", l_dest, calc_maskbitn(l_mask))
				l_cmd := exec.Command("route", "del", "-net", l_net)
				l_cmd.Stdout = bufreset(&l_out)
				l_cmd.Stderr = bufreset(&l_err)

				if err := l_cmd.Run(); err == nil {
					loglv.Inf.Printf("删除 [route del -net %s],output(%s)", l_net, l_out.String())
				} else if errtxt := strings.TrimSpace(l_err.String()); !l_once_err1 {
					l_once_err2 = true
					loglv.Err.Printf("[ONCE]CMD(route del -net %s) error:%v\toutput(%s)", l_net, err, errtxt)
				}
			}

			//增加到ppp拨号的远端子网路由
			switch cmdrun := func(subnet string) {
				l_cmd := exec.Command("route", "add", "-net", subnet, l_iface)
				l_cmd.Stdout = bufreset(&l_out)
				l_cmd.Stderr = bufreset(&l_err)

				if err := l_cmd.Run(); err == nil {
					loglv.Inf.Printf("增加 [route add -net %s %s],output(%s)", subnet, l_iface, l_out.String())
				} else if errtxt := strings.TrimSpace(l_err.String()); strings.Index(errtxt, "File exists") >= 0 {
					loglv.Inf.Printf("增加 [route add -net %s %s],output(%s)", subnet, l_iface, errtxt)
				} else if !l_once_err1 {
					l_once_err2 = true
					loglv.Err.Printf("[ONCE]CMD(route add -net %s %s) error:%v\toutput(%s)", subnet, l_iface, err, errtxt)
				}
			}; { //switch cmdrun :=
			case strings.HasPrefix(l_gateway, "192.168."):
				subnet := l_gateway[:strings.LastIndexByte(l_gateway, '.')] + ".0/24"
				cmdrun(subnet)
			case strings.HasPrefix(l_gateway, "172."), strings.HasPrefix(l_gateway, "10."):
				subnet := l_gateway[:strings.LastIndexByte(l_gateway, '.')] + ".0"
				cmdrun(fmt.Sprintf("%s/%d", subnet, zbutil.MaxInt(calc_maskbitn(subnet), 16)))
			}
		}
	}
}

// 提取子网掩码位数
func calc_maskbitn(v4dot string) int {
	l_v4dot := strings.FieldsFunc(v4dot, func(r rune) bool { return r == '.' })
	if len(l_v4dot) > 4 {
		return 32
	}
	var l_Genmask uint32 = 0
	for _, sec := range l_v4dot {
		l_Genmask <<= 8
		sec, _ := strconv.Atoi(sec)
		l_Genmask += uint32(sec)
	}
	l_Genmask <<= ((4 - len(l_v4dot)) * 8)
	l_maskbitn := 0 //计算子网掩码位数
	for ; l_Genmask != 0; l_Genmask <<= 1 {
		l_maskbitn++
	}
	return l_maskbitn
}
