
package page_chat

import (
	time "time"
	regexp  "regexp"

	json "github.com/KpnmServer/go-util/json"
	uuid "github.com/google/uuid"

	iris "github.com/kataras/iris/v12"
	websocket "github.com/kataras/iris/v12/websocket"

	kuser "github.com/KpnmServer/kpnm_web/src/user"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
	chater "github.com/KpnmServer/kpnm_web/src/chater"
)

var (
	reg_group_id *regexp.Regexp = regexp.MustCompile(`^[A-Za-z_-][0-9A-Za-z_-]{2,31}$`)
)

func GetConnectWS()(iris.Handler){
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message)(err error){
			defer page_mnr.RecoverToEmail("/chat/api/connectws:OnNativeMessage", page_mnr.OPERATIONS_EMAILS...)
			conn := nsConn.Conn
			var msgobj = make(json.JsonObj)
			err = json.DecodeJson(msg.Body, &msgobj)
			if err != nil { return }
			status := msgobj.GetString("status")
			switch status {
			case "ping":
				conn.Socket().WriteText(json.EncodeJson(json.JsonObj{"status": "pong"}), time.Second * 3)
				return nil
			case "close":
				conn.Close()
			}
			return nil
		},
	})

	ws.OnConnect = func(conn *websocket.Conn)(error){
		defer page_mnr.RecoverToEmail("/zcs/api/uploadstatusws:OnConnect", page_mnr.OPERATIONS_EMAILS...)
		return nil
	}

	ws.OnDisconnect = func(conn *websocket.Conn){
		defer page_mnr.RecoverToEmail("/zcs/api/uploadstatusws:OnDisconnect", page_mnr.OPERATIONS_EMAILS...)
	}

	wshandler := websocket.Handler(ws)
	return func(ctx iris.Context){
		// user := kuser.GetCtxLog(ctx)
		wshandler(ctx)
	}
}

func UnreadMsgApi(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	groupid := uuid.MustParse(ctx.Params().Get("group"))
	group := chater.GetGroup(groupid)
	if group == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	mem := group.GetMember(user.Id)
	if mem == nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
	msgs, err := mem.GetUnreadMsgs()
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "GetUnreadMsgsError",
			"errorMessage": err.Error(),
		})
		return
	}
	var data = make([]json.JsonObj, len(msgs))
	for i, m := range msgs {
		data[i] = m.Json()
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": data,
	})
}

func SendMsgApi(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	groupid := uuid.MustParse(ctx.Params().Get("group"))
	group := chater.GetGroup(groupid)
	if group == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	mem := group.GetMember(user.Id)
	if mem == nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
	reader := ctx.Request().Body
	defer reader.Close()
	msgobj, err := json.ReadJsonObj(reader)
	reader.Close()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	var msgidlist []string
	var msg chater.Message = nil
	if msgobj.Has("list") {
		mlist := msgobj.GetObjs("list")
		msgidlist = make([]string, len(mlist))
		for i, m := range mlist {
			mtype := (chater.MsgType)(m.GetInt8("type"))
			tg := mem.NewMessage(mtype, m.GetString("data"))
			if i == 0 {
				tg = msg
			}
			msg.Append(tg)
			msgidlist[i] = tg.MData().Id.String()
		}
	}else{
		msg = mem.NewMessage((chater.MsgType)(msgobj.GetInt8("type")), msgobj.GetString("data"))
	}
	_, err = msg.Now().InsertAll()
	ctx.JSON(iris.Map{
		"status": "ok",
		"msgs": msgidlist,
	})
}

func PutFileApi(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	groupid := uuid.MustParse(ctx.Params().Get("group"))
	group := chater.GetGroup(groupid)
	if group == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	fileid := uuid.MustParse(ctx.Params().Get("fileid"))
	filemsg := group.GetFileMsg(fileid)
	if filemsg == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	if filemsg.MData().Owner != user.Id {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
	if filemsg.HasData() {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "FileExists",
		})
		return
	}
	reader := ctx.Request().Body
	defer reader.Close()
	_, err := filemsg.ReadFrom(reader)
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "ReadError",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{"status": "ok"})
}

func GetFileApi(ctx iris.Context){
	groupid := uuid.MustParse(ctx.Params().Get("group"))
	group := chater.GetGroup(groupid)
	if group == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	fileid := uuid.MustParse(ctx.Params().Get("fileid"))
	filemsg := group.GetFileMsg(fileid)
	if filemsg == nil || !filemsg.HasData() {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	writer := ctx.ResponseWriter()
	if ctx.URLParam("nodownload") != "T" {
		writer.Header().Set("Content-Disposition", "attachment;filename=" + filemsg.GetData())
	}
	_, err := filemsg.WriteTo(writer)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

func CreateGroupApi(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	if !user.IsGood() {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
	gps := chater.GetGroupsByOwner(user.Id)
	if gps == nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if !user.IsTrusted(){
		if (user.Level < kuser.LEVEL_WELL && len(gps) >= 3) ||
		   (len(gps) >= 10) {
			ctx.JSON(iris.Map{
				"status": "error",
				"error": "ToMuchGroupError",
			})
			return
		}
	}
	name := ctx.PostValueTrim("name")
	if !reg_group_id.MatchString(name) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Group name is illegal data.",
		})
		return
	}
	desc := ctx.PostValueTrim("desc")
	group := chater.NewGroup(name, user.Id, chater.GROUP_NORMAL, desc)
	err := group.Insert()
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "InsertGroupError",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"group": group.Id.String(),
	})
}

func InitApi(group iris.Party){
	apigp := group.Party("/api")
	apigp.Get("/connectws", kuser.MustLogHandler, GetConnectWS()).ExcludeSitemap()
	apigp.Get("/{group:uuid}/unreadmsg", kuser.MustLogHandler, UnreadMsgApi).ExcludeSitemap()
	apigp.Post("/{group:uuid}/sendmsg", kuser.MustLogHandler, SendMsgApi).ExcludeSitemap()
	apigp.Put("/{group:uuid}/file/{fileid:uuid}", kuser.MustLogHandler, PutFileApi).ExcludeSitemap()
	apigp.Get("/{group:uuid}/file/{fileid:uuid}", GetFileApi).ExcludeSitemap()
	apigp.Post("/creategroup", kuser.MustLogHandler, CreateGroupApi).ExcludeSitemap()
}
