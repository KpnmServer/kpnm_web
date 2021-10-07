
package main

import (
	context "context"
	os "os"
	signal "os/signal"
	syscall "syscall"
	time "time"
	fmt "fmt"
	strings "strings"

	iris "github.com/kataras/iris/v12"

	ufile "github.com/KpnmServer/go-util/file"
	json "github.com/KpnmServer/go-util/json"

	kuser "github.com/KpnmServer/kpnm_web/src/user"
	kutil "github.com/KpnmServer/kpnm_web/src/util"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
	_ "github.com/KpnmServer/kpnm_web/src/pages/index"
	_ "github.com/KpnmServer/kpnm_web/src/pages/server"
	_ "github.com/KpnmServer/kpnm_web/src/pages/user"
	_ "github.com/KpnmServer/kpnm_web/src/pages/chat"
	_ "github.com/KpnmServer/kpnm_web/src/pages/zcs"
)

var (
	PORT uint16 = 0
	USE_HTTPS bool = false
	CRT_FILE string = ""
	KEY_FILE string = ""
)

var app *iris.Application

func init(){
	var langfiles []string
	{ // read config file
		var fd *os.File
		var err error
		fd, err = os.Open(ufile.JoinPath("config", "config.json"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		var obj = make(json.JsonObj)
		err = json.ReadJson(fd, &obj)
		if err != nil {
			panic(err)
		}
		page_mnr.DEBUG = obj.GetBool("debug")
		PORT = obj.GetUInt16("port")
		USE_HTTPS = obj.GetBool("https")
		if USE_HTTPS {
			CRT_FILE = obj.GetString("crt_file")
			KEY_FILE = obj.GetString("key_file")
		}
		langfiles = obj.GetStrings("languages")
	}
	app = NewApp()
	page_mnr.APPLICATION = app
	page_mnr.LOGGER = app.Logger()

	for _, item := range langfiles {
		i := strings.IndexByte(item, ':')
		if i == -1 {
			app.Logger().Errorf("Format error: '%s'", item)
			continue
		}
		page_mnr.GLOBAL_I18N_MAP.LoadLanguage(item[:i], item[i + 1:])
	}
}

func main(){
	defer kutil.OnClose()

	if !page_mnr.DEBUG {
		app.Use(page_mnr.RecoverToEmailHandler(page_mnr.OPERATIONS_EMAILS...))
	}
	{
		tmpl := iris.Django("./webs/", ".html")
		tmpl.Reload(page_mnr.DEBUG)
		page_mnr.RegisterI18N(app, tmpl)
		app.RegisterView(tmpl)
	}
	page_mnr.NoSitemap(app.Favicon("./webs/static/favicon.ico"))
	page_mnr.ServeStatic(app, "/robots.txt", "./webs/robots.txt", false)
	page_mnr.ServeStatic(app, "/google9aa38deb43e89452.html", "./google9aa38deb43e89452.html", false)
	page_mnr.NoSitemap(app.HandleDir("/static", iris.Dir("./webs/static"))...)

	registerErrorPages(app)
	page_mnr.InitAll(app, func(group iris.Party){})

	page_mnr.BindSiteMap(app, "https://kpnm.waerba.com")

	ipaddr := fmt.Sprintf("%s:%d", "0.0.0.0", PORT)

	exitch := make(chan struct{}, 1)

	var runner iris.Runner
	if USE_HTTPS {
		runner = iris.TLS(ipaddr, CRT_FILE, KEY_FILE)
	}else{
		runner = iris.Addr(ipaddr)
	}
	go func(){
		err := app.Run(runner)
		if err != nil {
			app.Logger().Errorf("Error: %v", err)
		}
		if !page_mnr.DEBUG {
			page_mnr.SendCloseEmail(page_mnr.OPERATIONS_EMAILS...)
			if err != nil {
				page_mnr.SendCloseErrEmail(err, page_mnr.OPERATIONS_EMAILS...)
			}
		}
		exitch <- struct{}{}
	}()

	bgcont := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		timeoutCtx, _ := context.WithTimeout(bgcont, 5 * time.Second)
		app.Logger().Warn("Closing server...")
		app.Shutdown(timeoutCtx)
		<-exitch
	case <-exitch:
	}
}

func registerErrorPages(group iris.Party){
	group.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context){
		url := ctx.Request().URL
		ctx.View("/errors/404.html", iris.Map{
			"path": url.String(),
		})
	})
	group.OnErrorCode(iris.StatusUnauthorized, func(ctx iris.Context){
		url := ctx.Request().URL
		ctx.View("/errors/401.html", iris.Map{
			"path": url.String(),
			"islogin": kuser.GetCtxLog(ctx) != nil,
		})
	})
	group.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context){
		url := ctx.Request().URL
		ctx.View("/errors/500.html", iris.Map{
			"path": url.String(),
		})
	})
}

func NewApp()(app *iris.Application){
	app = iris.New()
	if page_mnr.DEBUG {
		app.Logger().SetLevel("debug")
	}else{
		app.Logger().SetLevel("info")
	}
	app.Logger().Printer.SetSync(true)
	app.Logger().SetTimeFormat("2006-01-02 15:04:05.000:")

	page_mnr.BindLogger(app)

	return
}

