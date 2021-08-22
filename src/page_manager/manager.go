
package kweb_manager

import (
	iris "github.com/kataras/iris/v12"
)

type pageInfo struct{
	path string
	temppath string
	initcall func(iris.Party)
}

var _PAGES = make([]*pageInfo, 0)
var DEBUG = true

func RegisterHTML(group iris.Party, path string){
	tmpl := iris.HTML(path, ".html")
	tmpl.Reload(DEBUG)
	group.RegisterView(tmpl)
}

func Register(path string, temppath string, call func(iris.Party)){
	_PAGES = append(_PAGES, &pageInfo{
		path: path,
		temppath: temppath,
		initcall: call,
	})
}

func InitAll(app *iris.Application, initGroup func(iris.Party)){
	for _, info := range _PAGES {
		group := app.Party(info.path)
		initGroup(group)
		if info.temppath != "" {
			RegisterHTML(group, info.temppath)
		}
		info.initcall(group)
	}
}


