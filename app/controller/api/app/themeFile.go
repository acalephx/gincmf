/**
** @创建时间: 2021/1/7 3:08 下午
** @作者　　: return
** @描述　　:
 */
package app

import (
	"errors"
	"gincmf/app/model"
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/controller"
	"gorm.io/gorm"
)

type ThemeFile struct {
	rc controller.Rest
}

func (rest *ThemeFile) List(c *gin.Context) {

	theme := c.Query("theme")

	if theme == "" {
		rest.rc.Error(c,"主题不能为空！",nil)
		return
	}

	isPublic := c.Query("is_public")

	query := []string{"theme = ? AND is_public = ?"}
	queryArgs := []interface{}{theme, isPublic}

	data,err := new(model.ThemeFile).List(query,queryArgs)

	if err != nil && !errors.Is(err,gorm.ErrRecordNotFound) {
		rest.rc.Error(c,err.Error(),nil)
		return
	}

	rest.rc.Success(c,"获取成功！",data)


}

func (rest *ThemeFile) Detail(c *gin.Context) {

	theme := c.Query("theme")

	if theme == "" {
		rest.rc.Error(c,"主题不能为空！",nil)
		return
	}

	file := c.Query("file")
	if file == "" {
		rest.rc.Error(c,"文件不能为空！",nil)
		return
	}

	query := []string{"theme = ? AND file = ?"}
	queryArgs := []interface{}{theme, file}

	data,err := new(model.ThemeFile).Show(query,queryArgs)

	if err != nil {
		rest.rc.Error(c,err.Error(),nil)
		return
	}

	rest.rc.Success(c,"获取成功！",data)

}
