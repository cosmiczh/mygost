package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ginuerzh/gost"
)

func GetExeName() string {
	var l_basename string
	if idx := strings.LastIndexAny(os.Args[0], "/\\"); idx < 0 {
		l_basename = os.Args[0]
	} else {
		l_basename = os.Args[0][idx+1:]
	}
	return l_basename
}

func GetExeBaseName() string {
	var l_basename = GetExeName()

	if len(l_basename) > 4 && strings.ToLower(l_basename[len(l_basename)-4:]) == ".exe" {
		l_basename = l_basename[:len(l_basename)-4]
	}
	return l_basename
}

func GetExePath() string { return gost.GetExeDir() + "/" + GetExeName() }

func SearchFile(plist_file *[]os.FileInfo, dirname string, name_pattern string) ([]os.FileInfo, error) {
	return search_in_fs(plist_file, dirname, name_pattern, false)
}
func SearchDir(plist_file *[]os.FileInfo, dirname string, name_pattern string) ([]os.FileInfo, error) {
	return search_in_fs(plist_file, dirname, name_pattern, true)
}

func search_in_fs(plist_file *[]os.FileInfo, dirname string, name_pattern string, dirORfile bool) ([]os.FileInfo, error) {
	if len(name_pattern) < 1 {
		return nil, fmt.Errorf("The file name to search is null.")
	}
	if plist_file == nil {
		var l_tmp []os.FileInfo = nil
		plist_file = &l_tmp
	}
	if *plist_file == nil {
		finfo, err := ioutil.ReadDir(dirname)
		if err != nil {
			return nil, err
		}
		*plist_file = finfo
	}
	l_retlist := make([]os.FileInfo, 0)
	l_mode_arr := SplitFields(name_pattern, "*", false)
	if name_pattern[0] == '*' && len(l_mode_arr) > 0 && l_mode_arr[0] != "" {
		l_mode_arr = append([]string{""}, l_mode_arr...)
	}
	if name_pattern[len(name_pattern)-1] == '*' && len(l_mode_arr) > 0 && l_mode_arr[len(l_mode_arr)-1] != "" {
		l_mode_arr = append(l_mode_arr, "")
	}
	l_mode_arr1 := make([]int8, len(l_mode_arr))
	l_mode_arr2 := make([][]string, len(l_mode_arr))
	for i, wildstr := range l_mode_arr {
		if len(wildstr) < 1 {
			l_mode_arr1[i] = -1
			continue
		}
		l_mode_arr1[i] = 1 //l_incl := true
		if wildstr[0] == '|' {
			if len(wildstr) < 2 {
				continue
			}
			wildstr = wildstr[1:]
			l_mode_arr1[i] = 0 //l_incl = false
		}
		l_mode_arr2[i] = SplitFields(wildstr, "/", false)
	}
	for _, finfo := range *plist_file {
		if dirORfile != finfo.IsDir() {
			continue
		}
		l_xname, l_satisfy := finfo.Name(), true

		for j, incl := range l_mode_arr1 {
			if incl < 0 {
				continue
			}
			l_incl, l_wilds := incl > 0, l_mode_arr2[j]

			if len(l_mode_arr1) == 1 {
				l_ismatch := IndexFunc(len(l_wilds), func(i int) bool { return l_xname == l_wilds[i] }) >= 0
				l_satisfy = (l_incl && l_ismatch) || (!l_incl && !l_ismatch)
			} else if j == 0 {
				l_wildi := IndexFunc(len(l_wilds),
					func(i int) bool {
						if len(l_xname) >= len(l_wilds[i]) {
							return l_xname[:len(l_wilds[i])] == l_wilds[i]
						}
						return false
					})
				l_ismatch := l_wildi >= 0

				l_satisfy = (l_incl && l_ismatch) || (!l_incl && !l_ismatch)
				if !l_satisfy {
					break
				}
				if l_incl {
					l_xname = l_xname[len(l_wilds[l_wildi]):]
				}
			} else if j == len(l_mode_arr1)-1 {
				l_wildi := IndexFunc(len(l_wilds),
					func(i int) bool {
						if len(l_xname) >= len(l_wilds[i]) {
							return l_xname[len(l_xname)-len(l_wilds[i]):] == l_wilds[i]
						}
						return false
					})
				l_ismatch := l_wildi >= 0
				l_satisfy = (l_incl && l_ismatch) || (!l_incl && !l_ismatch)
			} else {
				l_wildi := IndexFunc(len(l_wilds), func(i int) bool { return strings.Index(l_xname, l_wilds[i]) >= 0 })
				l_ismatch := l_wildi >= 0

				l_satisfy = (l_incl && l_ismatch) || (!l_incl && !l_ismatch)
				if !l_satisfy {
					break
				}
				if l_incl {
					l_pos := strings.Index(l_xname, l_wilds[l_wildi])
					l_xname = l_xname[l_pos+len(l_wilds[l_wildi]):]
				}
			}
		}
		if l_satisfy {
			l_retlist = append(l_retlist, finfo)
		}
	}
	return l_retlist, nil
}
