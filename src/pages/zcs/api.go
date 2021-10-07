
package page_zcs

import (
	fmt "fmt"
	os "os"
	ioutil "io/ioutil"
	time "time"
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
		ok = false
	}
	for i := 0; i < mlen ;i++ {
		if _ZCS_CONN_TOKEN[i] != btk[i] {
			ok = false
		}
	}
	return
}

var svrstatus_getter_connects = make(map[string]map[string]*websocket.Conn)

func onUpdateServerStatus(svrname string, updatemap map[string]bool){
	getterlist, ok := svrstatus_getter_connects[svrname]
	if !ok { return }
	info, ok := ZCS_SVR_INFOS[svrname]
	if !ok { return }
	jobj := json.JsonObj{}
	for k, v := range updatemap {
		if v {
			switch k {
			case "status":
				jobj["status"] = info.status.String()
			case "interval":
				jobj["interval"] = info.interval
			case "ticks":
				jobj["ticks"] = info.ticks
			case "cpu_num":
				jobj["cpu_num"] = info.cpu_num
			case "java_version":
				jobj["java_version"] = info.java_version
			case "os":
				jobj["os"] = info.os
			case "max_mem":
				jobj["max_mem"] = info.max_mem
			case "total_mem":
				jobj["total_mem"] = info.total_mem
			case "used_mem":
				jobj["used_mem"] = info.used_mem
			case "cpu_load":
				jobj["cpu_load"] = info.cpu_load
			case "cpu_time":
				jobj["cpu_time"] = info.cpu_time
			case "errstr":
				jobj["errstr"] = info.errstr
			}
		}
	}
	message := json.EncodeJson(jobj)
	for _, getter := range getterlist {
		go getter.Socket().WriteText(message, time.Second * 5)
	}
}

func GetUpLoadStatusWS()(iris.Handler){
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message)(err error){
			defer page_mnr.RecoverToEmail("/zcs/api/uploadstatusws:OnNativeMessage", page_mnr.OPERATIONS_EMAILS...)
			conn := nsConn.Conn
			// connid := conn.ID()
			ctx := websocket.GetContext(conn)
			values := ctx.Values()
			fmt.Printf("values: (%p)%v\n", values, values)
			svrname := values.Get("kweb.server.monitorws.name").(string)
			info := values.Get("kweb.server.monitorws.info").(*ServerInfo)
			var msgobj json.JsonObj
			msgobj, err = json.DecodeJsonObj(msg.Body)
			if err != nil { return }
			status := msgobj.GetString("status")
			updatemap := map[string]bool{}
			switch status {
			case "ping":
				conn.Socket().WriteText(json.JsonObj{"status": "pong"}.Bytes(), time.Second * 3)
				return nil
			case "status":
				info.status = serverStatus(msgobj.GetUInt("code"))
				updatemap["status"] = true
			case "info":
				if msgobj.Has("interval") {
					updatemap["interval"] = true
					info.interval = msgobj.GetUInt("interval")
				}
				if msgobj.Has("ticks") {
					updatemap["ticks"] = true
					info.ticks = msgobj.GetInt("ticks")
				}
				if msgobj.Has("cpu_num") {
					updatemap["cpu_num"] = true
					info.cpu_num = msgobj.GetUInt("cpu_num")
				}
				if msgobj.Has("java_version") {
					updatemap["java_version"] = true
					info.java_version = msgobj.GetString("java_version")
				}
				if msgobj.Has("os") {
					updatemap["os"] = true
					info.os = msgobj.GetString("os")
				}
				if msgobj.Has("max_mem") {
					updatemap["max_mem"] = true
					info.max_mem = msgobj.GetUInt64("max_mem")
				}
				if msgobj.Has("total_mem") {
					updatemap["total_mem"] = true
					info.total_mem = msgobj.GetUInt64("total_mem")
				}
				if msgobj.Has("used_mem") {
					updatemap["used_mem"] = true
					info.used_mem = msgobj.GetUInt64("used_mem")
				}
				if msgobj.Has("cpu_load") {
					updatemap["cpu_load"] = true
					info.cpu_load = msgobj.GetFloat64("cpu_load")
				}
				if msgobj.Has("cpu_time") {
					updatemap["cpu_time"] = true
					info.cpu_time = msgobj.GetFloat64("cpu_time")
				}
			case "close_with_err":
				updatemap["errstr"] = true
				info.errstr = msgobj.GetString("error")
				fallthrough
			case "close":
				info.status = SERVER_STOPPED
				conn.Close()
			}
			onUpdateServerStatus(svrname, updatemap)
			return nil
		},
	})

	ws.OnConnect = func(conn *websocket.Conn)(error){
		defer page_mnr.RecoverToEmail("/zcs/api/uploadstatusws:OnConnect", page_mnr.OPERATIONS_EMAILS...)
		connid := conn.ID()
		info := websocket.GetContext(conn).Values().Get("kweb.server.monitorws.info").(*ServerInfo)
		info.id = connid
		info.errstr = ""
		return nil
	}

	ws.OnDisconnect = func(conn *websocket.Conn){
		defer page_mnr.RecoverToEmail("/zcs/api/uploadstatusws:OnDisconnect", page_mnr.OPERATIONS_EMAILS...)
		info, ok := websocket.GetContext(conn).Values().Get("kweb.server.monitorws.info").(*ServerInfo)
		if !ok {
			name := websocket.GetContext(conn).Values().Get("kweb.server.monitorws.name").(string)
			info = ZCS_SVR_INFOS[name]
		}else{
			info.id = ""
			info.interval = 0
		}
	}

	wshandler := websocket.Handler(ws)
	return func(ctx iris.Context){
		svrname := ctx.URLParamDefault("svr", "main")
		info, ok := ZCS_SVR_INFOS[svrname]
		if !ok {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}
		tk := ctx.URLParam("token")
		if !verifyToken(tk) || info.id != "" {
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}
		ctx.Values().SetImmutable("kweb.server.monitorws.name", svrname)
		ctx.Values().Set("kweb.server.monitorws.info", info)
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
	updatemap := map[string]bool{}
	switch status {
	case "status":
		info.status = serverStatus(msgobj.GetUInt("code"))
		updatemap["status"] = true
	case "info":
		if msgobj.Has("interval") {
			updatemap["interval"] = true
			info.interval = msgobj.GetUInt("interval")
		}
		if msgobj.Has("ticks") {
			updatemap["ticks"] = true
			info.ticks = msgobj.GetInt("ticks")
		}
		if msgobj.Has("cpu_num") {
			updatemap["cpu_num"] = true
			info.cpu_num = msgobj.GetUInt("cpu_num")
		}
		if msgobj.Has("java_version") {
			updatemap["java_version"] = true
			info.java_version = msgobj.GetString("java_version")
		}
		if msgobj.Has("os") {
			updatemap["os"] = true
			info.os = msgobj.GetString("os")
		}
		if msgobj.Has("max_mem") {
			updatemap["max_mem"] = true
			info.max_mem = msgobj.GetUInt64("max_mem")
		}
		if msgobj.Has("total_mem") {
			updatemap["total_mem"] = true
			info.total_mem = msgobj.GetUInt64("total_mem")
		}
		if msgobj.Has("used_mem") {
			updatemap["used_mem"] = true
			info.used_mem = msgobj.GetUInt64("used_mem")
		}
		if msgobj.Has("cpu_load") {
			updatemap["cpu_load"] = true
			info.cpu_load = msgobj.GetFloat64("cpu_load")
		}
		if msgobj.Has("cpu_time") {
			updatemap["cpu_time"] = true
			info.cpu_time = msgobj.GetFloat64("cpu_time")
		}
	case "close_with_err":
		updatemap["errstr"] = true
		info.errstr = msgobj.GetString("error")
		fallthrough
	case "close":
		info.id = ""
		info.interval = 0
	}
	onUpdateServerStatus(svrname, updatemap)
	ctx.JSON(iris.Map{
		"id": page_mnr.SVR_UUID.String(),
	})
}

