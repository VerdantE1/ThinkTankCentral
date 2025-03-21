package initialize

import (
	"ThinkTankCentral/global"
	"ThinkTankCentral/middleware"
	"ThinkTankCentral/router"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	//设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	//使用gin会话路由
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret)) // 创建一个基于 Cookie 的会话存储（Session Store）
	Router.Use(sessions.Sessions(global.Config.System.Env, store))

	// 将指定目录下的文件提供给客户端
	// "uploads" 是URL路径前缀，http.Dir("uploads")是实际文件系统中存储文件的目录
	Router.StaticFS(global.Config.Upload.Path, http.Dir(global.Config.Upload.Path))

	//路由组启动器 （无结构，用于Specify路由器）
	routerGroup := router.RouterGroupApp

	//一个路由组实例，api路由组
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	{
		routerGroup.InitBaseRouter(publicGroup) //将publicGroup初始化为baseGroup

	}

	return Router
}
