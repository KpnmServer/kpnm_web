
package page_user

import (
	// http "net/http"

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

func SetLoginPage(ctx iris.Context){
	ctx.View("index.html")
}

func SetRegistePage(ctx iris.Context){
	ctx.View("index.html")
}


func init(){page_mnr.Register("/user", "./webs/user", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/{user:string}", UserIndexPage)
	group.Get("/setting/login", SetLoginPage)
	group.Get("/setting/registe", SetRegistePage)
})}

