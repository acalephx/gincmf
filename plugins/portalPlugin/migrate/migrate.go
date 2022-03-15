 /**
** @创建时间: 2020/10/29 4:51 下午
** @作者　　: return
** @描述　　:
 */

package migrate

import (
	"gincmf/plugins/portalPlugin/model"
	cmf "github.com/gincmf/cmf/bootstrap"
)

func AutoMigrate () {
	cmf.Db().AutoMigrate(&model.PortalCategory{})
	cmf.Db().AutoMigrate(&model.PortalPost{})
	cmf.Db().AutoMigrate(&model.PortalCategoryPost{})
	new(model.PortalTag).AutoMigrate()
}
