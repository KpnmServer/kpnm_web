
package kweb_manager

import (
	os "os"
	ioutil "io/ioutil"
	htmltmpl "html/template"

	iris "github.com/kataras/iris/v12"
	router "github.com/kataras/iris/v12/core/router"
)

type pageInfo struct{
	path string
	temppath string
	initcall func(iris.Party)
}

var _PAGES = make([]*pageInfo, 0)
var _CLOSE_HANDLES = make([]func(), 0)

func RegisterHTML(group iris.Party, path string){
	tmpl := iris.HTML(path, ".html")
	tmpl.Reload(DEBUG)
	var i18nmap I18nMap = GetGlobalI18nMapCopy()
	var i18nmapp *I18nMap = &i18nmap
	group.Use(LocalHandle(i18nmapp))
	tmpl.AddFunc("a", i18nmapp.Localization)
	tmpl.AddFunc("getlang", i18nmapp.GetLocalLang)
	tmpl.AddFunc("noesc", func(dt string)(interface{}){ return htmltmpl.HTML(dt) })
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

func ServeStatic(group iris.Party, route string, path string, reload bool){
	var router *router.Route
	if reload {
		router = group.Get(route, func(ctx iris.Context){
			var fd *os.File
			var err error
			fd, err = os.Open(path)
			if err != nil {
				group.Logger().Errorf("Register static file error: %v", err)
				return
			}
			data, err := ioutil.ReadAll(fd)
			if err != nil {
				group.Logger().Errorf("Register static file error: %v", err)
				return
			}
			ctx.Write(data)
		})
	}else{
		var fd *os.File
		var err error
		fd, err = os.Open(path)
		if err != nil {
				group.Logger().Errorf("Register static file error: %v", err)
				return
		}
		data, err := ioutil.ReadAll(fd)
		if err != nil {
				group.Logger().Errorf("Register static file error: %v", err)
				return
		}
		router = group.Get(route, func(ctx iris.Context){
			ctx.Write(data)
		})
	}
	router.ExcludeSitemap()
}

func NoSitemap(routes ...*router.Route){
	for _, r := range routes {
		r.ExcludeSitemap()
	}
}

func RegisterClose(call func()){
	_CLOSE_HANDLES = append(_CLOSE_HANDLES, call)
}

func OnClose(){
	for _, h := range _CLOSE_HANDLES {
		h()
	}
}
