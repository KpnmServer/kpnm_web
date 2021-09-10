
package page_zcs

import (
	os "os"
	ioutil "io/ioutil"
	time "time"
	errors "errors"
	strings "strings"

	json "github.com/KpnmServer/go-util/json"
	ufile "github.com/KpnmServer/go-util/file"
	iris "github.com/kataras/iris/v12"
	websocket "github.com/kataras/iris/v12/websocket"

	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

var _ZCS_CONN_TOKEN []byte = []byte{}

func verifyToken(tk string)(ok bool){
	btk := ([]byte)(tk)
	mlen := len(_ZCS_CONN_TOKEN)
	ok = true
	if len(btk) < mlen {
		mlen = len(btk)
	}
	for i := 0; i < mlen ;i++ {
		if _ZCS_CONN_TOKEN[i] != btk[i] {
			ok = false
		}else if ok {
			ok = true
		}
	}
	return
}

func GetUpLoadStatusWS()(iris.Handler){
	mapconn := make(map[string]string)

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message)(err error){
			conn := nsConn.Conn
			connid := conn.ID()
			svrname := mapconn[connid]
			info := ZCS_SVR_INFOS[svrname]
			var msgobj = make(json.JsonObj)
			err = json.DecodeJson(msg.Body, &msgobj)
			if err != nil { return }
			status := msgobj.GetString("status")
			switch status {
			case "ping":
				conn.Socket().WriteText(json.EncodeJson(json.JsonObj{"status": "pong"}), time.Second * 3)
			case "status":
				info.status = serverStatus(msgobj.GetUInt("code"))
			case "info":
				if msgobj.Has("interval") {
					info.interval = msgobj.GetUInt("interval")
				}
				if msgobj.Has("ticks") {
					info.ticks = msgobj.GetInt("ticks")
				}
				if msgobj.Has("cpu_num") {
					info.cpu_num = msgobj.GetUInt("cpu_num")
				}
				if msgobj.Has("java_version") {
					info.java_version = msgobj.GetString("java_version")
				}
				if msgobj.Has("os") {
					info.os = msgobj.GetString("os")
				}
				if msgobj.Has("max_mem") {
					info.max_mem = msgobj.GetUInt64("max_mem")
				}
				if msgobj.Has("total_mem") {
					info.total_mem = msgobj.GetUInt64("total_mem")
				}
				if msgobj.Has("used_mem") {
					info.used_mem = msgobj.GetUInt64("used_mem")
				}
				if msgobj.Has("cpu_load") {
					info.cpu_load = msgobj.GetFloat64("cpu_load")
				}
				if msgobj.Has("cpu_time") {
					info.cpu_time = msgobj.GetFloat64("cpu_time")
				}
			case "close_with_err":
				info.errstr = msgobj.GetString("error")
				fallthrough
			case "close":
				info.status = SERVER_STOPPED
				conn.Close()
			}
			return nil
		},
	})

	ws.OnConnect = func(conn *websocket.Conn)(error){
		connid := conn.ID()
		svrname := conn.Socket().Request().URL.Query().Get("svr")
		if svrname == "" { svrname = "main" }
		info, ok := ZCS_SVR_INFOS[svrname]
		if !ok || info.id != "" {
			return errors.New("Server's monitor is already exists")
		}
		info.id = connid
		info.errstr = ""
		mapconn[connid] = svrname
		return nil
	}

	ws.OnDisconnect = func(conn *websocket.Conn){
		svrname := mapconn[conn.ID()]
		delete(mapconn, conn.ID())
		ZCS_SVR_INFOS[svrname].id = ""
	}

	wshandler := websocket.Handler(ws)
	return func(ctx iris.Context){
		tk := ctx.URLParam("token")
		if !verifyToken(tk) {
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}
		wshandler(ctx)
	}
}

