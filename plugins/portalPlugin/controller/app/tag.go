/**
** @创建时间: 2021/2/6 10:32 上午
** @作者　　: return
** @描述　　:
 */
package app

import (
	"errors"
	"gincmf/plugins/portalPlugin/model"
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	"github.com/gincmf/cmf/controller"
	"gorm.io/gorm"
)

type Tag struct {
	rc controller.Rest
}

func (rest *Tag) List(c *gin.Context) {

	var tag []model.PortalTag

	cmf.Db().Find(&tag)

	rest.rc.Success(c, "获取成功！", tag)

}

func (rest *Tag) Show(c *gin.Context) {

	var rewrite struct {
		Id int `uri:"id"`
	}

	if err := c.ShouldBindUri(&rewrite); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	var tag model.PortalTag

	tx := cmf.Db().Where("id",rewrite.Id).First(&tag)

	if tx.Error != nil && !errors.Is(tx.Error,gorm.ErrRecordNotFound) {
		rest.rc.Error(c,tx.Error.Error(),nil)
		return
	}

	if tx.RowsAffected == 0 {
		rest.rc.ErrorCode(c,404,"内容不存在！",nil)
		return
	}

	rest.rc.Success(c,"获取成功！",tag)

}

