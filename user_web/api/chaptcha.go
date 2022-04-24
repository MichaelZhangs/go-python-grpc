package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)


var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context)  {

	driver  := base64Captcha.DefaultDriverDigit
	cp := base64Captcha.NewCaptcha(driver, store )
	id, b64s, err := cp.Generate()
	if err != nil{
		zap.S().Errorf("生成验证吗错误, :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证吗错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath": b64s,
	})

}