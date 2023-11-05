package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/oss/global"
	"github.com/Yifangmo/micro-shop-web/oss/utils"
)

func Token(c *gin.Context) {
	response := utils.GetPolicyToken()
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(200, response)
}

func HandlerRequest(ctx *gin.Context) {
	fmt.Println("\nHandle Post Request ... ")
	// Get PublicKey bytes
	bytePublicKey, err := utils.GetPublicKey(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := utils.GetAuthorization(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// Get MD5 bytes from Newly Constructed Authrization String.
	byteMD5, bodyStr, err := utils.GetMD5FromNewAuthString(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	decodeurl, err := url.QueryUnescape(bodyStr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decodeurl)
	params := make(map[string]string)
	datas := strings.Split(decodeurl, "&")
	for _, v := range datas {
		sdatas := strings.Split(v, "=")
		fmt.Println(v)
		params[sdatas[0]] = sdatas[1]
	}
	fileName := params["filename"]
	fileUrl := fmt.Sprintf("%s/%s", global.ServerConfig.Oss.Host, fileName)

	// verifySignature and response to client
	if utils.VerifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		// do something you want accoding to callback_body ...
		ctx.JSON(http.StatusOK, gin.H{
			"url": fileUrl,
		})
		//utils.ResponseSuccess(ctx)  // response OK : 200
	} else {
		utils.ResponseFailed(ctx) // response FAILED : 400
	}
}

func InitOssRouter(Router *gin.RouterGroup) {
	OssRouter := Router.Group("oss")
	{
		//OssRouter.GET("token", middlewares.JWTAuth(), middlewares.IsAdminAuth(), handler.Token)
		OssRouter.GET("token", Token)
		OssRouter.POST("/callback", HandlerRequest)
	}
}
