
package page_user

import (
	iris "github.com/kataras/iris/v12"
	kuser "github.com/KpnmServer/kpnm_web/src/user"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
	data_mnr "github.com/KpnmServer/kpnm_web/src/data_manager"
)

var USER_DATA_FOLDER = data_mnr.GetDataFolder("user")

func IndexPage(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	ctx.View("/user/index.html", user)
}

func UserLookPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	user := kuser.GetUserDataByName(name)
	if user == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.View("/user/look.html", iris.Map{
		"name": name,
		"id": user.Id.String(),
		"desc": user.Desc,
		"email": user.Email,
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
	user := kuser.GetCtxLog(ctx)
	ctx.View("/user/setting.html", user)
}

func init(){page_mnr.Register("/user", func(group iris.Party){
	group.Get("/", kuser.LogOrRedirectHandler("/user/login"), IndexPage)
	group.Get("/@/{name:string}", UserLookPage)
	group.Get("/login", LoginPage)
	group.Get("/register", RegisterPage)
	group.Get("/setting", kuser.LogOrRedirectHandler("/user/login"), SettingPage).ExcludeSitemap()
	InitApi(group)
})}

