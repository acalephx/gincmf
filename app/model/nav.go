/**
** @创建时间: 2021/1/3 8:30 下午
** @作者　　: return
** @描述　　:
 */
package model

import (
	"errors"
	"gincmf/app/util"
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	cmfModel "github.com/gincmf/cmf/model"
	"gorm.io/gorm"
	"strings"
)

type Nav struct {
	Id       int    `json:"id"`
	IsMain   int    `gorm:"type:tinyint(3);comment:主导航;not null;default:0" json:"is_main"`
	Name     string `gorm:"type:varchar(20);comment:主导航名称" json:"name"`
	Remark   string `gorm:"type:varchar(255);comment:备注" json:"remark"`
	paginate cmfModel.Paginate
}

type NavItem struct {
	Id        int       `json:"id"`
	NavId     int       `gorm:"type:int(11);comment:导航id;not null" json:"nav_id"`
	ParentId  int       `gorm:"type:int(11);comment:所属父类id;default:0" json:"parent_id"`
	Status    int       `gorm:"type:tinyint(3);default:1" json:"status"`
	ListOrder float64   `gorm:"type:float;comment:排序;default:10000" json:"list_order"`
	Name      string    `gorm:"type:varchar(50);comment:路由名称" json:"name"`
	Target    string    `gorm:"type:varchar(10);comment:目标状态" json:"target"`
	Href      string    `gorm:"type:varchar(100);comment:路由路径" json:"href"`
	Icon      string    `gorm:"type:varchar(255);comment:图标地址" json:"icon"`
	IconPrev  string    `gorm:"-" json:"icon_prev"`
	Path      string    `gorm:"type:varchar(255);comment:路由路径" json:"path"`
	Children  []NavItem `gorm:"-" json:"children,omitempty"`
	paginate  cmfModel.Paginate
}

type NavItemOptions struct {
	Id       int    `json:"id"`
	ParentId int    `gorm:"type:int(11);comment:所属父类id;default:0" json:"parent_id"`
	Name     string `gorm:"type:varchar(50);comment:路由名称" json:"name"`
	Level    int    `json:"level"`
}

func (model *Nav) AutoMigrate() {
	cmf.Db().AutoMigrate(&Nav{})
	cmf.Db().AutoMigrate(&NavItem{})
}

func (model *Nav) Get(c *gin.Context, query []string, queryArgs []interface{}) (cmfModel.Paginate, error) {

	// 获取默认的系统分页
	current, pageSize, err := model.paginate.Default(c)

	if err != nil {
		return cmfModel.Paginate{}, err
	}

	// 合并参数合计
	queryStr := strings.Join(query, " AND ")
	var total int64 = 0

	var nav []Nav
	cmf.Db().Where(queryStr, queryArgs...).Find(&nav).Count(&total)

	result := cmf.Db().Where(queryStr, queryArgs...).Limit(pageSize).Offset((current - 1) * pageSize).Find(&nav)

	if result.Error != nil {

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return cmfModel.Paginate{}, result.Error
		}

	}

	paginate := cmfModel.Paginate{Data: nav, Current: current, PageSize: pageSize, Total: total}
	if len(nav) == 0 {
		paginate.Data = make([]string, 0)
	}

	return paginate, nil

}

func (model *NavItem) Show(query []string, queryArgs []interface{}) (NavItem, error) {

	navItem := NavItem{}

	// 合并参数合计
	queryStr := strings.Join(query, " AND ")

	tx := cmf.Db().Where(queryStr, queryArgs...).Find(&navItem)

	if tx.Error != nil {
		return navItem, tx.Error
	}

	if navItem.Icon != "" {
		navItem.IconPrev = util.GetFileUrl(navItem.Icon)
	}

	href := navItem.Href

	if href != "" && href[0:1] != "/" {
		href = "/" +href
	}

	navItem.Href = href

	return navItem, nil

}

func (model *NavItem) GetWithChild(c *gin.Context, query []string, queryArgs []interface{}) (cmfModel.Paginate, error) {

	// 获取默认的系统分页
	current, pageSize, err := model.paginate.Default(c)

	if err != nil {
		return cmfModel.Paginate{}, err
	}

	// 合并参数合计
	queryStr := strings.Join(query, " AND ")
	var total int64 = 0

	var navItem []NavItem
	cmf.Db().Where(queryStr, queryArgs...).Find(&navItem).Count(&total)
	tx := cmf.Db().Where(queryStr, queryArgs...).Limit(pageSize).Offset((current - 1) * pageSize).Find(&navItem)

	if tx.Error != nil {
		return cmfModel.Paginate{}, tx.Error
	}

	result := model.recursionNav(navItem, 0)

	paginate := cmfModel.Paginate{Data: result, Current: current, PageSize: pageSize, Total: total}
	if len(navItem) == 0 {
		paginate.Data = make([]string, 0)
	}

	return paginate, nil

}

func (model *NavItem) recursionNav(nav []NavItem, parentId int) []NavItem {

	var navItems []NavItem
	// 增加当前层级
	for _, v := range nav {

		href := v.Href

		if parentId == v.ParentId {
			ni := NavItem{
				Id:        v.Id,
				NavId:     v.NavId,
				ParentId:  parentId,
				Status:    v.Status,
				Name:      v.Name,
				Path:      v.Path,
				Target:    v.Target,
				Href:      href,
				Icon:      v.Icon,
				IconPrev: util.GetFileUrl(v.Icon),
				ListOrder: v.ListOrder,
			}

			childNav := model.recursionNav(nav, v.Id)
			ni.Children = childNav

			navItems = append(navItems, ni)
		}

	}

	return navItems
}

var navItemOptions []NavItemOptions

func (model *NavItem) OptionsList(query []string, queryArgs []interface{}) []NavItemOptions {

	var navItem []NavItem

	queryStr := strings.Join(query, " AND ")

	navItemOptions = make([]NavItemOptions, 0)
	cmf.Db().Where(queryStr, queryArgs...).Find(&navItem)

	data := model.recursionOptions(navItem, 0, 0)

	for k, v := range data {
		data[k].Name = model.indent(v.Level) + v.Name
	}

	return data
}

func (model *NavItem) indent(level int) string {

	indent := ""
	for i := 0; i < level; i++ {
		indent += "    |--"
	}

	return indent

}

func (model *NavItem) recursionOptions(nav []NavItem, parentId int, level int) []NavItemOptions {

	nextLevel := level + 1

	for _, v := range nav {

		if parentId == v.ParentId {

			ops := NavItemOptions{
				Id:       v.Id,
				ParentId: v.ParentId,
				Name:     v.Name,
				Level:    level,
			}

			navItemOptions = append(navItemOptions, ops)
			model.recursionOptions(nav, v.Id, nextLevel)
		}
	}

	return navItemOptions

}
