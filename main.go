package main

import (
	"ThinkTankCentral/core"
	"ThinkTankCentral/global"
)

func main() {
	global.Config = core.InitConf() //从 YAML 文件加载配置
	global.Log = core.InitLogger()  // 启动日志模块

	core.RunServer() //启动服务器
}
