package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	// "github.com/ginuerzh/gost/zbutil/loglv"
	"github.com/ginuerzh/gost/zbutil"
)

func GetOS() string { return "linux" }

type StrBuilder struct{ strings.Builder }

func NewBuilder(buf []byte) (ret *StrBuilder) {
	ret = &StrBuilder{}
	ret.Write(buf)
	return
}
func NewBuilderString(s string) (ret *StrBuilder) {
	ret = &StrBuilder{}
	ret.WriteString(s)
	return
}

type GetIOinfo interface {
	GetSockNum() int
	GetSendBytes() int64
	GetRecvBytes() int64
}

var (
	appName   string
	envName   string
	pidFile   string
	pidVal    int   = 0
	isdaemon  bool  = true
	stop_tmot int32 = 30 //服务stop的等待超时秒数，缺省30s，由命令行-f=xxx指定
)

func MainStartup(workdir string) func() {
	appName = zbutil.GetExeBaseName()
	pidFile = zbutil.GetExeDir() + "/pid/" + zbutil.GetExeBaseName() + ".pid"
	envName = zbutil.GetExeDir() + "/" + zbutil.GetExeName() + "__Daemon"

	if tmot := zbutil.CmdParmLike("-f="); len(tmot) > 0 {
		if tmot, err := zbutil.Atoi(tmot[len("-f="):]); err == nil {
			stop_tmot = tmot
		}
	}

forloop:
	for os.Getenv(envName) != "true" { //master
		switch {
		case len(zbutil.CmdParmLike("start")) == len("start"):
			if isRunning() {
				fmt.Printf("[%d] %s is running\n", pidVal, appName)
			} else { //fork daemon进程
				svcstart()
			}
		case len(zbutil.CmdParmLike("restart")) == len("restart"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
				svcstart()
			} else {
				fmt.Printf("[%d] %s restart now\n", pidVal, appName)
				restart(pidVal)
			}
		case len(zbutil.CmdParmLike("stop")) == len("stop"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else {
				svcstop(pidVal)
			}
		case len(zbutil.CmdParmLike("stat")) == len("stat"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else {
				fmt.Printf("[%d] %s is running\n", pidVal, appName)
			}
		case len(zbutil.CmdParmLike("sig47")) == len("sig47"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else if syscall.Kill(pidVal, 47) != nil {
				fmt.Printf("[%d] %s process not exist\n", pidVal, appName)
			} else {
				fmt.Printf("[%d] %s sig47(stackall) be sent successfully\n", pidVal, appName)
			}
		case len(zbutil.CmdParmLike("sig48")) == len("sig48"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else if syscall.Kill(pidVal, 48) != nil {
				fmt.Printf("[%d] %s process not exist\n", pidVal, appName)
			} else {
				fmt.Printf("[%d] %s sig48(cpuprof) be sent successfully\n", pidVal, appName)
			}
		case len(zbutil.CmdParmLike("sig50")) == len("sig50"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else if syscall.Kill(pidVal, 50) != nil {
				fmt.Printf("[%d] %s process not exist\n", pidVal, appName)
			} else {
				fmt.Printf("[%d] %s sig50(heaprof) be sent successfully\n", pidVal, appName)
			}
		case len(zbutil.CmdParmLike("sig51")) == len("sig51"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else if syscall.Kill(pidVal, 51) != nil {
				fmt.Printf("[%d] %s process not exist\n", pidVal, appName)
			} else {
				fmt.Printf("[%d] %s sig51(memlog) be sent successfully\n", pidVal, appName)
			}
		case len(zbutil.CmdParmLike("sig52")) == len("sig52"):
			if !isRunning() {
				fmt.Printf("%s not running\n", appName)
			} else if syscall.Kill(pidVal, 52) != nil {
				fmt.Printf("[%d] %s process not exist\n", pidVal, appName)
			} else {
				fmt.Printf("[%d] %s sig52(memlog switch interval) be sent successfully\n", pidVal, appName)
			}
		case len(zbutil.CmdParmLike("-h")) == len("-h"):
			fmt.Printf("Usage: %s [-log={logfile}]] [start|restart|stop|stat|sig48(cpuprof)|sig50(heaprof)|sig51(memlog)|sig52(memlog sw)]\n", appName)
		default:
			isdaemon = false
			break forloop
		}
		os.Exit(0)
	}
	os.Setenv(envName, "false")

	if len(workdir) > 0 {
		os.Chdir(workdir)
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("[%d] child report:%s\n", os.Getpid(), zbutil.GetExeDir()+"/"+zbutil.GetExeName())

	zbutil.InitLog(isdaemon)
	return func() { os.Remove(pidFile) }
}

//检查pidFile是否存在以及文件里的pid是否存活
func isRunning() bool {
	if mf, err := os.Open(pidFile); err == nil {
		cpid, err := ioutil.ReadAll(mf)
		if err != nil {
			pidVal = -1
		} else if pidVal, err = strconv.Atoi(string(cpid)); err != nil {
			pidVal = -1
		}
	}

	if pidVal > 0 {
		if err := syscall.Kill(pidVal, 0); err == nil { //发一个信号为0到指定进程ID，如果没有错误发生，表示进程存活
			return true
		}
	}
	os.Remove(pidFile)
	return false
}

func svcstart() {
	os.Setenv(envName, "true")

	func() { //保存pid
		os.Mkdir(zbutil.GetExeDir()+"/pid", 0755)
		file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("pidFile[%s] creation failure:%v\n", pidFile, err)
			return
		}
		defer file.Close()
	}()

	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	cpid, err := syscall.ForkExec(zbutil.GetExeDir()+"/"+zbutil.GetExeName(), os.Args, procAttr)
	if err != nil {
		fmt.Printf("%s start failure:%v,%s\n", appName, err, zbutil.GetExeDir()+"/"+zbutil.GetExeName())
		os.Exit(-2)
	}
	fmt.Printf("[%d] parent report child[%d]:%s\n", os.Getpid(), cpid, zbutil.GetExeDir()+"/"+zbutil.GetExeName())

	//wait for child startup
	l_printdot, l_had_printed := false, false
	for {
		f, err := os.Open(pidFile)
		if err == nil {
			pidVal, _ := ioutil.ReadAll(f)
			f.Close()
			if strconv.Itoa(cpid) == string(pidVal) {
				break //启动完全正常
			}
			if l_printdot {
				fmt.Print(".")
				l_had_printed = true
			}
		} else if os.IsNotExist(err) || os.IsPermission(err) { //文件还不存在
			if l_had_printed {
				fmt.Printf("\n")
			}
			fmt.Printf("[%d] %s start failure.----------------------\n\n", cpid, appName)
			os.Exit(-3)
		} else {
			if l_printdot {
				fmt.Print(".")
				l_had_printed = true
			}
		}
		time.Sleep(500 * time.Millisecond)
		if !l_printdot {
			l_printdot = true
		}
	}

	if l_had_printed {
		fmt.Printf("\n")
	}
	fmt.Printf("[%d] %s start daemon.----------------------\n\n", cpid, appName)
}
func waitexit(cpid int) int {
	l_printdot := false
	l_begin := time.Now()
	for { //循环，查看pidFile是否存在，不存在或值已改变，返回
		f, err := os.Open(pidFile)
		if err == nil {
			pidVal, _ := ioutil.ReadAll(f)
			f.Close()
			if strconv.Itoa(cpid) != string(pidVal) {
				return 0
			}
			if l_printdot {
				fmt.Print(".")
			}
		} else if os.IsNotExist(err) || os.IsPermission(err) { //文件已不存在
			if syscall.Kill(cpid, 0) != nil { //发一个信号0到指定进程ID，如果有错误发生，表示进程真的不存在了
				return 1
			}
			if l_printdot {
				fmt.Print(".")
			}
		} else {
			fmt.Println("err:", err)
		}
		if time.Now().Sub(l_begin) >= time.Duration(stop_tmot)*time.Second {
			if len(zbutil.CmdParmLike("-f")) >= len("-f") {
				return 2
			}
			return 2 //-1
		}
		time.Sleep(500 * time.Millisecond)
		if !l_printdot {
			l_printdot = true
		}
	}
}

//重启(先发送kill -HUP到运行进程，手工重启daemon ...当有运行的进程时，daemon不启动)
func restart(cpid int) {
	syscall.Kill(cpid, syscall.SIGHUP) //kill -HUP, daemon only时，会直接退出
	//处理结果
	switch waitexit(cpid) {
	case 2:
		syscall.Kill(cpid, syscall.SIGKILL)
		fmt.Printf("\n[%d] %s %ds timeout and killed. restarting...\n", cpid, appName, stop_tmot)
		os.Remove(pidFile)
		svcstart()
	case 1:
		fmt.Printf("\n[%d] %s restarting...\n", cpid, appName)
		svcstart()
	case 0:
		fmt.Printf("\n[%d] %s restart failure.\n\n", cpid, appName)
		os.Exit(-2)
	case -1:
		fmt.Printf("\n[%d] %s restart timeout.\n\n", cpid, appName)
		os.Exit(-3)
	}
}
func svcstop(cpid int) {
	syscall.Kill(cpid, syscall.SIGTERM) //kill
	fmt.Printf("\n[%d] %s stopping...\n", cpid, appName)
	//处理结果
	switch waitexit(cpid) {
	case 2:
		syscall.Kill(cpid, syscall.SIGKILL)
		fmt.Printf("\n[%d] %s %ds timeout and killed.\n", cpid, appName, stop_tmot)
		os.Remove(pidFile)
	case 1:
		fmt.Printf("\n[%d] %s stopped.\n", cpid, appName)
	case 0:
		fmt.Printf("\n[%d] %s stop failure.\n", cpid, appName)
		os.Exit(-2)
	case -1:
		fmt.Printf("\n[%d] %s stop timeout.\n", cpid, appName)
		os.Exit(-3)
	}
}

func MainWait(getio GetIOinfo) {
	if !isdaemon {
		var l_prev_recvbts, l_prev_sendbts int64 = 0, 0
		l_prev_now := time.Now()
		var l_instr string
		for l_instr != "exit" {
			time.Sleep(time.Duration(time.Millisecond * 100))
			fmt.Printf("请输入exit退出：")
			fmt.Scan(&l_instr)
			if l_instr == "getinfo" {
				fmt.Println("连接数:", getio.GetSockNum())
				l_recvbts, l_sendbts := getio.GetRecvBytes(), getio.GetSendBytes()
				fmt.Printf("接收:%d(KB)\t 发送:%d(KB)\n", (l_recvbts-l_prev_recvbts)/1024, (l_sendbts-l_prev_sendbts)/1024)
				l_now := time.Now()
				l_delta_sec := int64(l_now.Sub(l_prev_now) * 1024 / time.Second)
				fmt.Printf("接收:%d(KB/s)\t 发送:%d(KB/s)\n", (l_recvbts-l_prev_recvbts)/l_delta_sec, (l_sendbts-l_prev_sendbts)/l_delta_sec)
				l_prev_recvbts, l_prev_sendbts, l_prev_now = l_recvbts, l_sendbts, l_now
			}
		}
	} else {
		if !func() bool { //保存pid
			file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("pidFile[%s] creation failure:%v\n", pidFile, err)
				return false
			}
			defer file.Close()
			if _, err := file.WriteString(strconv.Itoa(os.Getpid())); err != nil {
				fmt.Printf("pidFile[%s] write pid failure:%v\n", pidFile, err)
				return false
			}
			return true
		}() {
			return
		}
		l_cpuf_fp, l_cpuf_name := (*os.File)(nil), ""
		lf_startcpuprof := func() {
			// defer loglv.Err.Recoverf("start cpuprof panic.")
			l_cpuf_name = fmt.Sprintf("%s/%s_%s.cpu.pprof", zbutil.GetLogDir(), zbutil.GetExeBaseName(), time.Now().Format("06-01-02_15.04.05"))

			os.Mkdir(zbutil.GetLogDir(), 0755)
			l_err := error(nil)
			if l_cpuf_fp, l_err = os.Create(l_cpuf_name); l_err != nil {
				/*loglv.Err.*/ fmt.Printf("\n\tCan't create CPU profile[%s] error:%v", l_cpuf_name, l_err)
				l_cpuf_fp, l_cpuf_name = nil, ""
				return
			}
			if l_err = pprof.StartCPUProfile(l_cpuf_fp); l_err != nil { //监控cpu
				/*loglv.Err.*/ fmt.Printf("\n\tCan't start CPU profile[%s] error:%v", l_cpuf_name, l_err)
				l_cpuf_fp.Close()
				os.Remove(l_cpuf_name)
				l_cpuf_fp, l_cpuf_name = nil, ""
				return
			}

			/*loglv.War.*/
			fmt.Printf("\n\tCPU Profile[%s] started.", l_cpuf_name)
		}
		lf_closecpuprof := func() {
			if len(l_cpuf_name) > 0 {
				// defer loglv.Err.Recoverf("stop cpuprof panic.")
				pprof.StopCPUProfile()
				l_cpuf_fp.Close()

				/*loglv.War.*/
				fmt.Printf("\n\tCPU Profile[%s] stopped.", l_cpuf_name)
				l_cpuf_fp, l_cpuf_name = nil, ""
			}
		}
		defer lf_closecpuprof()

		l_memlog_waitgrp, l_gap_chan, l_gap_switch := (*sync.WaitGroup)(nil), (chan time.Duration)(nil), 0
		lf_closememlog := func() {
			if l_memlog_waitgrp != nil {
				// defer loglv.Err.Recoverf("stop memlog panic.")
				close(l_gap_chan)
				l_memlog_waitgrp.Wait()
				/*loglv.War.*/ fmt.Printf("\n\tMEMlog stopped.")
				l_memlog_waitgrp, l_gap_chan, l_gap_switch = nil, nil, 0
			}
		}
		defer lf_closememlog()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT,
			syscall.Signal(47), syscall.Signal(48), syscall.Signal(50), syscall.Signal(51), syscall.Signal(52))
	FOR_LOOP:
		for {
			switch sig := <-sigs; sig {
			case syscall.SIGINT: //2:Ctrl+C
				break FOR_LOOP
			case syscall.SIGQUIT: //3:Ctrl+/
				break FOR_LOOP
			case syscall.SIGABRT: //6:调用abort函数
				break FOR_LOOP
			case syscall.SIGTERM: //15
				break FOR_LOOP
			case syscall.SIGHUP:
				// break FOR_LOOP	//某些路由器嵌入式linux终端退出daemon会收到这个挂起信号
			case syscall.Signal(47):
				l_file := zbutil.GetLogDir() + "/stack_" + zbutil.GetExeBaseName() + time.Now().Format("_06-01-02_15.04.05") + ".log"
				f, e := os.OpenFile(l_file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
				if e != nil {
					/*loglv.Err.*/ fmt.Printf("OpenStackFile(%s) error:%v", l_file, e)
				} else {
					// l_buf := [40*1024*1024 + 2]byte{'\n'}
					l_buf := make([]byte, 40*1024*1024+2)
					l_buf[0] = '\n'
					l_stack := l_buf[:runtime.Stack(l_buf[1:len(l_buf)-2], true)+2]
					l_stack[len(l_stack)-1] = '\n'
					f.Write(l_stack)
					f.Sync()
					f.Close()
				}
			case syscall.Signal(48): //cpuprof
				func() {
					if len(l_cpuf_name) < 1 { //开始cpuprof
						lf_startcpuprof()
					} else { //结束cpuprof
						lf_closecpuprof()
					}
				}()
			case syscall.Signal(50): //heapprof
				func() {
					// defer loglv.Err.Recoverf("write memprof panic.")
					l_memf_fp, l_memf_name := (*os.File)(nil),
						fmt.Sprintf("%s/%s_%s.mem.pprof", zbutil.GetLogDir(), zbutil.GetExeBaseName(), time.Now().Format("06-01-02_15.04.05"))

					l_err := error(nil)
					if l_memf_fp, l_err = os.Create(l_memf_name); l_err != nil {
						/*loglv.Err.*/ fmt.Printf("\n\tCan't create MEM profile[%s] error:%v", l_memf_name, l_err)
						l_memf_name, l_memf_fp = "", nil
						return
					}
					if l_err = pprof.WriteHeapProfile(l_memf_fp); l_err != nil { //写入heap
						/*loglv.Err.*/ fmt.Printf("\n\tCan't write MEM profile[%s] error:%v", l_memf_name, l_err)
						l_memf_fp.Close()
						os.Remove(l_memf_name)
						l_memf_fp, l_memf_name = nil, ""
						return
					}
					l_memf_fp.Close()
					l_memf_fp, l_memf_name = nil, ""
				}()
			case syscall.Signal(51): //memlog
				func() {
					if l_memlog_waitgrp == nil { //开始memlog
						// defer loglv.Err.Recoverf("start memlog panic.")
						l_memlog_waitgrp = new(sync.WaitGroup)
						l_memlog_waitgrp.Add(1)
						l_gap_switch = 1 //初始设置为10秒间隔
						l_gap_chan = make(chan time.Duration, 1)
						l_gap_chan <- 10 * time.Second
						go logmemstat(l_memlog_waitgrp, l_gap_chan)

						/*loglv.War.*/
						fmt.Printf("\n\tMEMlog started.")
					} else { //结束memlog
						lf_closememlog()
					}
				}()
			case syscall.Signal(52): //memlog间隔切换
				if l_memlog_waitgrp == nil {
					/*loglv.War.*/ fmt.Println("\n\tNot start memlog.")
				} else if l_gap_switch > 0 { //已经是10秒间隔，切换到5秒间隔
					l_gap_switch = 0
					l_gap_chan <- 5 * time.Second
					/*loglv.War.*/ fmt.Println("\n\tswitch to interval 5s.")
				} else if l_gap_switch == 0 { //已经是5秒间隔，切换到2秒间隔
					l_gap_switch = -1
					l_gap_chan <- 2 * time.Second
					/*loglv.War.*/ fmt.Println("\n\tswitch to interval 2s.")
				} else { //已经是2秒间隔，切换到10秒间隔
					l_gap_switch = 1
					l_gap_chan <- 10 * time.Second
					/*loglv.War.*/ fmt.Println("\n\tswitch to interval 10s.")
				}
			}
		}
	}
}

func logmemstat(waitgrp *sync.WaitGroup, gap_chan chan time.Duration) {
	defer waitgrp.Done()
	os.Mkdir(zbutil.GetLogDir(), 0755)

	l_memlogfile := fmt.Sprintf("%s/%s_%s.mem.log", zbutil.GetLogDir(), zbutil.GetExeBaseName(), time.Now().Format("06-01-02_15.04.05"))

	f, e := os.OpenFile(l_memlogfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if e != nil {
		/*loglv.Err.*/ fmt.Printf("\n\t Open memory stat file[%s] failure.error:%v", l_memlogfile, e)
		return
	}
	defer f.Close()

	l_gap := time.Duration(1 * time.Second)
	get_sign := func() (ok bool) {
		select {
		case l_gap, ok = <-gap_chan:
		case <-time.After(l_gap):
			ok = true
		}
		return
	}

	var l_stat runtime.MemStats
	for ok := get_sign(); ok; ok = get_sign() {
		runtime.GC()
		// debug.FreeOSMemory()
		runtime.ReadMemStats(&l_stat)
		f.WriteString(fmt.Sprintf("%s Sys:%02dg%03dm%03dk%03d,\tHeapSys:%02dg%03dm%03dk%03d,"+
			"\tHeapAlloc:%02dg%03dm%03dk%03d,\tHeapIdle:%02dg%03dm%03dk%03d,\tHeapReleased:%02dg%03dm%03dk%03d\n",

			time.Now().Format("06-01-02_15.04.05"),
			l_stat.Sys/(1<<30), l_stat.Sys%(1<<30)/(1<<20), l_stat.Sys%(1<<20)/(1<<10), l_stat.Sys%(1<<10),
			l_stat.HeapSys/(1<<30), l_stat.HeapSys%(1<<30)/(1<<20), l_stat.HeapSys%(1<<20)/(1<<10), l_stat.HeapSys%(1<<10),
			l_stat.HeapAlloc/(1<<30), l_stat.HeapAlloc%(1<<30)/(1<<20), l_stat.HeapAlloc%(1<<20)/(1<<10), l_stat.HeapAlloc%(1<<10),
			l_stat.HeapIdle/(1<<30), l_stat.HeapIdle%(1<<30)/(1<<20), l_stat.HeapIdle%(1<<20)/(1<<10), l_stat.HeapIdle%(1<<10),
			l_stat.HeapReleased/(1<<30), l_stat.HeapReleased%(1<<30)/(1<<20), l_stat.HeapReleased%(1<<20)/(1<<10), l_stat.HeapReleased%(1<<10),
		))
	}
}
