/**
** @创建时间: 2020/12/25 11:24 下午
** @作者　　: return
** @描述　　:
 */
package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"gincmf/app/util"
	"gincmf/plugins/portalPlugin/model"
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	"github.com/gincmf/cmf/controller"
	cmfModel "github.com/gincmf/cmf/model"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Post struct {
	rc controller.Rest
}

/**
 * @Author return <1140444693@qq.com>
 * @Description 获取文章列表
 * @Date 2021/1/10 17:10:34
 * @Param
 * @return
 **/

func (rest *Post) Get(c *gin.Context) {

	var rewrite struct {
		Id int `uri:"id"`
	}

	if err := c.ShouldBindUri(&rewrite); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	categoryId := rewrite.Id
	pc := model.PortalCategory{Id: categoryId}
	ids, err := pc.ChildIds()

	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}

	var query []string
	var queryArgs []interface{}

	for _, v := range ids {
		query = append(query, "cp.category_id = ?")
		queryArgs = append(queryArgs, v)
	}

	queryRes := []string{"p.post_type = 1 AND p.delete_at = 0"}

	queryStr := strings.Join(query, " OR ")
	queryRes = append(queryRes, "("+queryStr+")")

	data, err := new(model.PortalPost).IndexByCategory(c, queryRes, queryArgs)

	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}

	rest.rc.Success(c, "获取成功！", data)

}

/**
 * @Author return <1140444693@qq.com>
 * @Description 获取文章列表
 * @Date 2021/1/10 17:10:34
 * @Param
 * @return
 **/
func (rest *Post) ListWithCid(c *gin.Context) {

	var form struct {
		Ids []int `json:"ids"`
	}

	if err := c.BindJSON(&form); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	var query []string
	var queryArgs []interface{}

	for _, v := range form.Ids {
		query = append(query, "cp.category_id = ?")
		queryArgs = append(queryArgs, v)
	}

	queryRes := []string{"p.post_type = 1 AND p.delete_at = 0"}

	if len(query) > 0 {
		queryStr := strings.Join(query, " OR ")
		queryRes = append(queryRes, queryStr)
	}

	data, err := model.PortalPost{}.IndexByCategory(c, queryRes, queryArgs)

	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}

	rest.rc.Success(c, "获取成功！", data)

}

/**
 * @Author return <1140444693@qq.com>
 * @Description 根据标签获取文章列表
 * @Date 2021/2/5 18:32:6
 * @Param
 * @return
 **/

func (rest *Post) ListByTag(c *gin.Context) {

	var rewrite struct {
		Id int `uri:"id"`
	}

	if err := c.ShouldBindUri(&rewrite); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	id := rewrite.Id

	type temp struct {
		model.PortalPost
		ThumbnailPrev  string                 `json:"thumbnail_prev"`
		UserLogin      string                 `json:"user_login"`
		Category       []model.PortalCategory `gorm:"-" json:"category"`
		Tags           []model.PostTagResult  `gorm:"-" json:"tags"`
		MoreJson       model.More             `gorm:"-" json:"more_json"`
		ReferenceQuote string                 `gorm:"-" json:"reference_quote"`
	}

	current, pageSize, err := new(cmfModel.Paginate).Default(c)
	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}

	var tempArr []temp

	prefix := cmf.Conf().Database.Prefix

	var total int64 = 0
	cmf.Db().Debug().Table(prefix+"portal_tag_post tp").
		Joins("inner join "+prefix+"portal_post p ON tp.post_id = p.id").
		Where("tp.tag_id = ? and p.delete_at = 0", id).Count(&total)

	tx := cmf.Db().Table(prefix+"portal_tag_post tp").
		Joins("inner join "+prefix+"portal_post p ON tp.post_id = p.id").
		Where("tp.tag_id = ? and p.delete_at = 0", id).
		Limit(pageSize).Offset((current - 1) * pageSize).Scan(&tempArr)

	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		rest.rc.Error(c, tx.Error.Error(), nil)
		return
	}

	for k, v := range tempArr {

		tempArr[k].ThumbnailPrev = util.GetFileUrl(v.Thumbnail)

		category := model.PortalCategory{}
		categoryItem, _ := category.ListWithPost([]string{"p.id = ? AND  p.delete_at = ?"}, []interface{}{v.Id, 0})
		tempArr[k].Category = categoryItem

		createTime := time.Unix(v.CreateAt, 0).Format("2006-01-02 15:04:05")
		tempArr[k].CreateTime = createTime

		updateTime := time.Unix(v.UpdateAt, 0).Format("2006-01-02 15:04:05")
		tempArr[k].UpdateTime = updateTime

		m := model.More{}
		json.Unmarshal([]byte(v.More), &m)
		tempArr[k].MoreJson = m

		referenceQuote := m.ReferenceQuote

		tempArr[k].ReferenceQuote = referenceQuote

		data, _ := new(model.PortalTag).ListByPostId(v.Id)

		tempArr[k].Tags = data

	}

	paginate := cmfModel.Paginate{Data: tempArr, Current: current, PageSize: pageSize, Total: total}
	if len(tempArr) == 0 {
		paginate.Data = make([]string, 0)
	}

	rest.rc.Success(c, "获取成功！", paginate)

}

