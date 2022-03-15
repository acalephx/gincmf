package router

import (
	"gincmf/app/controller/web/home"
	"gincmf/app/middleware"
	cmf "github.com/gincmf/cmf/bootstrap"
)

//web路由初始化

func WebListenRouter() {
	cmf.Get("/",new(home.Index).Index, middleware.HomeBaseController)
	cmf.Get("/about",new(home.About).Index, middleware.HomeBaseController)
}
