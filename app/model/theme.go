/**
** @创建时间: 2021/1/7 2:12 下午
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

type Theme struct {
	Id        int     `json:"id"`
	Name      string  `gorm:"type:varchar(40);comment:主题名称;not null" json:"name"`
	Version   string  `gorm:"type:varchar(10);comment:主题版本;not null" json:"version"`
	Thumbnail string  `gorm:"type:varchar(255);comment:主题缩略图;not null" json:"thumbnail"`
	CreateAt  int64   `gorm:"type:int(10);comment:创建时间;default:0" json:"create_at"`
	UpdateAt  int64   `gorm:"type:int(10);comment:更新时间;default:0" json:"update_at"`
	ListOrder float64 `gorm:"type:float;comment:排序;default:10000" json:"list_order"`
	DeleteAt  int64   `gorm:"type:int(10);comment:删除时间;default:0" json:"delete_at"`
}

type ThemeFile struct {
	Id          int     `json:"id"`
	Theme       string  `gorm:"type:varchar(40);comment:主题名称;not null" json:"theme"`
	ListOrder   float64 `gorm:"type:float;comment:排序;default:10000" json:"list_order"`
	IsPublic    int     `gorm:"type:tinyint(3);comment:是否公告部分;not null" json:"is_public"`
	Name        string  `gorm:"type:varchar(20);comment:模板文件名;not null" json:"name"`
	File        string  `gorm:"type:varchar(50);comment:模板文件,对应的魔板渲染页面;not null" json:"file"`
	Description string  `gorm:"type:varchar(255);comment:模板描述;not null" json:"description"`
	More        string  `gorm:"type:longtext;comment:主题文件用户配置文件" json:"more"`
	ConfigMore  string  `gorm:"type:longtext;comment:主题文件默认配置文件" json:"config_more"`
	CreateAt    int64   `gorm:"type:int(11)" json:"create_at"`
	UpdateAt    int64   `gorm:"type:int(11)" json:"update_at"`
}

func (model *Theme) AutoMigrate() {
	cmf.Db().AutoMigrate(&Theme{})
	cmf.Db().AutoMigrate(&ThemeFile{})
}

func (model *ThemeFile) List(query []string, queryArgs []interface{}) ([]ThemeFile, error) {

	var themeFile []ThemeFile

	queryStr := strings.Join(query, " AND ")

	result := cmf.Db().Where(queryStr, queryArgs...).Find(&themeFile)

	if result.Error != nil {
		return themeFile, result.Error
	}

	return themeFile, nil
}

func (model *ThemeFile) Show(query []string, queryArgs []interface{}) (ThemeFile, error) {

	themeFile := ThemeFile{}

	queryStr := strings.Join(query, " AND ")

	result := cmf.Db().Where(queryStr, queryArgs...).First(&themeFile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return themeFile, errors.New("该分类不存在！")
		}
		return themeFile, result.Error
	}

	return themeFile, nil
}
