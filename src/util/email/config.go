
package kpnmmail

import (
	os "os"
)

func init(){
	{ // load mail template files
		var templateFiles []string = make([]string, 0)
		basePath := util.GetAbsPath(TEMPLATE_PATH)
		var findFunc func(path string)
		findFunc = func(path string){
			finfos, err := ioutil.ReadDir(path)
			if err != nil {
				panic(err)
			}
			for _, info := range finfos {
				fpath := util.JoinPath(path, info.Name())
				if info.IsDir() {
					findFunc(fpath)
				}else{
					templateFiles = append(templateFiles, fpath)
				}
			}
		}
		findFunc(basePath)

		if len(templateFiles) > 0 {
			LoadHtmlFiles(templateFiles...)
		}
	}
}
