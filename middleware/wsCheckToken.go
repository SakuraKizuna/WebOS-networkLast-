package middleware

import (
	"drop_os_back/dao"
	"drop_os_back/module"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WsCheckToken(c *gin.Context){
	token := module.WsGetTokenFromHeader(c.Request.Header)
	//fmt.Println("header:",token)
	if token == "none"{
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "未登录",
		})
		c.Abort()
		return
	}
	stuNum := dao.Redis.Get(token).Val()
	//fmt.Println(stuNum)
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		c.Abort()
		return
	}
	module.ResetTokenTime(token, stuNum)
	//c.Header("Sec-WebSocket-Protocol", token)
	c.Next()
}
