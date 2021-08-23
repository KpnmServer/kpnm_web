
package page_server

import (
	fmt "fmt"
	ioutil "io/ioutil"
	http "net/http"

	iris "github.com/kataras/iris/v12"
	mc_util "github.com/zyxgad/go-mc_util"
	page_mnr "github.com/zyxgad/kpnm_svr/src/page_manager"
)


func IndexPage(group iris.Party)(iris.Handler){
	return func(ctx iris.Context){
	var svrList []*ServerInfo
	if files, err := ioutil.ReadDir(SERVER_DATA_PATH); err == nil {
		svrList = make([]*ServerInfo, 0, len(files))
		for _, file := range files {
			svr, e := GetServerInfo(file.Name())
			if e == nil {
				svrList = append(svrList, svr)
			}
		}
	}else{
		svrList = make([]*ServerInfo, 0)
	}
	ctx.View("index.html", svrList)
}}

func ServerPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	svr, err := GetServerInfo(name)
	if err != nil {
		page_mnr.LOGGER.Debugf("Get server \"%s\" error: %v", name, err)
		ctx.StatusCode(http.StatusNotFound)
		return
	}
	ctx.View("info.html", iris.Map{
		"name": svr.Name,
		"version": svr.Version,
		"desc": svr.Description,
	})
}

func InfoMePage(ctx iris.Context){
	name := ctx.Params().Get("name")
	data, err := GetServerReadme(name)
	if err != nil {
		ctx.StatusCode(http.StatusNotFound)
		return
	}
	ctx.Markdown(data)
}

func StatusPagePost(ctx iris.Context){
	name := ctx.Params().Get("name")
	svr, err := GetServerInfo(name)
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"errorMessage": err.Error(),
		})
		return
	}
	var (
		status *mc_util.ServerStatus
		host string
		port uint16
	)
	for _, addr := range svr.Addrs {
		host = addr.GetString(0)
		port = addr.GetUInt16(1)
		page_mnr.LOGGER.Debugf("Pinging \"%s:%d\"", host, port)
		status, err = mc_util.Ping(host, port)
		if err == nil {
			page_mnr.LOGGER.Debugf("Ping \"%s:%d\" success", host, port)
			break
		}
		page_mnr.LOGGER.Debugf("Ping \"%s:%d\" failed: %v", host, port, err)
	}
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"ping": status.Delay,
		"desc": status.Description,
		"ip": fmt.Sprintf("%s:%d", host, port),
		"players": status.Players,
		"player_count": status.Online_player,
		"player_max_count": status.Max_player,
		"version": status.Version,
		"favicon": status.Favicon,
	})
}

func init(){page_mnr.Register("/server", "./webs/server", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/{name:string}", ServerPage)
	group.Get("/{name:string}/infome", InfoMePage)
	group.Get("/{name:string}/status", StatusPagePost)
})}

