package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	myvalidator "mxshop-api/user-web/validator"
)

func main() {

	port := 8083

	// 1. 初始化logger
	initialize.InitLogger()
	//初始化配置文件
	initialize.InitConfig()
	// 初始化routers
	Router := initialize.Routers()

	// 初始化翻译
	_ = initialize.InitTrans("zh")

	zap.S().Infof("手机验证....")
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	//logger, _ := zap.NewProduction()
	//
	//defer logger.Sync()
	//sugar := logger.Sugar()
	/*
		S()可以获取一个全局的sugar, 可以让自己设置全局Logger
		S, L 安全
	 */
	zap.S().Debugf("启动服务器, 端口号: %d", global.ServerConfig.Port)

	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil{
		//sugar.Infof("启动失败...")
		zap.S().Panic("启动失败：", err.Error())
	}
}
