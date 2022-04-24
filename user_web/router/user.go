package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup)  {
	UserRouter := Router.Group("user") //.Use(middlewares.JWTAuth()) 此配置使得所有的接口都登录
	zap.S().Infof("配置用户的相关信息")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth() ,api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}

}
