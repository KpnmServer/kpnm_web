
package page_index

import (
	iris "github.com/kataras/iris/v12"
	kuser "github.com/KpnmServer/kpnm_web/src/user"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

func IndexPage(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	data := iris.Map{
		"islog": false,
	}
	if user != nil {
		data["islog"] = true
		data["user"] = iris.Map{
			"id": user.Id,
			"name": user.Username,
			"email": user.Email,
		}
	}
	ctx.View("/index/index.html", data)
}

func init(){page_mnr.Register("/", func(group iris.Party){
	group.Get("/", IndexPage)
})}
