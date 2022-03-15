package home

import (
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/view"
)

type About struct {
	view.Template
}

//首页控制器
func (web *About) Index(c *gin.Context) {
	view := web.GetView(c)
	view.Assign("seoTitle","首页")
	view.Fetch("about.html")
}
