package middleware

import (
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	"github.com/gincmf/cmf/view"
)

func BaseController(c *gin.Context) {

}

func HomeBaseController(c *gin.Context) {
	tmpl :=  "/"+cmf.TemplateMap.ThemePath+"/"+cmf.TemplateMap.Theme
	template := new(view.Template).Assign("tmpl",tmpl) //静态资源路径
	c.Set("template",template)
	c.Next()
}

func ApiBaseController(c *gin.Context) {
	ApiController(c)
	AdminController(c)
}

func ApiController(c *gin.Context) {
	BaseController(c)
}

func AdminController(c *gin.Context) {

}
