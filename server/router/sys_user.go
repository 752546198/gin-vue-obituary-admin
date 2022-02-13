package router

import (
	v1 "forever.love/api/v1"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	userApi := new(v1.User)
	{
		userRouter.POST("register", userApi.Register)
	}
}
