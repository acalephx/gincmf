/**
** @创建时间: 2020/10/29 4:33 下午
** @作者　　: return
** @描述　　:
 */

package router

import (
	"gincmf/app/middleware"
	"gincmf/plugins/portalPlugin/controller/admin"
	"gincmf/plugins/portalPlugin/controller/app"
	cmf "github.com/gincmf/cmf/bootstrap"
)

func ApiListenRouter() {

	adminGroup := cmf.Group("api/admin/v1", middleware.ValidationBearerToken, middleware.ApiBaseController)
	{
		adminGroup.Rest("/portal/category", new(admin.Category))
		adminGroup.Get("/portal/category_list", new(admin.Category).List)
		adminGroup.Get("/portal/category_options", new(admin.Category).Options)
		adminGroup.Rest("/portal/article", new(admin.PortalPost))
		adminGroup.Get("/themeFile/list", new(admin.ThemeFile).List)
		adminGroup.Rest("/portal/tag", new(admin.Tag))
	}

	appGroup := cmf.Group("api/app/v1", middleware.ApiBaseController)
	{

		appGroup.Get("/portal/category", new(app.Category).List)
		appGroup.Get("/portal/category/:id", new(app.Category).Show)

		appGroup.Get("/portal/top_category_id/:id", new(app.Category).GetTopId)
		appGroup.Get("/portal/list/:id", new(app.Post).Get) // 根据id获取分页列表
		appGroup.Post("/portal/list_with_id", new(app.Post).ListWithCid)

		appGroup.Get("/portal/tag_id/:id", new(app.Post).ListByTag)

		appGroup.Get("/portal/top_category_list", new(app.Category).GetTopList)

		appGroup.Get("/portal/article/:id", new(app.Post).Show)

		appGroup.Get("/portal/page/:id", new(app.Post).Page)

		appGroup.Get("/portal/tag", new(app.Tag).List)

		appGroup.Get("/portal/tag/:id", new(app.Tag).Show)

	}

}
