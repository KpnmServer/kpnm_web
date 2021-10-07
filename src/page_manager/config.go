
package kweb_manager

import (
	iris "github.com/kataras/iris/v12"
	golog "github.com/kataras/golog"
)

var (
	APPLICATION *iris.Application
	LOGGER *golog.Logger
	DEBUG = true
)
