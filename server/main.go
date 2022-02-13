package main

import (
	"fmt"
	"forever.love/initialize"
	"forever.love/initialize/pika"
	"forever.love/initialize/redis"
	viper2 "forever.love/initialize/viper"
	"github.com/gin-gonic/gin"
)

func main() {
	// 设置 gin 框架的模式（viper 初始化会根据模式读取不同的 yaml 配置文件），本地开发可以使用调试模式（debug），上线必须使用 release
	gin.SetMode(gin.DebugMode)
	// 初始化 viper
	viper2.Setup()
	// 初始化 redis
	redis.Setup()
	// 初始化 pika
	pika.Setup()
	// 初始化路由
	routers := initialize.Routers()
	_ = routers.Run(fmt.Sprintf(":%d", viper2.ServerConf.Port))
}
