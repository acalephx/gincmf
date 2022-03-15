/**
** @创建时间: 2021/3/10 9:37 上午
** @作者　　: return
** @描述　　:
 */
package app

import (
	"encoding/json"
	"gincmf/app/model"
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	"github.com/gincmf/cmf/controller"
)

type Settings struct {
	rc controller.Rest
}

func (rest *Settings) Get(c *gin.Context) {
	option := &model.Option{}
	tx := cmf.Db().First(option, "option_name = ?", "site_info") // 查询

	info := model.SiteInfo{}

	json.Unmarshal([]byte(option.OptionValue),&info)

	if tx.RowsAffected > 0 {
		rest.rc.Success(c, "获取成功", info)
	} else {
		rest.rc.Error(c, "获取失败", nil)
	}
}
