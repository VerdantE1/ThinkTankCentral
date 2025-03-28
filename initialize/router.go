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

	//基础公共路由
	publicGroup := Router.Group(global.Config.System.RouterPrefix)

	//特定路由
	privateGroup := Router.Group(global.Config.System.RouterPrefix)
	privateGroup.Use(middleware.JWTAuth())

	//管理员路由
	adminGroup := Router.Group(global.Config.System.RouterPrefix)
	adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
	{
		routerGroup.InitBaseRouter(publicGroup) //将publicGroup初始化为baseGroup
	}
	{
		routerGroup.InitUserRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitArticleRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitCommentRouter(privateGroup, publicGroup, adminGroup)
	}
	{
		routerGroup.InitImageRouter(adminGroup)
	}
	return Router
}