func UpLoadStatusApi(ctx iris.Context){
	tk := ctx.URLParam("token")
	if !verifyToken(tk) {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
	svrname := ctx.URLParamDefault("svr", "main")
	info, ok := ZCS_SVR_INFOS[svrname]
	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	var msgobj = make(json.JsonObj)
	bodyrd := ctx.Request().Body
	defer bodyrd.Close()
	err := json.ReadJson(bodyrd, &msgobj)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		ctx.Application().Logger().Debugf("Parse status json err: %v", err)
		return
	}
	info.id = msgobj.GetString("id")
	status := msgobj.GetString("status")
	switch status {
	case "status":
		info.status = serverStatus(msgobj.GetUInt("code"))
	case "info":
		if msgobj.Has("interval") {
			info.interval = msgobj.GetUInt("interval")
		}
		if msgobj.Has("ticks") {
			info.ticks = msgobj.GetInt("ticks")
		}
		if msgobj.Has("cpu_num") {
			info.cpu_num = msgobj.GetUInt("cpu_num")
		}
		if msgobj.Has("java_version") {
			info.java_version = msgobj.GetString("java_version")
		}
		if msgobj.Has("os") {
			info.os = msgobj.GetString("os")
		}
		if msgobj.Has("max_mem") {
			info.max_mem = msgobj.GetUInt64("max_mem")
		}
		if msgobj.Has("total_mem") {
			info.total_mem = msgobj.GetUInt64("total_mem")
		}
		if msgobj.Has("used_mem") {
			info.used_mem = msgobj.GetUInt64("used_mem")
		}
		if msgobj.Has("cpu_load") {
			info.cpu_load = msgobj.GetFloat64("cpu_load")
		}
		if msgobj.Has("cpu_time") {
			info.cpu_time = msgobj.GetFloat64("cpu_time")
		}
	case "close_with_err":
		info.errstr = msgobj.GetString("error")
		fallthrough
	case "close":
		// Do nothing
	}
	ctx.JSON(iris.Map{
		"id": page_mnr.SVR_UUID.String(),
	})
}

func SvrInfoApi(ctx iris.Context){
	svrname := ctx.URLParamDefault("svr", "main")
	svrinfo, ok := ZCS_SVR_INFOS[svrname]
	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.JSON(iris.Map{
		"id": svrinfo.id,
		"interval": svrinfo.interval,
		"cpu_num": svrinfo.cpu_num,
		"java_version": svrinfo.java_version,
		"os": svrinfo.os,
		"max_mem": svrinfo.max_mem,
	})
}

func SvrStatusApi(ctx iris.Context){
	svrname := ctx.URLParamDefault("svr", "main")
	svrinfo, ok := ZCS_SVR_INFOS[svrname]
	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.JSON(iris.Map{
		"id": svrinfo.id,
		"status": svrinfo.status.String(),
		"ticks": svrinfo.ticks,
		"total_mem": svrinfo.total_mem,
		"used_mem": svrinfo.used_mem,
		"cpu_load": svrinfo.cpu_load,
		"cpu_time": svrinfo.cpu_time,
		"errstr": svrinfo.errstr,
	})
}

func InitApi(group iris.Party){
	apigp := group.Party("/api")
	apigp.Get("/uploadstatusws", GetUpLoadStatusWS()).ExcludeSitemap()
	apigp.Post("/uploadstatus", page_mnr.SkipLogHandle, UpLoadStatusApi).ExcludeSitemap()
	apigp.Get("/svrinfo", page_mnr.SkipLogHandle, SvrInfoApi).ExcludeSitemap()
	apigp.Get("/svrstatus", page_mnr.SkipLogHandle, SvrStatusApi).ExcludeSitemap()
}

func init(){
	{ // read token
		var fd *os.File
		var err error
		fd, err = os.Open(ufile.JoinPath("data", "zcs", "token.txt"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		var data []byte
		data, err = ioutil.ReadAll(fd)
		if err != nil {
			panic(err)
		}

		_ZCS_CONN_TOKEN = ([]byte)(strings.TrimSpace((string)(data)))
	}
}
