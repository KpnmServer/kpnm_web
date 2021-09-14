
package kweb_manager

import (
	os "os"
	ioutil "io/ioutil"
	htmltmpl "html/template"

	uuid "github.com/google/uuid"
	ufile "github.com/KpnmServer/go-util/file"
	iris "github.com/kataras/iris/v12"
	view "github.com/kataras/iris/v12/view"
	router "github.com/kataras/iris/v12/core/router"
)

var SVR_UUID uuid.UUID = uuid.New()

type pageInfo struct{
	path string
	temppath string
	initcall func(iris.Party)
}

var _PAGES = make([]*pageInfo, 0)
var _CLOSE_HANDLES = make([]func(), 0)
type tmplCache struct{
	tmpl *htmltmpl.Template
	body []byte
}
var _GLOBAL_TMPLS = make([]tmplCache, 0)

func GlobalHTML(path string){
	var tmpl *view.HTMLEngine = iris.HTML(path, ".html")
	tmpl.Load()
	for _, t := range tmpl.Templates.Templates() {
		name := t.Name()
		f := ufile.JoinPath(path, name)
		var err error
		fd, err := os.Open(f)
		if err != nil { continue }
		data, err := ioutil.ReadAll(fd)
		fd.Close()
		if err != nil { continue }
		_GLOBAL_TMPLS = append(_GLOBAL_TMPLS, tmplCache{tmpl: t, body:data})
	}
}

type htmlHandler func(group iris.Party, tmpl view.EngineFuncer)

func RegisterI18N(group iris.Party, tmpl view.EngineFuncer){
	var i18nmap I18nMap = GetGlobalI18nMapCopy()
	var i18nmapp *I18nMap = &i18nmap
	group.Use(LocalHandle(i18nmapp))
	tmpl.AddFunc("a", i18nmapp.Localization)
	tmpl.AddFunc("getlang", i18nmapp.GetLocalLang)
}

func RegisterHTML(group iris.Party, path string, handlers ...htmlHandler){
	var tmpl *view.HTMLEngine = iris.HTML(path, ".html")
	tmpl.Reload(DEBUG)
	tmpl.Load()
	for _, t := range _GLOBAL_TMPLS {
		tmpl.ParseTemplate(t.tmpl.Name(), t.body, htmltmpl.FuncMap{})
	}
	for _, h := range handlers {
		h(group, tmpl)
	}
	group.RegisterView(tmpl)
}

func Register(path string, call func(iris.Party)){
	_PAGES = append(_PAGES, &pageInfo{
		path: path,
		temppath: "",
		initcall: call,
	})
}

func RegisterTemp(path string, temppath string, call func(iris.Party)){
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
			RegisterHTML(group, info.temppath, RegisterI18N)
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

func NoSitemap(routes ...*router.Route)([]*router.Route){
	for _, r := range routes {
		r.ExcludeSitemap()
	}
	return routes
}

func RegisterClose(call func()){
	_CLOSE_HANDLES = append(_CLOSE_HANDLES, call)
}

func OnClose(){
	for _, h := range _CLOSE_HANDLES {
		h()
	}
}
