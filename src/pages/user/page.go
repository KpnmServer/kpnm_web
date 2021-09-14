
package page_user

import (
	iris "github.com/kataras/iris/v12"
	uuid "github.com/google/uuid"
	ufile "github.com/KpnmServer/go-util/file"
	kses "github.com/KpnmServer/kpnm_web/src/session"
	kuser "github.com/KpnmServer/kpnm_web/src/user"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

var USER_DATA_PATH string = ufile.JoinPath(page_mnr.DATA_PATH, "user")

func IndexPage(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	if user == nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	ctx.View("/user/index.html", user)
}

func UserIndexPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	user := kuser.GetUserDataByName(name)
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
	if kuser.GetCtxLog(ctx) != nil {
		ctx.Redirect("/user", iris.StatusFound)
		return
	}
	ctx.View("/user/login.html")
}

func RegisterPage(ctx iris.Context){
	ctx.View("/user/register.html")
}

func SettingPage(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	userid, err := uuid.Parse(kses.GetSessionStr(uid, "loginuser"))
	if err != nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	user := kuser.GetUserData(userid)
	if user == nil {
		ctx.Redirect("/user/login", iris.StatusFound)
		return
	}
	ctx.View("/user/setting.html", user)
}

func init(){page_mnr.Register("/user", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/look/{name:string}", UserIndexPage)
	group.Get("/login", LoginPage)
	group.Get("/register", RegisterPage)
	group.Get("/setting", SettingPage).ExcludeSitemap()
	InitApi(group)
})}

