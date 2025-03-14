package core

import (
	"ThinkTankCentral/global"
	"ThinkTankCentral/initialize"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type server interface {
	ListenAndServe() error
}

// RunServer 用于启动服务器
func RunServer() {
	addr := global.Config.System.Addr() //初始化地址
	Router := initialize.InitRouter()   //初始化路由

	// 加载所有的 JWT 黑名单，存入本地缓存
	// TODO service.LoadAll()

	// 初始化服务器并启动
	s := initServer(addr, Router)
	global.Log.Info("server run success on ", zap.String("address", addr))
	global.Log.Error(s.ListenAndServe().Error())
}

// initServer 函数初始化一个标准的 HTTP 服务器
func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,          // 设置服务器监听的地址
		Handler:        router,           // 设置请求处理器（路由）
		ReadTimeout:    10 * time.Minute, // 设置请求的读取超时时间为 10 分钟
		WriteTimeout:   10 * time.Minute, // 设置响应的写入超时时间为 10 分钟
		MaxHeaderBytes: 1 << 20,          // 设置最大请求头的大小（1MB）
	}
}
