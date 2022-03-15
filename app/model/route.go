/**
** @创建时间: 2020/12/26 3:00 下午
** @作者　　: return
** @描述　　:
 */
package model

import (
	"errors"
	cmf "github.com/gincmf/cmf/bootstrap"
	"gorm.io/gorm"
	"strings"
)

type Route struct {
	Id           int     `json:"id"`
	ListOrder    float64 `gorm:"type:float;comment:排序;default:10000" json:"list_order"`
	Status       int     `gorm:"type:tinyint(3);comment:状态;default:1;;not null" json:"status"`
	Type         int     `gorm:"type:tinyint(4);comment:URL规则类型;1:用户自定义;2:别名添加';default:1;;not null" json:"type"`
	FullUrl      string  `gorm:"type:varchar(255);comment:完整url;not null" json:"full_url"`
	Url          string  `gorm:"type:varchar(255);comment:实际显示的url;not null" json:"url"`
	Template     string  `gorm:"type:varchar(40);comment:模板;not null" json:"template"`
	TemplateType string  `gorm:"type:varchar(10);comment:模板类型;not null" json:"template_type"`
}

type RouteResult struct {
	Route
	Template string `json:"template"`
}

func (model *Route) AutoMigrate() {
	cmf.Db().AutoMigrate(&model)
}

func (model *Route) Show(query []string, queryArgs []interface{}) (Route, error) {

	route := Route{}
	queryStr := strings.Join(query, " AND ")
	result := cmf.Db().Where(queryStr, queryArgs...).First(&route)

	if result.Error != nil {
		return route, nil
	}

	return route, nil
}

func (model *Route) List(query []string, queryArgs []interface{}) ([]Route, error) {

	var route []Route
	queryStr := strings.Join(query, " AND ")
	result := cmf.Db().Where(queryStr, queryArgs...).Find(&route)

	if result.Error != nil {
		return route, nil
	}

	return route, nil
}

func (model *Route) Set() error {

	route := Route{
		Type:         2,
		FullUrl:      model.FullUrl,
		Url:          model.Url,
		Template:     model.Template,
		TemplateType: model.TemplateType,
	}

	tx := cmf.Db().Where("full_url", route.FullUrl).First(&route)

	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {

		return tx.Error
	}

	if tx.RowsAffected == 0 {
		cmf.Db().Create(&route)
	} else {
		route.Type = 2
		route.FullUrl =  model.FullUrl
		route.Url =  model.Url
		route.Template = model.Template
		route.TemplateType = model.TemplateType
		tx := cmf.Db().Save(&route)
		if tx.Error != nil {
			return  tx.Error
		}
	}

	return nil
}
