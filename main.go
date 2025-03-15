package main

import (
	"ThinkTankCentral/core"
	"ThinkTankCentral/global"
	"ThinkTankCentral/initialize"
)

func main() {
	global.Config = core.InitConf() //从 YAML 文件加载配置
	global.Log = core.InitLogger()  // 启动日志模块
	initialize.OtherInit()
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()

	defer global.Redis.Close()
	initialize.InitCron()
	core.RunServer() //启动服务器
}
