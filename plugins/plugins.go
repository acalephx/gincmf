/**
** @创建时间: 2020/10/29 4:53 下午
** @作者　　: return
** @描述　　:
 */
package plugins

import (
	portalPlugin "gincmf/plugins/portalPlugin"
)

func AutoRegister()  {
	// 注册路由
	portalPlugin.Router()
	AutoMigrate()
}

func AutoMigrate()  {
	// 注册数据库迁移
	portalPlugin.AutoMigrate()
}