func (rest *Post) Show(c *gin.Context) {

	var rewrite struct {
		Id int `uri:"id"`
	}

	if err := c.ShouldBindUri(&rewrite); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	id := rewrite.Id

	var query = []string{"id = ? AND post_type = 1", "delete_at = ?"}
	var queryArgs = []interface{}{id, 0}

	post, err := model.PortalPost{}.Show(query, queryArgs)

	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}


	var result struct {
		model.PortalPost
		ThumbPrevPath string                 `json:"thumb_prev_path"`
		AudioPrevPath string                 `json:"audio_prev_path"`
		VideoPrevPath string                 `json:"video_prev_path"`
		Template      string                 `json:"template"`
		Address       string                 `json:"address"`
		Keywords      []string               `json:"keywords"`
		Photos        []model.Path           `json:"photos"`
		Accessories   []model.Path           `json:"accessories"`
		Files         []model.Path           `json:"files"`
		Audio         string                 `json:"audio"`
		Video         string                 `json:"video"`
		Category      []model.PortalCategory `json:"category"`
		Tags          []model.PostTagResult  `json:"tags"`
		Slug          string                 `json:"slug"`
	}

	json.Unmarshal([]byte(post.More), &result)

	if len(result.Photos) == 0 {
		result.Photos = make([]model.Path, 0)
	}

	result.PortalPost = post

	for k, v := range result.Photos {
		result.Photos[k].PrevPath = util.GetFileUrl(v.FilePath)
		result.MoreJson.Photos[k].PrevPath = util.GetFileUrl(v.FilePath)
	}

	for k, v := range result.MoreJson.Photos {
		result.MoreJson.Photos[k].PrevPath =  util.GetFileUrl(v.FilePath)
	}

	if len(result.Accessories) == 0 {
		result.Accessories = make([]model.Path, 0)
	}

	if len(result.Files) == 0 {
		result.Files = make([]model.Path, 0)
	}

	if post.PostKeywords != "" {
		result.Keywords = strings.Split(post.PostKeywords, ",")
	}

	result.ThumbPrevPath = util.GetFileUrl(result.Thumbnail)
	result.AudioPrevPath = util.GetFileUrl(result.Audio)
	result.VideoPrevPath = util.GetFileUrl(result.Video)

	category, err := new(model.PortalCategory).GetCategoryByPostId(post.Id)
	result.Category = category

	data, _ := new(model.PortalTag).ListByPostId(post.Id)
	result.Tags = data

	rest.rc.Success(c, "获取成功！", result)
}

func (rest *Post) Page(c *gin.Context) {

	var rewrite struct {
		Id int `uri:"id"`
	}

	if err := c.ShouldBindUri(&rewrite); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	id := rewrite.Id

	var query = []string{"id = ? AND post_type = 2", "delete_at = ?"}
	var queryArgs = []interface{}{id, 0}

	post, err := model.PortalPost{}.Show(query, queryArgs)

	if err != nil {
		rest.rc.Error(c, err.Error(), nil)
		return
	}

	var result struct {
		model.PortalPost
		ThumbPrevPath string       `json:"thumb_prev_path"`
		AudioPrevPath string       `json:"audio_prev_path"`
		VideoPrevPath string       `json:"video_prev_path"`
		Template      string       `json:"template"`
		Address       string       `json:"address"`
		Keywords      []string     `json:"keywords"`
		Photos        []model.Path `json:"photos"`
		Files         []model.Path `json:"files"`
		Audio         string       `json:"audio"`
		Video         string       `json:"video"`
		Slug          string       `json:"slug"`
	}

	json.Unmarshal([]byte(post.More), &result)

	if len(result.Photos) == 0 {
		result.Photos = make([]model.Path, 0)
	}

	if len(result.Files) == 0 {
		result.Files = make([]model.Path, 0)
	}

	if post.PostKeywords != "" {
		result.Keywords = strings.Split(post.PostKeywords, ",")
	}

	result.PortalPost = post
	result.ThumbPrevPath = util.GetFileUrl(result.Thumbnail)
	result.AudioPrevPath = util.GetFileUrl(result.Audio)
	result.VideoPrevPath = util.GetFileUrl(result.Video)

	rest.rc.Success(c, "获取成功！", result)
}
