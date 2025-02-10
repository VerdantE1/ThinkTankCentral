package main

import (
	"ThinkTankCentral/config"
	"ThinkTankCentral/routes"
)

func main() {

	// 加载环境变量
	config.LoadEnv()

	// 初始化数据库
	config.InitDatabase()

	// 设置路由
	r := routes.SetupRouter()

	// 获取端口号
	port := config.GetEnv("APP_PORT")

	// 运行服务
	r.Run(":" + port)

}
