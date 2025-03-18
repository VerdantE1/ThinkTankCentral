package initialize

import (
	"ThinkTankCentral/global"
	"ThinkTankCentral/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	//路由组启动器 （无结构，用于Specify路由器）
	routerGroup := router.RouterGroupApp

	//一个路由组实例，api路由组
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	{
		routerGroup.InitBaseRouter(publicGroup) //将publicGroup初始化为baseGroup

	}

	return Router
}
