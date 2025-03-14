package initialize

import (
	"ThinkTankCentral/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	/* 默认路由逻辑 */
	Router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	Router.LoadHTMLGlob("static/*")

	return Router
}
