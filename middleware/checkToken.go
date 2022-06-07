package middleware

import (
	"drop_os_back/dao"
	"drop_os_back/module"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context){
	token := module.GetTokenFromHeader(c.Request.Header)
	//fmt.Println(token)
	if token == "none" || token == "null" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "未登录",
		})
		c.Abort()
		return
	}
	stuNum := dao.Redis.Get(token).Val()
	//fmt.Println("stuNum",stuNum)
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		c.Abort()
		return
	}
	module.ResetTokenTime(token, stuNum)
	//c.Header("sssss","sssss")
	c.Next()
}
