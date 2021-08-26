
package page_user

import (
	iris "github.com/kataras/iris/v12"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

func IndexPage(ctx iris.Context){
	ctx.View("index.html")
}

func UserIndexPage(ctx iris.Context){
	user := ctx.Params().Get("user")
	ctx.View("user.html", iris.Map{
		"user": user
	})
}

func LoginPage(ctx iris.Context){
	ctx.View("index.html")
}

func RegistePage(ctx iris.Context){
	ctx.View("index.html")
}


func init(){page_mnr.Register("/user", "./webs/user", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/{user:string}", UserIndexPage)
	group.Get("/setting/login", LoginPage)
	group.Get("/setting/registe", RegistePage)
	InitApi(group)
})}

