package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
)

func InitConfig()  {

	configFilename := fmt.Sprintf("user-web/config-debug.yaml")
	v := viper.New()
	v.SetConfigFile(configFilename)
	if err := v.ReadInConfig(); err != nil{
		panic(err)
	}

	serverConfig := global.ServerConfig

	if  err := v.Unmarshal(serverConfig); err != nil{
		panic(err)
	}
	zap.S().Infof("host = %s port = %d,",serverConfig.UserSrvInfo.Host, serverConfig.UserSrvInfo.Port)
}
