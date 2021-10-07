
package page_chat

import (
	iris "github.com/kataras/iris/v12"

	kuser "github.com/KpnmServer/kpnm_web/src/user"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

func IndexPage(ctx iris.Context){
	user := kuser.GetCtxLog(ctx)
	ctx.View("/chat/index.html", iris.Map{
		"id": user.Id.String(),
	})
}

func init(){page_mnr.Register("/chat", func(group iris.Party){
	group.Use(kuser.MustLogHandler)
	InitApi(group)
	group.Get("/", IndexPage).ExcludeSitemap()
})}
