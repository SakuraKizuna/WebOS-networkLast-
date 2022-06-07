package module

import (
	"drop_os_back/dao"
	"drop_os_back/models"
	"net/http"
	"time"
)

func TokenToLevelAndBelong(token interface{}) (level int, belong string) {
	adminnameStruct := dao.Redis.Get(token.(string))
	adminname := adminnameStruct.Val()
	//dao.Redis.Del(token.(string))
	//dao.Redis.Set(token.(string), adminname, 1200*time.Second)
	dao.Redis.Expire(token.(string), 14400*time.Second)
	admin, err := models.GetAdminInfo(adminname)
	if err != nil {
		return -1, "error"
	}
	return admin.Level, admin.Belong
}

//http请求
func GetTokenFromHeader(header http.Header) (token string) {
	token = "none"
	for k, v := range header {
		if k == "Token" {
			token = v[0]
		}
	}
	return token
}

//websocket请求
func WsGetTokenFromHeader(header http.Header) (token string) {
	//token = "none"
	for k, v := range header {
		//fmt.Println(k,":",v[0])
		if k == "Sec-Websocket-Protocol" {
			token = v[0]
		}
	}
	return token
}

//重置token时间
func ResetTokenTime(token, username string) {
	dao.Redis.Expire(token, 14400*time.Second)
	//dao.Redis.Del(token)
	//dao.Redis.Set(token, username, 1200*time.Second)
}