func GetSvrStatusWS()(iris.Handler){
	mapconn := make(map[string]string)

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message)(err error){
			defer page_mnr.RecoverToEmail("/zcs/api/svrstatusws:OnNativeMessage", page_mnr.OPERATIONS_EMAILS...)
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
				return nil
			case "init":
				conn.Socket().WriteText(json.EncodeJson(json.JsonObj{
					"status": info.status.String(),
					"interval": info.interval,
					"ticks": info.ticks,
					"cpu_num": info.cpu_num,
					"java_version": info.java_version,
					"os": info.os,
					"max_mem": info.max_mem,
					"total_mem": info.total_mem,
					"used_mem": info.used_mem,
					"cpu_load": info.cpu_load,
					"cpu_time": info.cpu_time,
					"errstr": info.errstr,
				}), time.Second * 5)
			}
			return nil
		},
	})

	ws.OnConnect = func(conn *websocket.Conn)(error){
		defer page_mnr.RecoverToEmail("/zcs/api/svrstatusws:OnConnect", page_mnr.OPERATIONS_EMAILS...)
		svrname := conn.Socket().Request().URL.Query().Get("svr")
		if svrname == "" { svrname = "main" }
		getterlist, ok := svrstatus_getter_connects[svrname]
		if !ok {
			getterlist = make(map[string]*websocket.Conn)
			svrstatus_getter_connects[svrname] = getterlist
		}
		getterlist[conn.ID()] = conn
		mapconn[conn.ID()] = svrname
		return nil
	}

	ws.OnDisconnect = func(conn *websocket.Conn){
		defer page_mnr.RecoverToEmail("/zcs/api/svrstatusws:OnDisconnect", page_mnr.OPERATIONS_EMAILS...)
		svrname := conn.Socket().Request().URL.Query().Get("svr")
		if svrname == "" { svrname = "main" }
		if getterlist, ok := svrstatus_getter_connects[svrname]; ok {
			delete(getterlist, conn.ID())
		}
		delete(mapconn, conn.ID())
	}

	wshandler := websocket.Handler(ws)
	return func(ctx iris.Context){
		if _, ok := ZCS_SVR_INFOS[ctx.URLParamDefault("svr", "main")]; !ok {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}
		wshandler(ctx)
	}
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
	apigp.Get("/svrstatusws", GetSvrStatusWS()).ExcludeSitemap()
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
