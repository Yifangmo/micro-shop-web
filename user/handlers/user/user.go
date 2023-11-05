package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Yifangmo/micro-shop-services/common"
	"github.com/Yifangmo/micro-shop-services/user/proto"
	"github.com/Yifangmo/micro-shop-web/user/global"
	"github.com/Yifangmo/micro-shop-web/user/middlewares"
	"github.com/Yifangmo/micro-shop-web/user/models"
	"github.com/Yifangmo/micro-shop-web/user/utils"
)

func List(ctx *gin.Context) {
	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.JWTClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)
	//生成grpc的client并调用接口

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	resp, err := global.UserSrvClient.GetUserList(context.Background(), &common.PageInfo{
		PageNumber: uint32(pnInt),
		PageSize:   uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": resp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range resp.Data {
		user := models.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Birthday: models.JsonTime(value.Birthday.AsTime()),
			Gender:   value.Gender.Display(),
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}

	reMap["data"] = result
	ctx.JSON(http.StatusOK, reMap)
}

func LoginByPassword(c *gin.Context) {
	passwordLoginForm := models.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}

	// if store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"captcha": "验证码错误",
	// 	})
	// 	return
	// }

	if resp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		//只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.CheckPasswordRequest{
			RawPassword:       passwordLoginForm.PassWord,
			EncryptedPassword: resp.Password,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.JWTClaims{
					ID:          uint(resp.Id),
					NickName:    resp.NickName,
					AuthorityId: uint(resp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "test",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":         resp.Id,
					"nick_name":  resp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败",
				})
			}
		}
	}
}

func Register(c *gin.Context) {
	registerForm := models.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}

	//验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.Redis.Host, global.ServerConfig.Redis.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}
	// TODO:
	nickName := ""
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserRequest{
		NickName: registerForm.Mobile,
		Password: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.JWTClaims{
		ID:          uint(user.Id),
		NickName:    registerForm.Mobile,
		AuthorityId: 1,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "test",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  nickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

func Detail(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.JWTClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	resp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.UserIDRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name":     resp.NickName,
		"birthday": resp.Birthday.AsTime().Format("2006-01-02"),
		"gender":   resp.Gender,
		"mobile":   resp.Mobile,
	})
}

func Update(ctx *gin.Context) {
	updateUserForm := models.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.JWTClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserRequest{
		Id:       int32(currentUser.ID),
		NickName: updateUserForm.Name,
		Gender:   proto.NewGender(updateUserForm.Gender),
		Birthday: timestamppb.New(birthDay),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
