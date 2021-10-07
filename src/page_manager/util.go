
package kweb_manager

import (
	fmt "fmt"
	time "time"

	kutil "github.com/KpnmServer/kpnm_web/src/util"
	kmail "github.com/KpnmServer/kpnm_web/src/email"

	iris "github.com/kataras/iris/v12"
)

var OPERATIONS_EMAILS = []string{"kupond@outlook.com"}

func SendCloseEmail(to ...string)(error){
	return kmail.SendHtml(to, "Server closed", "server/close.html", time.Now().Format("2006-01-02 15:04:05 -0700"))
}

func SendCloseErrEmail(err error, to ...string)(error){
	return kmail.SendHtml(to, "Server closed with error", "server/close_err.html", iris.Map{
		"time": time.Now().Format("2006-01-02 15:04:05 -0700"),
		"error": err.Error(),
	})
}

func RecoverToEmailHandler(to ...string)(iris.Handler){
	return func(ctx iris.Context){
		defer func(){
			if err := recover(); err != nil {
				path := ctx.Request().URL.String()
				kmail.SendHtml(to, "Server panic", "server/panic.html", iris.Map{
					"path": path,
					"error": fmt.Sprint(err),
					"stacks": kutil.GetStacks(),
				})
				ctx.Application().Logger().Errorf("Error in '%s': %v", path, err)

				ctx.StatusCode(iris.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}

func RecoverToEmail(desc string, to ...string){
	if err := recover(); err != nil {
		kmail.SendHtml(to, "Server panic", "server/panic.html", iris.Map{
			"path": "desc: " + desc,
			"error": fmt.Sprint(err),
			"stacks": kutil.GetStacks(),
		})
		LOGGER.Errorf("Error in '%s': %v", desc, err)
	}
}
