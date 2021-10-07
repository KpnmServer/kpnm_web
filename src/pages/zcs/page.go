
package page_zcs

import (
	iris "github.com/kataras/iris/v12"

	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)


func IndexPage(ctx iris.Context){
	ctx.View("/zcs/index.html", iris.Map{
		"cycle_imgs": [][2]string{
			{"/static/images/zcs/index/cycle-1.png", "cycle1"},
			{"/static/images/zcs/index/cycle-2.png", "cycle2"},
			{"/static/images/zcs/index/cycle-3.png", "cycle3"},
			{"/static/images/zcs/index/cycle-4.png", "cycle4"},
			{"/static/images/zcs/index/cycle-5.png", "cycle5"},
			{"/static/images/zcs/index/cycle-6.png", "cycle6"},
		},
	})
}

func StatusPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	if _, ok := ZCS_SVR_INFOS[name]; !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.View("/zcs/status.html", iris.Map{
		"svrname": name,
	})
}

func init(){page_mnr.Register("/zcs", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/status", func(ctx iris.Context){ ctx.Redirect("./status/main") })
	group.Get("/status/{name:string}", StatusPage)
	InitApi(group)
})}

