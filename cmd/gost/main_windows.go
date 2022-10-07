package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"

	_ "net/http/pprof"

	"github.com/ginuerzh/gost"
	"github.com/go-log/log"
	"github.com/kardianos/service"
)

var (
	configureFile string
	baseCfg       = &baseConfig{}
	pprofAddr     string
	pprofEnabled  = os.Getenv("PROFILING") != ""

	install, remove, runsvc bool
	svcname                 string
)

func init() {
	gost.SetLogger(&gost.LogLogger{})
	var l_delidx []int
	for i, arg := range os.Args {
		switch arg {
		case "-install":
			l_delidx, install = append(l_delidx, i), true
			if i+1 < len(os.Args) && os.Args[i+1][0] != '-' {
				l_delidx = append(l_delidx, i+1)
				svcname = os.Args[i+1]
			}
		case "-remove":
			l_delidx, remove = append(l_delidx, i), true
			if i+1 < len(os.Args) && os.Args[i+1][0] != '-' {
				l_delidx = append(l_delidx, i+1)
				svcname = os.Args[i+1]
			}
		case "-runsvc":
			l_delidx, runsvc = append(l_delidx, i), true
			if i+1 < len(os.Args) && os.Args[i+1][0] != '-' {
				l_delidx = append(l_delidx, i+1)
				svcname = os.Args[i+1]
			}
		}
	}
	if len(l_delidx) > 0 {
		l_delcount := DelMulti(len(os.Args), func(i, j int) { os.Args[i] = os.Args[j] }, l_delidx...)
		os.Args = os.Args[:len(os.Args)-l_delcount]
	}
	if install {
		create_svc().Install()
		log.Log("服务安装成功")
		return
	} else if remove {
		create_svc().Uninstall()
		log.Log("服务卸载成功")
		return
	} else if runsvc {
		for i, arg := range os.Args {
			arg = strings.ReplaceAll(arg, "\\\"", "")
			os.Args[i] = strings.ReplaceAll(arg, "\"", "")
		}
	}

	var (
		printVersion bool
		_lfname      stringList
	)

	flag.String("install", "", "install windows service")
	flag.String("remove", "", "remove windows service")
	flag.String("runsvc", "", "running windows service in service control")
	flag.IntVar(&baseCfg.route.Mark, "M", 0, "Specify out connection mark")
	flag.Var(&baseCfg.route.ChainNodes, "F", "forward address, can make a forward chain")
	flag.Var(&baseCfg.route.ServeNodes, "L", "listen address, can listen on multiple ports (required)")
	flag.Var(&_lfname, "LF", "file name which read options -F and -L from")
	flag.Var(&_lfname, "FL", "file name which read options -F and -L from")
	flag.StringVar(&configureFile, "C", "", "configure file")
	flag.BoolVar(&baseCfg.Debug, "D", false, "enable debug log")
	flag.BoolVar(&printVersion, "V", false, "print version")
	if pprofEnabled {
		flag.StringVar(&pprofAddr, "P", ":6060", "profiling HTTP server address")
	}
	flag.Parse()

	if printVersion {
		fmt.Fprintf(os.Stdout, "gost %s (%s %s/%s)\n",
			gost.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if configureFile != "" {
		_, err := parseBaseConfig(configureFile)
		if err != nil {
			log.Log(err)
			os.Exit(1)
		}
	}

	for _, fn := range _lfname {
		L, F := parseLF(fn)
		baseCfg.route.ServeNodes = append(baseCfg.route.ServeNodes, L...)
		baseCfg.route.ChainNodes = append(baseCfg.route.ChainNodes, F...)
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
func main() {
	if install || remove {
		return
	}
	if pprofEnabled {
		go func() {
			log.Log("profiling server on", pprofAddr)
			log.Log(http.ListenAndServe(pprofAddr, nil))
		}()
	}

	// NOTE: as of 2.6, you can use custom cert/key files to initialize the default certificate.
	tlsConfig, err := tlsConfig(defaultCertFile, defaultKeyFile, "")
	if err != nil {
		// generate random self-signed certificate.
		cert, err := gost.GenCertificate()
		if err != nil {
			log.Log(err)
			os.Exit(1)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	} else {
		log.Log("load TLS certificate files OK")
	}

	gost.DefaultTLSConfig = tlsConfig

	if runsvc {
		create_svc().Run()
		return
	}

	if err := start(); err != nil {
		log.Log(err)
		os.Exit(1)
	}

	select {}
}

var g_routers []router

func start() error {
	gost.Debug = baseCfg.Debug

	rts, err := baseCfg.route.GenRouters()
	if err != nil {
		return err
	}
	g_routers = append(g_routers, rts...)

	for _, route := range baseCfg.Routes {
		rts, err := route.GenRouters()
		if err != nil {
			return err
		}
		g_routers = append(g_routers, rts...)
	}

	if len(g_routers) == 0 {
		return errors.New("invalid config")
	}
	for i := range g_routers {
		go g_routers[i].Serve()
	}

	return nil
}

type program struct{}

func (p *program) Start(s service.Service) error {
	err := start()
	if err != nil {
		log.Log(err)
	}
	return err
}
func (p *program) Stop(s service.Service) error {
	for i := range g_routers {
		g_routers[i].Close()
	}
	g_routers = nil
	return nil
}

func create_svc() service.Service {
	svcConfig := &service.Config{
		Name: func() string { //服务显示名称
			if len(svcname) < 1 {
				return "gost"
			} else {
				return "gost_" + svcname
			}
		}(),
		DisplayName: func() string { //服务名称
			if len(svcname) < 1 {
				return "gost proxy"
			} else {
				return "gost_" + svcname + " proxy"
			}
		}(),
		Description: "https://docs.ginuerzh.xyz/gost/", //服务描述
		Arguments: func() []string {
			if len(svcname) < 1 {
				return append([]string{"-runsvc"}, os.Args[1:]...)
			} else {
				return append([]string{"-runsvc", svcname}, os.Args[1:]...)
			}
		}(),
	}
	for i, arg := range svcConfig.Arguments {
		l_isFL := false
		if len(arg) >= 3 {
			switch arg[:3] {
			case "-LF", "-FL":
				if l_isFL = true; len(arg) == 3 {
					if i+1 < len(svcConfig.Arguments) && svcConfig.Arguments[i+1][0] != '"' && svcConfig.Arguments[i+1][0] != '-' {
						svcConfig.Arguments[i+1] = "\"" + svcConfig.Arguments[i+1] + "\""
					}
				} else if arg[3] != '=' {
					svcConfig.Arguments[i] = arg[:3] + "\"" + arg[3:] + "\""
				} else if len(arg) >= 5 {
					svcConfig.Arguments[i] = arg[:4] + "\"" + arg[4:] + "\""
				}
			}
		}
		if !l_isFL && len(arg) >= 2 {
			switch arg[:2] {
			case "-L", "-F":
				if len(arg) == 2 {
					if i+1 < len(svcConfig.Arguments) && svcConfig.Arguments[i+1][0] != '"' && svcConfig.Arguments[i+1][0] != '-' {
						svcConfig.Arguments[i+1] = "\"" + svcConfig.Arguments[i+1] + "\""
					}
				} else if arg[2] != '=' {
					svcConfig.Arguments[i] = arg[:2] + "\"" + arg[2:] + "\""
				} else if len(arg) >= 4 {
					svcConfig.Arguments[i] = arg[:3] + "\"" + arg[3:] + "\""
				}
			}
		}
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Log(err)
		return nil
	}
	return s
}
func DelMulti(nlen int, fassign func(i, j int), deletedidx ...int) (deletedcount int) {
	if len(deletedidx) < 1 {
		return 0
	} else if len(deletedidx) > 2 {
		deletedidx = UniqueInt(deletedidx...)
	}
	for i := 0; i < len(deletedidx); i++ {
		if deletedidx[i] >= 0 {
			deletedidx = deletedidx[i:]
			break
		}
	}
	if len(deletedidx) < 1 {
		return 0
	}
	if deletedidx[len(deletedidx)-1] < nlen {
		deletedidx = append(deletedidx, nlen) //末尾一定要追加一个大数
	}
	deletedcount = 0
	for i := deletedidx[0]; i < nlen-deletedcount; i++ {
		if i == deletedidx[deletedcount]-deletedcount {
			deletedcount++
			i--
			continue
		}
		fassign(i, i+deletedcount)
	}
	return
}
func Unique2(nlen int, fless func(i, j int) bool, fswap func(i, j int), fequal func(i, j int) bool, fassign func(i, j int)) (deletedcount int) {
	if nlen < 2 {
		return 0
	}
	Sort(nlen, fless, fswap)

	deletedcount = 0
	for i := 1; i < nlen-deletedcount; i++ {
		if fequal(i-1, i+deletedcount) {
			deletedcount++
			i--
			continue
		}
		if deletedcount > 0 {
			fassign(i, i+deletedcount)
		}
	}
	return
}

func UniqueInt(arr ...int) []int {
	l_deletedcount := Unique2(len(arr),
		func(i, j int) bool { return arr[i] < arr[j] },
		func(i, j int) { arr[i], arr[j] = arr[j], arr[i] },
		func(i, j int) bool { return arr[i] == arr[j] },
		func(i, j int) { arr[i] = arr[j] })
	return arr[:len(arr)-l_deletedcount]
}

type struct4sort struct {
	m_len   int
	mf_less func(i, j int) bool
	mf_swap func(i, j int)
}

func (this *struct4sort) Len() int {
	return this.m_len
}
func (this *struct4sort) Less(i, j int) bool {
	return this.mf_less(i, j)
}
func (this *struct4sort) Swap(i, j int) {
	this.mf_swap(i, j)
}

func Sort(nlen int, fless func(i, j int) bool, fswap func(i, j int)) {
	sort.Sort(&struct4sort{m_len: nlen, mf_less: fless, mf_swap: fswap})
}
