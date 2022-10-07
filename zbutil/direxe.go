package zbutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var g_exe_basedir string = ""

func init() {
	if g_exe_basedir != "" {
		return
	}
	g_exe_basedir = GetExeDir()
}

func GetExeDir() string {
	if g_exe_basedir != "" {
		return g_exe_basedir
	}
	l_execPath, l_err := exec.LookPath(os.Args[0])
	if l_err != nil {
		if strings.IndexAny(os.Args[0], "/\\") >= 0 {
			l_execPath = os.Args[0]
		} else {
			return "."
		}
	}
	// fmt.Printf("l_execPath1:%s\n", l_execPath)
	//    Is Symlink
	l_fileinfo, l_err := os.Lstat(l_execPath)
	if l_err != nil {
		return "."
	}
	if l_fileinfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		l_lnkPath, l_err := os.Readlink(l_execPath)
		// fmt.Printf("l_lnkPath2:%s\n", l_lnkPath)
		if l_err != nil {
			return "."
		}
		if l_lnkPath[0] == '/' {
			l_execPath = l_lnkPath
		} else {
			l_execPath = filepath.Dir(l_execPath) + "/" + l_lnkPath
		}
	}
	l_execDir := filepath.Dir(l_execPath)
	// fmt.Printf("l_execDir3:%s\n", l_execDir)
	if l_execDir == "." {
		l_execDir, l_err = os.Getwd()
		// fmt.Printf("l_execDir4:%s\n", l_execDir)
		if l_err != nil {
			return "."
		}
	}
	if len(l_execDir) < 1 || l_execDir == "." {
		l_execDir, _ = filepath.Abs(os.Args[0])
		// fmt.Printf("l_execDir5:%s\n", l_execDir)
		l_execDir = filepath.Dir(l_execDir)
	} else if (l_execDir[0] != '/') && (len(l_execDir) < 2 || l_execDir[1] != ':') {
		l_execDir, _ = filepath.Abs(l_execDir)
	}
	// fmt.Printf("l_execDir6:%s\n", l_execDir)
	return l_execDir
}
func GetLogDir() string {
	return GetExeDir() + "/log"
}
