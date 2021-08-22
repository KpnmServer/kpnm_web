
package kweb_util_file

import (
	os "os"
)

func JoinPath(paths ...string)(allpath string){
	allpath = ""
	for _, p := range paths {
		if len(p) == 0 {
			continue
		}
		if p[0] == '/' {
			allpath = p
			continue
		}
		if len(allpath) != 0 && allpath[len(allpath) - 1] != '/' {
			allpath += "/"
		}
		allpath += p
	}
	return allpath
}

func JoinPathWithoutAbs(paths ...string)(allpath string){
	allpath = ""
	for _, p := range paths {
		if len(p) == 0 {
			continue
		}
		if len(allpath) != 0 && allpath[len(allpath) - 1] != '/' && p[0] != '/' {
			allpath += "/"
		}
		allpath += p
	}
	return allpath
}

func GetRunPath()(cwdPath string){
	var err error
	cwdPath, err = os.Getwd()
	if err != nil {
		panic(err)
		return "."
	}
	return cwdPath
}

func GetAbsPath(path string)(string){
	return JoinPath(GetRunPath(), path)
}
