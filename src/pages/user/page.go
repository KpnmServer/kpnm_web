
package page_server

import (
	// http "net/http"

	iris "github.com/kataras/iris/v12"
	page_mnr "github.com/zyxgad/kpnm_svr/src/page_manager"
)

func IndexPage(ctx iris.Context){
	ctx.View("index.html")
}

func UserIndexPage(ctx iris.Context){
	ctx.View("user.html", iris.Map{
		
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

