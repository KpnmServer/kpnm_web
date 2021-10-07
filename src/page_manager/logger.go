
package kweb_manager

import (
	os "os"
	time "time"
	iris "github.com/kataras/iris/v12"
	ufile "github.com/KpnmServer/go-util/file"
)

const (
	SKIP_LOGGER_CONTEXT_KEY = "kweb.logger.skip"
)

func SkipLog(ctx iris.Context){
	ctx.Values().Set(SKIP_LOGGER_CONTEXT_KEY, struct{}{})
}

func SkipLogHandle(ctx iris.Context){
	SkipLog(ctx)
	ctx.Next()
}

func isSkipLog(ctx iris.Context)(bool){
	return ctx.Values().Get(SKIP_LOGGER_CONTEXT_KEY) != nil
}

func BindLogger(app *iris.Application){
	var logFile *os.File
	ufile.CreateDir("./logs")

	app.UseRouter(func(ctx iris.Context){
		request := ctx.Request()
		var (
			ipaddr string
			method string
			code int
			startTime time.Time
			useTime time.Duration
			path string
			query string
		)
		startTime = time.Now()
		ctx.Next()
		useTime = time.Since(startTime)
		code = ctx.GetStatusCode()
		if code / 100 != 5 && isSkipLog(ctx) {
			return
		}
		ipaddr = request.RemoteAddr
		method = request.Method
		path = ctx.RequestPath(true)
		query = request.URL.RawQuery

		if rip, ok := request.Header["X-Real-Ip"]; ok && len(rip) > 0 {
			ipaddr = rip[0]
		}
		ctx.Application().Logger().Infof("[%s %s %d %v]:%s:%s", ipaddr, method, code, useTime, path, query)
	})

	checkEmptyLog := func(){
		if logFile != nil {
			logFile.Close()
			name := logFile.Name()
			logFile = nil
			logstat, err := os.Stat(name)
			if err == nil {
				logsize := logstat.Size()
				if logsize == 0 {
					app.Logger().Debugf("remove empty log \"%s\"", name)
					os.Remove(name)
				}
			}else{
				app.Logger().Debugf("Get log file stat err: %s", err.Error())
			}
		}
	}

	app.ConfigureHost(func(su *iris.Supervisor){
		su.RegisterOnShutdown(func(){
			if logFile != nil {
				logFile.Close()
			}
			checkEmptyLog()
		})
	})

	changeLogFileFunc := func(){
		logf, err := os.OpenFile("logs/" + time.Now().Format("20060102-15.log"),
			os.O_CREATE | os.O_WRONLY | os.O_APPEND | os.O_SYNC, os.ModePerm)
		if err != nil {
			app.Logger().Errorf("Create log file error: %s", err.Error())
			return
		}
		checkEmptyLog()
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
