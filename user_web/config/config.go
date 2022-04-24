package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int32  `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type AliSmsConfig struct {
	ApiKey     string `mapstructure:"key"`
	ApiSecrect string `mapstructure:"secrect"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int  `mapstructure:"port"`
	Expire int `mapstructure:"expire"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"host"`
	Port        int32         `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
	AliSmsInfo  AliSmsConfig  `mapstructure:"sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis"`
}
