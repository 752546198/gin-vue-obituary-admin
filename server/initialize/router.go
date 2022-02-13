package initialize

import (
	"forever.love/router"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Routers 初始化路由
func Routers() *gin.Engine {
	Router := gin.Default()
	// 集成 swagger 自动生成文档
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiV1 := Router.Group("api/v1")
	{
		router.InitUserRouter(apiV1) // 注册用户路由
	}
	return Router
}
