/**
** @创建时间: 2021/1/3 11:23 下午
** @作者　　: return
** @描述　　:
 */

package app

import (
	"gincmf/app/model"
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/controller"
)

type Route struct {
	rc controller.Rest
}

func (rest *Route) List(c *gin.Context) {
	var rewrite struct {
		Id int `uri:"id"`
	}
	if err := c.ShouldBindUri(&rewrite); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	data,err := new(model.Route).List(nil,nil)

	if err != nil {
		rest.rc.Error(c,err.Error(),nil)
		return
	}

	rest.rc.Success(c,"获取成功！",data)

}
