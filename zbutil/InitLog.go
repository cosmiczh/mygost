package zbutil

import (
	"os"
	"time"

	"github.com/ginuerzh/gost/zbutil/loglv"
)

func InitLog(isdaemon bool) {
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
	f, e := os.OpenFile(l_logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if e == nil {
		loglv.SetOutput(f)
		if isdaemon {
			std2null()
		}
	}
}
