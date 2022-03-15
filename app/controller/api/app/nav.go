/**
** @创建时间: 2021/1/3 8:29 下午
** @作者　　: return
** @描述　　:
 */
package app

import (
	"gincmf/app/model"
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/controller"
)

type Nav struct{
	rc controller.Rest
}

func (rest *NavItem) Get(c *gin.Context) {

	navId := c.Query("nav_id")

	var query = []string{"nav_id  = ?"}
	var queryArgs = []interface{}{navId}

	navItem, err := new(model.NavItem).GetWithChild(c, query, queryArgs)

	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}

	rest.rc.Success(c, "获取成功！", navItem)
}


