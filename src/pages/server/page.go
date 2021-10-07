
package page_server

import (
	fmt "fmt"
	bytes "bytes"
	strings "strings"
	strconv "strconv"

	iris "github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	mc_util "github.com/KpnmServer/go-mc_util"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
	svr_mnr "github.com/KpnmServer/kpnm_web/src/server_manager"
)

func IndexPage(ctx iris.Context){
	ctx.View("/server/index.html")
}

func ServerPage(ctx iris.Context){
	id := ctx.Params().Get("id")
	svr := svr_mnr.GetServer(id)
	if svr == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	var (
		err error
		readme_data []byte
		readme_buf *bytes.Buffer = bytes.NewBuffer([]byte{})
	)
	readme_file := svr.GetDataFolder().File("README.MD")
	if readme_file.IsExist() {
		readme_data, err = readme_file.ReadAll()
		if err == nil {
			iris_context.WriteMarkdown(readme_buf, readme_data, iris_context.DefaultMarkdownOptions)
		}
	}
	ctx.View("/server/info.html", iris.Map{
		"id": id,
		"name": svr.Name,
		"version": svr.Version,
		"desc": svr.Description,
		"readme": readme_buf.String(),
	})
}

func StatusPagePost(ctx iris.Context){
	id := ctx.Params().Get("id")
	svr := svr_mnr.GetServer(id)
	if svr == nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "NoServerFoundError",
		})
		return
	}
	var (
		status *mc_util.ServerStatus
		host string
		port uint16
		err error
	)
	svraddrs := strings.Split(svr.Addrstr, ";")
	for _, addr := range svraddrs {
		host, port = splitAddr(addr)
		ctx.Application().Logger().Debugf("Pinging \"%s\"", host, port)
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

func splitAddr(addr string)(host string, port uint16){
	b := strings.SplitN(addr, ":", 2)
	if len(b) < 2 {
		return
	}
	host = b[0]
	p0, err := strconv.Atoi(b[1])
	if err != nil {
		return
	}
	port = (uint16)(p0)
	return
}

func init(){page_mnr.Register("/server", func(group iris.Party){
	group.Get("/", IndexPage)
	group.Get("/{id:string}", ServerPage)
	group.Get("/{id:string}/status", page_mnr.SkipLogHandle, StatusPagePost)
	InitApi(group)
})}

