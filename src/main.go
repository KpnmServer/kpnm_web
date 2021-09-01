
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
	defer func(){
		page_mnr.OnClose()
	}()

	page_mnr.RegisterHTML(app, "./webs/globals")
	page_mnr.NoSitemap(app.Favicon("./webs/static/favicon.ico"))
	page_mnr.ServeStatic(app, "/robots.txt", "./webs/robots.txt", false)
	page_mnr.ServeStatic(app, "/google9aa38deb43e89452.html", "./google9aa38deb43e89452.html", false)
	page_mnr.NoSitemap(app.HandleDir("/static", iris.Dir("./webs/static"))...)

	registerErrorPages(app)
	page_mnr.InitAll(app, func(group iris.Party){})

	page_mnr.BindSiteMap(app, "https://kpnm.waerba.com")

	ipaddr := fmt.Sprintf("%s:%d", "0.0.0.0", PORT)

	go func(){
		var runner iris.Runner
		if USE_HTTPS {
			runner = iris.TLS(ipaddr, CRT_FILE, KEY_FILE)
		}else{
			runner = iris.Addr(ipaddr)
		}
		app.Run(runner)
	}()

	bgcont := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		timeoutCtx, _ := context.WithTimeout(bgcont, 16 * time.Second)
		app.Logger().Warn("Closing server...")
		app.Shutdown(timeoutCtx)
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
	ufile.CreateDir("./logs")

	app.ConfigureHost(func(su *iris.Supervisor){
		su.RegisterOnShutdown(func(){
			if logFile != nil {
				logFile.Close()
			}
		})
	})
	app.UseRouter(func(ctx iris.Context){
		request := ctx.Request()
		var (
			ipaddr string = request.RemoteAddr
			method string = request.Method
			code int
			startTime time.Time
			useTime time.Duration
			path string = ctx.RequestPath(true)
			query string = request.URL.RawQuery
		)
		startTime = time.Now()
		ctx.Next()
		useTime = time.Since(startTime)
		code = ctx.GetStatusCode()

		if rip, ok := request.Header["X-Real-Ip"]; ok && len(rip) > 0 {
			ipaddr = rip[0]
		}
		ctx.Application().Logger().Infof("[%s %s %d %v]:%s:%s", ipaddr, method, code, useTime, path, query)
	})

	changeLogFileFunc := func(){
		logf, err := os.OpenFile("logs/" + time.Now().Format("20060102-15.log"),
			os.O_CREATE | os.O_WRONLY | os.O_APPEND | os.O_SYNC, os.ModePerm)
		if err != nil {
			app.Logger().Errorf("Create log file error: %s", err.Error())
			return
		}
		if logFile != nil {
			logFile.Close()
			logstat, err := logFile.Stat()
			if err == nil {
				logsize := logstat.Size()
				if logsize == 0 {
					os.Remove(logFile.Name())
				}
			}
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

