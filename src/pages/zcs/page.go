
package page_zcs

import (
	iris "github.com/kataras/iris/v12"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)


func IndexPage(ctx iris.Context){
	ctx.View("index.html", iris.Map{
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

func init(){page_mnr.Register("/zcs", "./webs/zcs", func(group iris.Party){
	group.Get("/", IndexPage)
})}

