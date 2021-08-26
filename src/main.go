
package main

import (
	os "os"
	time "time"
	fmt "fmt"

	iris "github.com/kataras/iris/v12"
	accesslog "github.com/kataras/iris/v12/middleware/accesslog"

	kfutil "github.com/KpnmServer/go-util/file"
	json "github.com/KpnmServer/go-util/json"

	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
	_ "github.com/KpnmServer/kpnm_web/src/pages/index"
	_ "github.com/KpnmServer/kpnm_web/src/pages/server"
	_ "github.com/KpnmServer/kpnm_web/src/pages/user"
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
	{ // read config file
		var fd *os.File
		var err error
		fd, err = os.Open(kfutil.JoinPath("config", "config.json"))
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
	}
	app = NewApp()
	page_mnr.APPLICATION = app
	page_mnr.LOGGER = app.Logger()

	page_mnr.GLOBAL_I18N_MAP.LoadLanguage("en-us", kfutil.JoinPath("language", "en-us", "lang.json"))
	page_mnr.GLOBAL_I18N_MAP.LoadLanguage("zh-cn", kfutil.JoinPath("language", "zh-cn", "lang.json"))
}

func main(){
	page_mnr.RegisterHTML(app, "./webs/globals")
	app.Favicon("./webs/static/favicon.ico")
	page_mnr.RegisterStatic(app, "/robots.txt", "./webs/robots.txt", false)
	page_mnr.RegisterStatic(app, "/sitemap.xml", "./webs/sitemap.xml", true)
	page_mnr.RegisterStatic(app, "/google9aa38deb43e89452.html", "./google9aa38deb43e89452.html", false)
	app.HandleDir("/static", iris.Dir("./webs/static"))

	registerErrorPages(app)
	page_mnr.InitAll(app, func(group iris.Party){})

	ipaddr := fmt.Sprintf("%s:%d", "0.0.0.0", PORT)
	if USE_HTTPS {
		app.Run(iris.TLS(ipaddr, CRT_FILE, KEY_FILE))
	}else{
		app.Run(iris.Addr(ipaddr))
	}
}

func registerErrorPages(group iris.Party){
	group.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context){
		url := ctx.Request().URL
		ctx.View("404.html", iris.Map{
			"path": url.Path,
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

	bindLogger(app)

	return
}

func bindLogger(app *iris.Application){
	var logFile *os.File
	kfutil.CreateDir("./logs")

	app.ConfigureHost(func(su *iris.Supervisor){
		su.RegisterOnShutdown(func(){
			if logFile != nil {
				logFile.Close()
			}
		})
	})
	app.UseRouter(accesslog.New(kfutil.HandleWriter(func(bts []byte)(int, error){
		app.Logger().Info(string(bts))
		return len(bts), nil
	})).SetFormatter(&accesslog.Template{
		Text: "[{{.IP}} {{.Method}} {{.Code}} {{.Latency}}]:{{.Path}}:|{{.RequestValuesLine}} |{{.Request}} |{{.Response}}",
	}).Handler)

	changeLogFileFunc := func(){
		logf, err := os.OpenFile("./logs/" + time.Now().Format("20060102-15.log"),
			os.O_CREATE | os.O_WRONLY | os.O_APPEND | os.O_SYNC, os.ModePerm)
		if err != nil {
			app.Logger().Errorf("Create log file error: %s", err.Error())
			return
		}
		if logFile != nil {
			logFile.Close()
		}
		logFile = logf
		app.Logger().Printer.SetOutput(os.Stdout, logFile)
		app.Logger().Debugf("Using \"%s\" to log requests", logFile.Name())
	}
	changeLogFileFunc()
	go func(){
		for {
			select{
			case <-time.After(time.Duration(60 - (time.Now().Unix() / 60) % 60) * time.Minute):
				changeLogFileFunc()
			}
		}
	}()
}

