package controller

import (
	"drop_os_back/dao"
	"drop_os_back/models"
	"drop_os_back/module"
	"drop_os_back/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var mux sync.Mutex

//聊一聊接口
func Handler(c *gin.Context) {
	//获取发送者uid和被发送者uid
	token := module.WsGetTokenFromHeader(c.Request.Header)
	fmt.Println(token)
	stuNum := dao.Redis.Get(token).Val()
	toStuNum := c.Query("toStuNum")
	if stuNum == toStuNum {
		c.JSON(http.StatusOK, gin.H{
			"status": 203,
			"msg":    "自己无法与自己建立沟通",
		})
		return
	}
	//升级websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//将获取的参数放进这个数组
		Subprotocols: []string{token},
	}).Upgrade(c.Writer, c.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户实例
	client := &models.Client{
		ID:           createID(stuNum, toStuNum), //1->2
		SendID:       createID(toStuNum, stuNum), //2->1
		FromUsername: stuNum,
		ToUsername:   toStuNum,
		Socket:       conn,
		Send:         make(chan []byte),
	}
	//用户注册到用户管理上
	mux.Lock()
	models.Manager.Register <- client
	mux.Unlock()
	go client.Read()
	go client.Write()
	go models.BatchRead(toStuNum, stuNum)
}

//websocket发送图片
func PostPic(c *gin.Context) {
	//获取发送者uid和被发送者uid
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	//uid := c.Query("uid")
	toStuNum := c.Query("toStuNum")
	//获取文件头
	file, err := c.FormFile("uploadPicture")
	if err != nil {
		fmt.Println(file)
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	fmt.Println(file.Filename)
	fileFormat := file.Filename[len(file.Filename)-4:]
	if file.Filename[len(file.Filename)-3:] != "jpg" && file.Filename[len(file.Filename)-3:] != "png" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "仅支持上传jpg或者png格式",
		})
		return
	}
	fileUuid := module.GetUUID()
	filePath := "./uploadPic/" + fileUuid + fileFormat
	fmt.Println(filePath)
	realPicFormat := "http://" + c.Request.Host + "/uploadPic/" + fileUuid + fileFormat
	userid := stuNum + "->" + toStuNum
	//_ = models.InsertMsg(userid,realPicFormat,0,int64(3*month))
	fmt.Println(userid)
	var client *models.Client
	for id, conn := range models.Manager.Clients {
		if userid != id {
			continue
		}
		client = conn
	}
	if client == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "发送失败，服务器逻辑层错误",
		})
		return
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	mux.Lock()
	//fmt.Println(client.ID)
	models.Manager.Broadcast <- &models.Broadcast{
		Client:  client,
		Message: []byte(realPicFormat), //发送过来的消息
		Type:    1,
	}
	mux.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "发送成功",
	})
}

//websocket发送文件
func PostFile(c *gin.Context) {
	//获取发送者uid和被发送者uid
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	//uid := c.Query("uid")
	toStuNum := c.Query("toStuNum")
	//获取文件头
	file, err := c.FormFile("uploadFile")
	if err != nil {
		fmt.Println(file)
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	//fmt.Println(file.Filename)
	//fileFormat := file.Filename[len(file.Filename)-4:]
	fileFormatInt := strings.Index(file.Filename, ".")
	fileFormat := file.Filename[fileFormatInt:]
	if fileFormat != ".doc" && fileFormat != ".pdf" && fileFormat != ".ppt" && fileFormat != ".pptx" && fileFormat != ".docx" && fileFormat != ".zip" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "格式不符合，支持格式：doc/docx，pdf，ppt/pptx，zip",
		})
		return
	}
	fileUuid := module.GetUUID()
	filePath := "./uploadPic/" + fileUuid + fileFormat
	fmt.Println(filePath)
	realPicFormat := "http://" + c.Request.Host + "/uploadPic/" + fileUuid + fileFormat
	userid := stuNum + "->" + toStuNum
	//_ = models.InsertMsg(userid,realPicFormat,0,int64(3*month))
	fmt.Println(userid)
	var client *models.Client
	for id, conn := range models.Manager.Clients {
		if userid != id {
			continue
		}
		client = conn
	}
	if client == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "发送失败，服务器逻辑层错误",
		})
		return
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	mux.Lock()
	//fmt.Println(client.ID)
	models.Manager.Broadcast <- &models.Broadcast{
		Client:  client,
		Message: []byte(realPicFormat), //发送过来的消息
		Type:    2,
	}
	mux.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"msg":      "发送成功",
		"fileAddr": realPicFormat,
	})
}

//websocket登陆后建立消息推送长连接(message publisher)
func ReceiveMsgWebsocket(c *gin.Context) {
	//获取发送者uid和被发送者uid
	token := module.WsGetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	//username = "3"
	//升级websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//将获取的参数放进这个数组
		Subprotocols: []string{token},
	}).Upgrade(c.Writer, c.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户实例
	recClient := &models.ClientRecMsg{
		ID:     stuNum,
		Socket: conn,
		Send:   make(chan []byte),
	}
	models.AddMsgPusher(stuNum, recClient)
	go recClient.PushMsg()
	go recClient.CheckOnline()
}

//用户获取消息列表
func GetMsgList(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "未登录",
		})
		return
	}
	msgList := models.GetMsgList(stuNum)
	//fmt.Println(msgList)
	finalMsgList := models.QueryEveryLastMsg(msgList)
	finalMsgList = MsgListSortByStartTime(finalMsgList)
	finalData := []interface{}{}
	for _, msg := range finalMsgList {
		msg.ToUserHeadPic = "http://" + c.Request.Host + "/uploadPic/headPic/" + msg.ToUserHeadPic
		if msg.FromUsername == stuNum && msg.Read == 0 {
			msg.Read = 1
		}
		if msg.Read == 0 {
			finalData = util.SliAddFromHead(finalData, msg)
		} else {
			finalData = append(finalData, msg)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "消息列表获取成功",
		"data":   finalData,
	})
}

//创建trainer的userid
func createID(uid, toUid string) string {
	return uid + "->" + toUid // 1 -> 2
}

func MsgListSortByStartTime(qelm []models.QELM) []models.QELM {
	sort.Slice(qelm, func(i, j int) bool { // desc
		return qelm[i].Start_time > qelm[j].Start_time
	})
	return qelm
}
