package main

import (
	"errors"
	"os"

	"github.com/ginuerzh/gost"
	"github.com/ginuerzh/gost/zbutil"
	"github.com/go-log/log"
	"github.com/kardianos/service"
)

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
	zbutil.InitLog(true)
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
