
package page_user

import (
	iris "github.com/kataras/iris/v12"
	uuid "github.com/google/uuid"
	ufile "github.com/KpnmServer/go-util/file"
	kses "github.com/KpnmServer/kpnm_web/src/session"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

var USER_DATA_PATH string = ufile.JoinPath(page_mnr.DATA_PATH, "user")

func IndexPage(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	userid, err := uuid.Parse(kses.GetSessionStr(uid, "loginuser"))
	if err != nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	user := GetUserData(userid)
	if user == nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	ctx.View("index.html", user)
}

func UserIndexPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	user := GetUserDataByName(name)
	if user == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.View("user.html", iris.Map{
		"name": name,
		"desc": user.Desc,
		"frozen": user.Frozen,
	})
}

func LoginPage(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	userid, err := uuid.Parse(kses.GetSessionStr(uid, "loginuser"))
	if err == nil && GetUserData(userid) != nil {
		ctx.Redirect("/user", iris.StatusFound)
		return
	}
	ctx.View("login.html")
}

func RegisterPage(ctx iris.Context){
	ctx.View("register.html")
}

func SettingPage(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	userid, err := uuid.Parse(kses.GetSessionStr(uid, "loginuser"))
	if err != nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	user := GetUserData(userid)
	if user == nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	ctx.View("setting.html", user)
}

func init(){page_mnr.Register("/user", "./webs/user", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/look/{name:string}", UserIndexPage)
	group.Get("/login", LoginPage)
	group.Get("/register", RegisterPage)
	group.Get("/setting", SettingPage).ExcludeSitemap()
	InitApi(group)
})}

