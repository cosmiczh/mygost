package loglv

import (
	"log"
	"os"
	"time"
)

var GetLogDir func() string
var GetExeBaseName func() string
var CmdParmLike func(leftpart string) string
var SearchFile func(plist_file *[]os.FileInfo, dirname string, name_pattern string) ([]os.FileInfo, error)

func InitLOG(isdaemon bool) {
	l_logfile := ""
	if l_logfile = CmdParmLike("-log="); len(l_logfile) > 0 {
		l_logfile = l_logfile[len("-log="):]
	} else if l_logfile = CmdParmLike("--logfile="); len(l_logfile) > 0 {
		l_logfile = l_logfile[len("--logfile="):]
	}
	if len(l_logfile) < 1 {
		os.Mkdir(GetLogDir(), 0755)
		l_logfile = GetLogDir() + "/" + GetExeBaseName() + ".log"
		l_bakfile := GetLogDir() + "/" + time.Now().Format("06-01-02_15.04.05_") + GetExeBaseName() + ".log"
		os.Rename(l_logfile, l_bakfile)

		l_osfils, _ := SearchFile(nil, GetLogDir(), GetExeBaseName()+".out")
		if len(l_osfils) > 0 && l_osfils[0].Size() > 0 {
			l_logfile = GetLogDir() + "/" + GetExeBaseName() + ".out"
			l_bakfile = GetLogDir() + "/" + time.Now().Format("06-01-02_15.04.05_") + GetExeBaseName() + ".out"
			os.Rename(l_logfile, l_bakfile)
		}
	}
	if err := SetOutput("", l_logfile, nil); err == nil && isdaemon {
		std2null()
	} else {
		log.Printf("----------daemon:%v,err:%v-----------------------", isdaemon, err)
	}
}
