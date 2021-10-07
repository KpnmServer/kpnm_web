
package page_server

import (
	time "time"
	base64 "encoding/base64"

	// json "github.com/KpnmServer/go-util/json"
	uuid "github.com/google/uuid"

	iris "github.com/kataras/iris/v12"
	websocket "github.com/kataras/iris/v12/websocket"

	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
	svr_mnr "github.com/KpnmServer/kpnm_web/src/server_manager"
)

const (
	_CONNWS_SERVERDATA_KEY = "kweb.server.api.connws.serverdata"
	_CONNWS_SERVERCONN_KEY = "kweb.server.api.connws.serverconn"
)

func GetServerConnWs()(iris.Handler){
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message)(err error){
			defer page_mnr.RecoverToEmail("/server/api/connws:OnNativeMessage", page_mnr.OPERATIONS_EMAILS...)
			conn := nsConn.Conn
			svrconn := websocket.GetContext(conn).Values().Get(_CONNWS_SERVERCONN_KEY).(*svr_mnr.ServerConn)
			re := svrconn.OnMessage((string)(msg.Body))
			if re != nil {
				conn.Socket().WriteText(re.Bytes(), time.Second * 3)
			}
			return nil
		},
	})

	ws.OnConnect = func(conn *websocket.Conn)(error){
		defer page_mnr.RecoverToEmail("/server/api/connws:OnConnect", page_mnr.OPERATIONS_EMAILS...)
		ctx := websocket.GetContext(conn)
		ctx_values := ctx.Values()
		server := ctx_values.Get(_CONNWS_SERVERDATA_KEY).(*svr_mnr.ServerData)
		svrconn := svr_mnr.NewServerConn(server, conn)
		ctx_values.Set(_CONNWS_SERVERCONN_KEY, svrconn)
		svrconn.Init()
		return nil
	}

	ws.OnDisconnect = func(conn *websocket.Conn){
		defer page_mnr.RecoverToEmail("/server/api/connws:OnDisconnect", page_mnr.OPERATIONS_EMAILS...)
		svrconn := websocket.GetContext(conn).Values().Get(_CONNWS_SERVERCONN_KEY).(*svr_mnr.ServerConn)
		svrconn.Close()
	}

	wshandler := websocket.Handler(ws)
	return func(ctx iris.Context){
		groupid, err := uuid.Parse(ctx.URLParam("G"))
		if err != nil {
			ctx.StatusCode(iris.StatusNotFound)
		}
		server := svr_mnr.GetServerByGroup(groupid)
		if server == nil {
			ctx.StatusCode(iris.StatusNotFound)
		}
		tk := ctx.URLParam("T")
		btk, err := base64.URLEncoding.DecodeString(tk)
		if err != nil || !server.VerifyToken(btk) {
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}
		ctx.Values().Set(_CONNWS_SERVERDATA_KEY, server)
		wshandler(ctx)
	}
}

func ServerSendEventApi(ctx iris.Context){
	groupid, err := uuid.Parse(ctx.URLParam("G"))
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	server := svr_mnr.GetServerByGroup(groupid)
	if server == nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	tk := ctx.URLParam("T")
	if !server.VerifyToken(([]byte)(tk)) {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
	ctx.StatusCode(iris.StatusNotFound)
}

func SearchApi(ctx iris.Context){
	keyword := ctx.URLParam("key")
	svrs := svr_mnr.SearchServer(keyword)
	var data = make([]string, len(svrs))
	for i, s := range svrs {
		data[i] = s.Id
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": data,
	})
}

func GetInfoApi(ctx iris.Context){
	id := ctx.URLParam("id")
	svr := svr_mnr.GetServer(id)
	if svr == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": iris.Map{
			"id": svr.Id,
			"name": svr.Name,
			"version": svr.Version,
			"desc": svr.Description,
		},
	})
}

func InitApi(group iris.Party){
	apigp := group.Party("/api")
	apigp.Get("/connws", GetServerConnWs()).ExcludeSitemap()
	apigp.Post("/sendevent", ServerSendEventApi).ExcludeSitemap()
	apigp.Get("/search", SearchApi).ExcludeSitemap()
	apigp.Get("/info", GetInfoApi).ExcludeSitemap()
}

