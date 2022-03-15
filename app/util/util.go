package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	cmf "github.com/gincmf/cmf/bootstrap"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//获取当前登录管理员id
func CurrentAdminId(c *gin.Context) string {
	userId, _ := c.Get("user_id")
	return userId.(string)
}

type role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// 获取当前用户角色
func CurrentRole(c *gin.Context) []role {
	userId, _ := c.Get("user_id")
	userIdInt, _ := strconv.Atoi(userId.(string))
	return GetRoleById(userIdInt)
}

// 根据用户id获取所有角色
func GetRoleById(userId int) []role {
	var result []role
	prefix := cmf.Conf().Database.Prefix
	cmf.Db().Table(prefix+"role_user ru").Select("r.id,r.name").
		Joins("INNER JOIN "+prefix+"role r ON ru.role_id = r.id").
		Where("user_id = ?", userId).
		Scan(&result)
	return result
}

// 是否超级管理员
func SuperRole(c *gin.Context, t int) bool {
	type resultStruct struct {
		Id   int    `json:"id"`
		name string `json:"name"`
	}
	var result []resultStruct
	userId, _ := c.Get("user_id")

	if userId == "1" {
		return true
	}

	prefix := cmf.Conf().Database.Prefix
	cmf.Db().Table(prefix+"role_user ru").Select("r.id,r.name").
		Joins("INNER JOIN "+prefix+"role r ON ru.role_id = r.id").
		Where("ru.user_id = ?", userId).
		Scan(&result)
	for _, v := range result {
		if v.Id == t {
			return true
		}
	}
	return false
}

// 获取真实路径
func CurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// 获取文件是否存在
func ExistPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		fmt.Println(err.Error())
		return false, nil
	}
	fmt.Println(err.Error())
	return false, err
}

// 获取真实url
func GetFileUrl(path string) string {

	if path == "" {
		return ""
	}
	domain := cmf.Conf().App.Domain
	prevPath := domain + "/uploads/" + path
	return prevPath
}

// 获取绝对地址
func GetAbsPath(path string) string {
	if path == "" {
		return ""
	}

	prevPath := CurrentPath() + "/public/uploads/" + path
	return prevPath
}

// 去除空格回车
func TrimAll(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return s
}

func ToLowerInArray(search string, arr []string) bool {
	for _, item := range arr {
		if strings.ToLower(search) == strings.ToLower(item) {
			return true
		}
	}
	return false
}

/**
 * @Author return <1140444693@qq.com>
 * @Description 获取唯一redis生成的编号
 * @Date 2020/11/2 12:27:02
 * @Param
 * @return
 **/

func CurrentDate() (string, string, string) {
	year, month, day := time.Now().Date()
	yearStr := strconv.Itoa(year)
	monthStr := strconv.Itoa(int(month))
	if month < 10 {
		monthStr = "0" + monthStr
	}
	dayStr := strconv.Itoa(day)
	if day < 10 {
		dayStr = "0" + dayStr
	}

	return yearStr, monthStr, dayStr
}

type files struct {
	Article []File `json:"article"`
	List    []File `json:"list"`
	Page    []File `json:"page"`
}

type File struct {
	Name string `json:"name"`
	File string `json:"file"`
}

type sManifest struct {
	Files files `json:"files"`
}

func ReadThemeFiles(theme string) sManifest {
	manifest, err := ioutil.ReadFile("./data/theme/"+theme+"/manifest.json")
	if err != nil {
		fmt.Println(err)
	}

	res := sManifest{}

	json.Unmarshal(manifest, &res)

	return res
}
