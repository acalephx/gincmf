/**
** @创建时间: 2020/12/25 2:08 下午
** @作者　　: return
** @描述　　:
 */
package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	cmfModel "github.com/gincmf/cmf/model"
	"gorm.io/gorm"
	"strings"
)

// 标签内容
type PortalTag struct {
	Id          int    `json:"id"`
	Status      int    `gorm:"type:tinyint(3);comment:状态,1:发布,0:不发布;default:1;not null;" json:"status"`
	Recommended int    `gorm:"type:tinyint(3);comment:是否推荐,1:推荐;0:不推荐;default:0;not null;" json:"recommended"`
	PostCount   int64  `gorm:"type:bigint(20);comment:标签文章数;default:0;not null;" json:"post_count"`
	Name        string `gorm:"type:varchar(20);comment:标签名称;not null;" json:"name"`
}

// 标签关系
type PortalTagPost struct {
	Id     int `json:"id"`
	TagId  int `gorm:"type:bigint(20);comment:标签id;not null;" json:"tag_id"`
	PostId int `gorm:"type:bigint(20);comment:文章id;not null;" json:"post_id"`
	Status int `gorm:"type:tinyint(3);comment:状态,1:发布,0:不发布;default:1;not null;" json:"status"`
}

type PostTagResult struct {
	PortalTagPost
	Name string `gorm:"->" json:"name"`
}

func (model *PortalTag) Error() string {
	panic("implement me")
}

func (model *PortalTag) AutoMigrate() {
	cmf.Db().AutoMigrate(&model)
	cmf.Db().AutoMigrate(&PortalTagPost{})
}

func (model *PortalTag) Index(c *gin.Context, query []string, queryArgs []interface{}) (cmfModel.Paginate, error) {
	// 获取默认的系统分页
	current, pageSize, err := new(cmfModel.Paginate).Default(c)

	if err != nil {
		return cmfModel.Paginate{}, err
	}

	// 合并参数合计
	queryStr := strings.Join(query, " AND ")
	var total int64 = 0

	var tag []PortalTag
	cmf.Db().Where(queryStr, queryArgs...).Find(&tag).Count(&total)

	tx := cmf.Db().Where(queryStr, queryArgs...).Limit(pageSize).Offset((current - 1) * pageSize).Find(&tag)

	if tx.Error != nil {

		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return cmfModel.Paginate{}, tx.Error
		}

	}

	paginate := cmfModel.Paginate{Data: tag, Current: current, PageSize: pageSize, Total: total}
	if len(tag) == 0 {
		paginate.Data = make([]string, 0)
	}
	return paginate, nil
}

/**
 * @Author return <1140444693@qq.com>
 * @Description // 根据文章id获取标签数
 * @Date 2021/1/22 16:15:36
 * @Param
 * @return
 **/
func (model *PortalTag) ListByPostId(postId int) ([]PostTagResult, error) {

	prefix := cmf.Conf().Database.Prefix

	var postTagResult []PostTagResult
	tx := cmf.Db().Table(prefix+"portal_tag_post tp").
		Joins("INNER JOIN "+prefix+"portal_tag t ON t.id = tp.tag_id").
		Joins("INNER JOIN "+prefix+"portal_post p ON p.id = tp.post_id").
		Where("tp.post_id = ? AND  tp.status = 1",postId).
		Order("p.list_order desc,p.id desc").
		Scan(&postTagResult)

	if tx.Error != nil {
		return postTagResult,tx.Error
	}

	return postTagResult,nil

}

func (model *PortalTag) Show(query []string, queryArgs []interface{}) (PortalTag, error) {
	tag := PortalTag{}
	queryStr := strings.Join(query, " AND ")
	result := cmf.Db().Where(queryStr, queryArgs...).Find(&tag)

	if result.Error != nil {
		return tag, nil
	}

	return tag, nil
}

func (model PortalTag) Save(postId int) error {

	var count int64
	cmf.Db().Model(&PortalTagPost{}).Where("post_id = ?", postId).Count(&count)

	cmf.Db().Where("id = ?",model.Id).First(&model)

	model.PostCount = count
	tx := cmf.Db().Save(&model)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (model PortalTag) FirstOrSave() (PortalTag, error) {
	// 新建

	var tx *gorm.DB

	if model.Id == 0 {
		tx = cmf.Db().Create(&model)
	} else {
		// 更新
		// 统计文章标签数
		var count int64
		cmf.Db().Model(&PortalTagPost{}).Where("tag_id = ?", model.Id).Count(&count)

		model.PostCount = count

		tx = cmf.Db().Where("id = ?", model.Id).Save(model)
	}

	if tx.Error != nil {
		return PortalTag{}, tx.Error
	}

	return model, nil

}

func (model *PortalTagPost) FirstOrSave(kId []int) error {

	// [0,1,2]  [1,3,4]

	postId := model.PostId

	// 查出原来的
	var tagPost []PortalTagPost
	cmf.Db().Where("post_id = ?", postId).Find(&tagPost)

	// 待添加的
	var toAdd []PortalTagPost

	for _, v := range kId {
		if !new(PortalTagPost).inAddArray(v,tagPost) || len(tagPost) == 0 {
			toAdd = append(toAdd,PortalTagPost{
				TagId: v,
				PostId: postId,
			})
		}
	}


	//待删除的
	var toDel []string
	var toDelArgs []interface{}

	for _, v := range tagPost {
		if !new(PortalTagPost).inDelArray(v.Id, kId) {
			toDel = append(toDel,"id = ?")
			toDelArgs = append(toDelArgs,v.Id)
		}

		if len(kId) == 0 {
			toDel = append(toDel,"id = ?")
			toDelArgs = append(toDelArgs,v.Id)
		}
	}

	// 删除要删除的
	if len(toDel) > 0 {
		delStr := strings.Join(toDel," OR ")
		cmf.Db().Where(delStr,toDelArgs...).Delete(&PortalTagPost{})
	}

	if len(toAdd) > 0 {
		// 增加待增加的
		cmf.Db().Debug().Create(toAdd)
	}

	// 统计当前标签文章数

	for _, v := range kId {
		err := PortalTag{Id: v}.Save(postId)
		if err != nil {
			return err
		}

	}
	return nil
}

func (model *PortalTagPost) inDelArray(s int, kId []int) bool {

	for _, v := range kId {
		if s == v {
			return true
		}
	}
	return false

}

func (model *PortalTagPost) inAddArray(s int, tagPost []PortalTagPost) bool {

	for _, v := range tagPost {
		if s == v.Id {
			return true
		}
	}
	return false

}
