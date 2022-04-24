package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context)  {
	// 将grpc 的code 转换成http的状态码
	if err != nil{
		if e, ok := status.FromError(err); ok{
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
				"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误" ,
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError,
					gin.H{
					"msg": "用户服务不可用",
					})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.Code(),
				})
			}

		}
		return
	}
}


func HandleValidatorError(err error, c *gin.Context){

	// 如何返回错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON( http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errs.Translate(global.Trans),
	})
}

func GetUserList(ctx * gin.Context)  {
	zap.S().Debugf("获取用户列表页。。。")

	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d",global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg ", err.Error())
	}
	claims ,_ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户的id: %d", currentUser.ID)
	// 调用接口
	// 生成grpc 的client 并调用接口

	pn := ctx.DefaultQuery("pn", "0")
	pnInt , _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt , _ := strconv.Atoi(pSize)

	userClient := proto.NewUserClient(userConn)
	rsp, err :=userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn: uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil{
		zap.S().Errorw("[GetUserList] 查询用户列表失败...")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data{
		//data := make(map[string]interface{})
		user := response.UserResponse{
			Id: value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender: value.Gender,
			Mobile: value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthDay"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}


func PassWordLogin(c *gin.Context){
	zap.S().Infof("用户登录校验....")
	// 表单验证
	passwordLogForm := forms.PassWordLoginForm{}

	if err := c.ShouldBindJSON(&passwordLogForm); err != nil{
		HandleValidatorError(err, c)
		return
	}
	fmt.Println("验证码: ", passwordLogForm.Captcha, "id ", passwordLogForm.CaptchaId)
	// 登录前先用验证码验证
	if  !store.Verify(passwordLogForm.CaptchaId, passwordLogForm.Captcha,  true){
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d",global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg ", err.Error())
	}
	userClient := proto.NewUserClient(userConn)
	// 登录逻辑
	rsp , err := userClient.GetUserByMobile(context.Background(), &proto.MobileReques{
		Mobile: passwordLogForm.Mobile,
	})
	if err != nil{
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest,map[string]string{
				"mobile": "用户手机不存在",
				} )
			default:
				c.JSON(http.StatusInternalServerError,gin.H{
				"mobile":"登录失败",
				})

			}
			return
		}
	}else {
		// 只是查询到了用户， 未校验密码
		//fmt.Println("密码 === ", passwordLogForm.PassWord)
		//fmt.Println("加密密码 === ", rsp.PassWord)
		passRsp, pasErr := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password: passwordLogForm.PassWord,
			EncryptPassword: rsp.PassWord,
		})
		if pasErr != nil{
			c.JSON(http.StatusOK, gin.H{
				"password": "登录失败",
			})
		}
		if passRsp.Success {
			// 生成token
			j := middlewares.NewJWT()
			claims := models.CustomClaims{
				ID: uint(rsp.Id),
				NickName: rsp.NickName,
				AuthorityId: uint(rsp.Role),
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(), // 签名的生效时间
					ExpiresAt: time.Now().Unix() + 60 * 60 *24*30,// 30天过期
					Issuer: "imooc",
				},
			}
			token, err := j.CreateToken(claims)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":"生成token失败",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id": rsp.Id,
				"nick_name": rsp.NickName,
				"token": token,
				"expired_at": (time.Now().Unix() + 60 * 60 *24*30) * 1000,
			})
		}else {
				c.JSON(http.StatusOK, map[string]string {
					"msg": "登录失败",
				})
			}
		}
}

func Register(c *gin.Context){
	// 用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBindJSON(&registerForm); err != nil{
		HandleValidatorError(err, c)
		return
	}
	// 验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(registerForm.Mobile).Result()
	if err == redis.Nil{
		zap.S().Errorf("key 不存在")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"验证码错误",
		})
		return
	}else {
		if value != registerForm.Code{
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":"验证码错误",
			})
		}
	}

	fmt.Println("验证码value = ", value)

	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d",global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg ", err.Error())
	}
	userClient := proto.NewUserClient(userConn)

	user, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.Mobile,
		Mobile: registerForm.Mobile,

	})
	if err != nil{
		zap.S().Errorw("[Register] 注册失败", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	fmt.Println("注册的用户 : ", user)
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: uint(user.Id),
		NickName: user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60 * 60 *24*30,// 30天过期
			Issuer: "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":"生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.Id,
		"nick_name": user.NickName,
		"token": token,
		"expired_at": (time.Now().Unix() + 60 * 60 *24*30) * 1000,
	})
}