package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	_ "net/http/pprof"

	"github.com/ginuerzh/gost"
	"github.com/ginuerzh/gost/zbutil"
	"github.com/go-log/log"
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
		l_delcount := zbutil.DelMulti(len(os.Args), func(i, j int) { os.Args[i] = os.Args[j] }, l_delidx...)
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
