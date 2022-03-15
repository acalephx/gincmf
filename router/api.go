package router

import (
	"gincmf/app/controller/api/admin"
	"gincmf/app/controller/api/app"
	"gincmf/app/controller/common"
	"gincmf/app/middleware"
	"gincmf/app/migrate"
	"gincmf/app/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"

	cmf "github.com/gincmf/cmf/bootstrap"
)

//web路由初始化
func ApiListenRouter() {

	// 全局中间件
	cmf.HandleFunc = append(cmf.HandleFunc, middleware.AllowCors)

	adminGroup := cmf.Group("api/admin", middleware.ValidationBearerToken, middleware.ValidationAdmin, middleware.ApiBaseController, middleware.Rbac)
	{
		adminGroup.Rest("/settings", new(admin.Settings))
		adminGroup.Rest("/assets", new(admin.Assets))
		adminGroup.Rest("/upload", new(admin.Upload))
		adminGroup.Rest("/role", new(admin.Role))
		adminGroup.Rest("/user", new(admin.User))
		adminGroup.Get("/admin_menu", new(admin.Menu).Get)
		adminGroup.Get("/authorize", new(admin.Authorize).Get)
		adminGroup.Get("/authorize/:id", new(admin.Authorize).Show)
		adminGroup.Get("/auth_access/:id", new(admin.AuthAccess).Show)
		adminGroup.Post("/auth_access/:id", new(admin.AuthAccess).Edit)
		adminGroup.Post("/auth_access", new(admin.AuthAccess).Store)

		adminGroup.Rest("/nav", new(admin.Nav))
		adminGroup.Rest("/nav_item", new(admin.NavItem))
		adminGroup.Get("/nav_item_options", new(admin.NavItem).OptionsList)
		adminGroup.Get("/nav_item_urls", new(admin.NavItem).OptionsUrls)

		adminGroup.Get("/theme/init", new(admin.Theme).Init)
		adminGroup.Post("/theme_file/:id", new(admin.ThemeFile).Save)
	}

	appGroup := cmf.Group("api/app/v1", middleware.ApiBaseController)
	{
		appGroup.Get("/settings", new(app.Settings).Get)
		appGroup.Get("/route", new(app.Route).List)
		appGroup.Get("/nav_item", new(app.NavItem).Get)
		appGroup.Get("/theme_file", new(app.ThemeFile).Detail)
		appGroup.Get("/theme_file/list", new(app.ThemeFile).List)

	}

	// 清除缓存
	cmf.Get("/api/clear", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.JSON(200, model.ReturnData{
			Code: 1,
			Data: nil,
			Msg:  "清除成功！",
		})
	})

	// 获取当前用户信息
	cmf.Get("/api/currentUser", new(admin.User).CurrentUser, middleware.ValidationBearerToken, middleware.ValidationAdmin)

	cmf.Get("/sync", func(c *gin.Context) {
		migrate.StartMigrate()
		c.JSON(http.StatusOK, gin.H{
			"msg": "操作成功！",
		})
	})

	cmf.Get("/test", new(admin.Test).Get)
	cmf.Get("/api/v1/region", new(common.RegionController).Get)
	cmf.Get("/api/v1/region/:id", new(common.RegionController).Show)

	common.RegisterOauthRouter()
}
