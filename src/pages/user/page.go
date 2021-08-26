
package page_user

import (
	iris "github.com/kataras/iris/v12"
	ufile "github.com/KpnmServer/go-util/file"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

var USER_DATA_PATH string = ufile.JoinPath(page_mnr.DATA_PATH, "user")

func IndexPage(ctx iris.Context){
	ctx.View("index.html")
}

func UserIndexPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	user := GetUserDataByName(name)
	ctx.View("user.html", iris.Map{
		"name": name,
		"desc": user.Desc,
		"frozen": user.Frozen,
	})
}

func LoginPage(ctx iris.Context){
	ctx.View("login.html")
}

func RegisterPage(ctx iris.Context){
	ctx.View("register.html")
}


func init(){page_mnr.Register("/user", "./webs/user", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/{name:string}", UserIndexPage)
	group.Get("/setting/login", LoginPage)
	group.Get("/setting/register", RegisterPage)
	InitApi(group)
})}

