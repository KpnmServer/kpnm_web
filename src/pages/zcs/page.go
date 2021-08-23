
package page_zcs

import (
	iris "github.com/kataras/iris/v12"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)


func IndexPage(ctx iris.Context){
	ctx.View("index.html")
}

func init(){page_mnr.Register("/zcs", "./webs/zcs", func(group iris.Party){
	group.Get("/", IndexPage)
})}

