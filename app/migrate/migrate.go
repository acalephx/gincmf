package migrate

import (
	"fmt"
	"gincmf/app/model"
	"gincmf/plugins"
	portalModel "gincmf/plugins/portalPlugin/model"
	"os"
)

type Migrate interface {
	AutoMigrate()
}

func AutoMigrate() {
	_, err := os.Stat("./data/install.lock")
	if err != nil || true {
		StartMigrate()
	}
	// 改为已安装
	file, error := os.Create("./data/install.lock")
	if error != nil {
		fmt.Println(error)
	}
	file.Close()
}

func initTheme()  {

	//bytes, err := ioutil.ReadFile("./data/conf/menu.json")
	//if err != nil {
	//	return
	//}

}

func StartMigrate()  {
	new(option).AutoMigrate()
	new(user).AutoMigrate()
	new(asset).AutoMigrate()
	new(role).AutoMigrate()
	new(authAccess).AutoMigrate()
	new(AdminMenu).AutoMigrate()
	new(Region).AutoMigrate()
	new(portalModel.PortalTag).AutoMigrate()
	new(model.Route).AutoMigrate()
	new(model.Nav).AutoMigrate()
	new(model.Theme).AutoMigrate()
	// 插件数据库迁移注册
	plugins.AutoMigrate()
}
