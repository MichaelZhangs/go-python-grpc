package api

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"math/rand"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"time"
)

func CreateCaptcha() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func SendSms(c * gin.Context) {
	// sms表单验证
	// 表单验证
	sendSmsForm := forms.SendSmsForm{}

	if err := c.ShouldBindJSON(&sendSmsForm); err != nil{
		HandleValidatorError(err, c)
		return
	}

	//client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect)
	//if err != nil {
	//	panic(err)
	//}

	smsCode := CreateCaptcha()
	mobile := sendSmsForm.Mobile
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = "18782222221"                          //手机号
	request.QueryParams["SignName"] = "开发者"                              //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_18782222221"                          //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + "777777" + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	//response, err := client.ProcessCommonRequest(request)
	//fmt.Print(client.DoAction(request, response))
	////  fmt.Print(response)
	//if err != nil {
	//	fmt.Print(err.Error())
	//}
	//fmt.Printf("response is %#v\n", response)
	//json数据解析

	fmt.Println("mobile = ", mobile)
	fmt.Println("smsCode = ", smsCode)
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(mobile, smsCode , time.Duration(global.ServerConfig.RedisInfo.Expire)* time.Second)

	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
