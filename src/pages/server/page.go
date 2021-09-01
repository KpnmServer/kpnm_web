
package page_server

import (
	fmt "fmt"
	bytes "bytes"
	ioutil "io/ioutil"

	iris "github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	mc_util "github.com/KpnmServer/go-mc_util"
	ufile "github.com/KpnmServer/go-util/file"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

var SERVER_DATA_PATH string = ufile.JoinPath(page_mnr.DATA_PATH, "server")

func IndexPage(ctx iris.Context){
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
}

func ServerPage(ctx iris.Context){
	name := ctx.Params().Get("name")
	svr, err := GetServerInfo(name)
	if err != nil {
		ctx.Application().Logger().Debugf("Get server \"%s\" error: %v", name, err)
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	var (
		readme_data []byte
		readme_buf *bytes.Buffer = bytes.NewBuffer([]byte{})
	)
	readme_data, err = GetServerReadme(name)
	if err == nil {
		iris_context.WriteMarkdown(readme_buf, readme_data, iris_context.DefaultMarkdownOptions)
	}
	ctx.View("info.html", iris.Map{
		"name": svr.Name,
		"version": svr.Version,
		"desc": svr.Description,
		"readme": readme_buf.String(),
	})
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
		ctx.Application().Logger().Debugf("Pinging \"%s:%d\"", host, port)
		status, err = mc_util.Ping(host, port)
		if err == nil {
			ctx.Application().Logger().Debugf("Ping \"%s:%d\" success", host, port)
			break
		}
		ctx.Application().Logger().Debugf("Ping \"%s:%d\" failed: %v", host, port, err)
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
	group.Get("/{name:string}/status", StatusPagePost)
})}

