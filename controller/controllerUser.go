package controller

import (
	"drop_os_back/dao"
	"drop_os_back/models"
	"drop_os_back/module"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
	//"go_service/config"
	//errLog "go_service/log"
	//"go_service/pkg"
)

//定义用户登录
func UserLogin(c *gin.Context) {
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	fmt.Println(json)
	if json["username"] == nil || json["password"] == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 222,
			"msg":    "服务器参数接收不完整",
		})
		return
	}
	username := json["username"].(string)
	password := json["password"].(string)
	//fmt.Println(username, password)
	statusMsg := models.UserLogin(username, password)
	if statusMsg == "noRecord" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "该账号信息未注册，去注册一下吧",
		})
		return
	}
	if statusMsg == "noAgree" {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "该账号信息未被管理员同意，请联系管理员",
		})
		return
	}
	if statusMsg == "passwordErr" {
		c.JSON(http.StatusOK, gin.H{
			"status": 203,
			"msg":    "密码错误，请重新输入密码",
		})
		return
	}
	token := dao.RandString(50)
	dao.Redis.Set(token, username, 14400*time.Second)
	//fmt.Println(token)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"token":  token,
		"msg":    "登录成功",
	})
}

//定义结束签到
func UserEndSign(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	if token == "none" || token == "null" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	json := map[string]interface{}{}
	c.BindJSON(&json)
	content := json["content"].(string)
	stuNum := dao.Redis.Get(token).Val()
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	module.ResetTokenTime(token, stuNum)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	EndTime := module.TimeStrToTime(nowTime)
	_, daytime := models.UserEndSign(stuNum, content, EndTime)
	if daytime == "normal" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "结束签到失败，并不处于签到中状态",
		})
		return
	}
	err2 := models.WeekTimeAdd(stuNum, daytime)
	err3 := models.TotalTimeAdd(stuNum, daytime)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 206,
			"msg":    "服务器错误,err2",
		})
		return
	}
	if err3 != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 206,
			"msg":    "服务器错误,err3",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "结束签到成功",
	})
}

//定义开始签到
func UserStartSign(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	if token == "none" || token == "null" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	stuNum := dao.Redis.Get(token).Val()
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	module.ResetTokenTime(token, stuNum)
	status, _ := models.JudgeLast(stuNum)
	if status == "abnormal" {
		c.JSON(http.StatusOK, gin.H{
			"status": 206,
			"msg":    "开始签到失败，该用户签到中",
		})
		return
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	date := nowTime[:10]
	StartTime := module.TimeStrToTime(nowTime)
	models.StartSign(stuNum, date, StartTime)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "开始签到成功",
	})
}

//忘记密码（重置密码）
func ResetPassword(c *gin.Context) {
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	email := json["email"].(string)
	newPassword := json["newPassword"].(string)
	username := json["username"].(string)
	code := json["code"].(string)
	codex := dao.Redis.Get(email).Val()
	if codex == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "验证码失效，请重新发送",
		})
		return
	}
	if code != codex {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "验证码错误，请重新发送",
		})
		return
	}
	userInfo := models.QueryPersonalInfo(username)
	if userInfo["email"].(string) != email {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "请输入该用户绑定的邮箱:" + userInfo["email"].(string),
		})
		return
	}
	models.ResetPassword(email, newPassword, username)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "密码重置成功",
	})
}

//定义发送注册邮箱
func SendEmail(c *gin.Context) {
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	if json == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 204,
			"msg":    "未能正确接收json信息",
		})
		return
	}
	email := json["email"].(string)
	fmt.Println(email)
	code := module.GenValidateCode(4)
	err := module.SendEmail(email, code)
	dao.Redis.Set(email, code, 600*time.Second)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误，或邮箱不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "验证码发送成功",
	})
}

//定义用户注册
func UserRegister(c *gin.Context) {
	json := map[string]interface{}{}
	c.BindJSON(&json)
	//fmt.Println(json)
	username := json["username"].(string)
	password := json["password"].(string)
	password2 := json["password2"].(string)
	grade := json["grade"].(string)
	email := json["email"].(string)
	code := json["code"].(string)
	realname := json["realname"].(string)
	if password2 != password {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "两次输入的密码不一致，请重新输入",
		})
		return
	}
	codex := dao.Redis.Get(email).Val()
	if code != codex {
		c.JSON(http.StatusOK, gin.H{
			"status": 203,
			"msg":    "验证码失效或错误，请重新发送验证码",
		})
		return
	}
	userInfo := models.QueryPersonalInfo(username)
	fmt.Println(userInfo["username"].(string))
	if userInfo["username"].(string) != "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "该用户名已存在",
		})
		return
	}
	models.UserRegister(username, password, email, grade, realname)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "注册成功，待管理员同意后方可登录",
	})
}

//获取时间信息（总时间，本周时间，本周时间排名，今日时长）
func UserGetTimeInfo(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	if token == "none" {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	//token := "PDXWMZLYNEXXNDPOOWGBPOFNTEOLBRUCWIALBHCFRTQFQHCKIF"
	stuNum := dao.Redis.Get(token).Val()
	grade := models.QueryGrade(stuNum)
	gradeStr := strconv.Itoa(grade)
	dataInfo := map[string]interface{}{}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	date := nowTime[:10]
	totalTime := models.GetPersonalTotal(stuNum)
	//weekTime := models.GetWeekTime(stuNum)
	allDayTime := models.GetAllDayTime(stuNum, date)
	weekRank, weekTime := models.GetWeekRank(stuNum, gradeStr)
	dataInfo["totalTime"] = totalTime
	dataInfo["weekTime"] = weekTime
	dataInfo["allDayTime"] = allDayTime
	dataInfo["weekRank"] = weekRank
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   dataInfo,
		"msg":    "获取成功",
	})
}

//发表帖子(帖子对应的图片api)
// TODO 先请求文章部分api再请求图片部分api
func UserPublishCommentPic(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	if token == "none" {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	//获取文件头
	file, err := c.FormFile("uploadPicture")
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	//获取文件名
	stuNum := dao.Redis.Get(token).Val()
	//stuNum := "20062111"
	fileName := file.Filename
	fmt.Println("文件名：", fileName)
	//保存文件到服务器本地
	//SaveUploadedFile(文件头，保存路径)
	//生成独一份的uuid
	fileUuid := module.GetUUID()
	filePath := "./uploadPic/" + fileUuid + ".jpg"
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	artId := dao.Redis.Get(stuNum).Val()
	artIdInt, _ := strconv.Atoi(artId)
	err2 := models.PublishArtPic(artIdInt, filePath)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "图片数据上传失败，服务器错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "图片上传成功",
	})

}

//发表帖子(帖子对应的文章api)
func UserPublishCommentArt(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	//stuNum := "20062111"
	//获取文件头
	file, err := c.FormFile("uploadPicture")
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	if stuNum == "游客" {
		c.JSON(http.StatusOK, gin.H{
			"status": 210,
			"msg":    "此账号为游客账号，没有发帖权限",
		})
		return
	}
	//fmt.Println(file.Filename)
	fileFormat := file.Filename[len(file.Filename)-4:]
	if file.Filename[len(file.Filename)-3:] == "jpg" || file.Filename[len(file.Filename)-3:] == "png" {
		//fmt.Println(fileFormat)
		comment, _ := c.GetPostForm("comment")
		classification, _ := c.GetPostForm("classification")
		title, _ := c.GetPostForm("title")
		//获取文件名
		//stuNum := "20062111"
		//保存文件到服务器本地
		//生成独一份的uuid
		fileUuid := module.GetUUID()
		filePath := "./uploadPic/artPic/" + fileUuid + "-" + stuNum + fileFormat
		//fmt.Println(filePath)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
			return
		}
		artPicAdd := fileUuid + "-" + stuNum + fileFormat
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		nowTime2 := module.TimeStrToTime(nowTime)
		err2 := models.PublishArt(stuNum, comment, classification, artPicAdd, title, nowTime2)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 202,
				"msg":    "文章发表失败，服务器错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "文章发表成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 208,
		"msg":    "发布失败，只支持jpg或png格式图像",
	})
	return

}

//评论帖子
// TODO 评论后给所属文章的评论数+1如果是评论评论就把评论的评论数+1
func UserDiscuss(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	id := json["id"].(string)
	idInt, _ := strconv.Atoi(id)
	comment := json["comment"].(string)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	nowTime2 := module.TimeStrToTime(nowTime)
	//stuNum := "20062111"
	err := models.PublishComment(idInt, stuNum, comment, nowTime2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "评论失败，服务器错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "评论成功",
	})
}

//获取所有帖子
// TODO discussions表里的pic_address是封面图片，后期新建库专门保存内容的如片
func GetAllArticles(c *gin.Context) {
	//json := map[string]interface{}{}
	//c.BindJSON(&json)
	//pageNum := json["pageNum"].(string)
	classification := c.Param("classification")
	pageNum := c.Param("pageNum")
	pageNumInt, _ := strconv.Atoi(pageNum)
	disInfo := models.GetAllArt(pageNumInt, classification)
	totalNum := models.GetDiscussionsTotalNum(classification)
	// TODO 及时修改地址
	for _, item := range disInfo {
		artMap := item.(map[string]interface{})
		artPic := artMap["pic_address"].(string)
		if artPic != "" {
			artPic2 := "http://" + c.Request.Host + "/uploadPic/artPic/" + artPic
			//artPic2 := "http://5epzud.natappfree.cc" + "/uploadPic/artPic/" + headPic
			artMap["pic_address"] = artPic2
		} else {
			artMap["pic_address"] = "null"
		}
	}
	totalPage := 0
	if totalNum%15 == 0 {
		totalPage = totalNum / 15
	} else {
		totalPage = totalNum/15 + 1
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"dataList": disInfo,
		"total":    totalPage,
		"msg":      "帖子信息获取成功",
	})
}

//获取单个用户的所有文章
func GetPersonalAllArticles(c *gin.Context) {
	pageNum := c.Param("pageNum")
	username := c.Param("username")
	pageNumInt, _ := strconv.Atoi(pageNum)
	disInfo := models.UserGetAllArt(pageNumInt, username)
	totalNum := models.GetPersonalDiscussionsTotalNum(username)
	for _, item := range disInfo {
		artMap := item.(map[string]interface{})
		artPic := artMap["pic_address"].(string)
		if artPic != "" {
			artPic2 := "http://" + c.Request.Host + "/uploadPic/artPic/" + artPic
			//artPic2 := "http://5epzud.natappfree.cc" + "/uploadPic/artPic/" + headPic
			artMap["pic_address"] = artPic2
		} else {
			artMap["pic_address"] = "null"
		}
	}
	totalPage := 0
	if totalNum%15 == 0 {
		totalPage = totalNum / 15
	} else {
		totalPage = totalNum/15 + 1
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"dataList": disInfo,
		"total":    totalPage,
		"msg":      "个人帖子信息获取成功",
	})
}

//获取单个帖子
func GetArticle(c *gin.Context) {
	articleID := c.Param("articleId")
	articleIDInt, _ := strconv.Atoi(articleID)
	disInfo := models.GetArticleIn(articleIDInt)
	picAdd := disInfo["pic_address"].(string)
	if picAdd != "" {
		picAdd = "http://" + c.Request.Host + "/uploadPic/artPic/" + picAdd
	} else {
		picAdd = "null"
	}
	disInfo["pic_address"] = picAdd
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   disInfo,
		"msg":    "文章信息获取成功",
	})
}

//获取评论区   GET
func GetDiscussionArea(c *gin.Context) {
	articleID := c.Param("id")
	pageNum := c.Param("pageNum")
	articleIDInt, _ := strconv.Atoi(articleID)
	pageNumInt, _ := strconv.Atoi(pageNum)
	//dataList := models.MakeDiscussionData(articleIDInt)
	dataList := models.DealFatherDiscussionArea(pageNumInt, articleIDInt)
	totalNum := models.GetDiscussionsAreaTotalNum(articleID)
	totalPage := 0
	if totalNum%10 == 0 {
		totalPage = totalNum / 10
	} else {
		totalPage = totalNum/10 + 1
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"dataList": dataList,
		"total":    totalPage,
		"msg":      "评论区信息获取成功",
	})
}

//获取未读消息数量 *扩展
func GetUnreadMessage(c *gin.Context) {

}

//通知中心 *扩展
func GetNotificationInfo(c *gin.Context) {

}

//TODO 测试上传文件  prepare for my personalWeb and labPhotos
func TestReceiveFile(c *gin.Context) {
	//获取文件头
	file, err := c.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	//获取文件名
	fileName := file.Filename
	fmt.Println("文件名：", fileName)
	//保存文件到服务器本地
	//SaveUploadedFile(文件头，保存路径)
	filePath := "./uploadPic/" + fileName
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "test",
	})
}

func TestHeaders(c *gin.Context) {
	//token := "none"
	//body, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println("---body/--- \r\n " + string(body))
	//fmt.Println("---header/--- \r\n")
	//fmt.Printf("11111:%T\n", c.Request.Header)
	//for k, v := range c.Request.Header {
	//	if k == "Token" {
	//		token = v[0]
	//		fmt.Println(token)
	//	}
	//}
	token := module.GetTokenFromHeader(c.Request.Header)
	fmt.Println(token)
	if token == "none" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "未登录",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":        200,
		"msg":           "token获取成功",
		"receivedToken": token,
	})
}

func GetStatus(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	fmt.Println(token)
	if token == "none" || token == "null" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "未登录",
		})
		return
	}
	stuNum := dao.Redis.Get(token).Val()
	fmt.Println(stuNum)
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	module.ResetTokenTime(token, stuNum)
	status, _ := models.JudgeLast(stuNum)
	if status == "abnormal" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "签到中",
		})
		return
	}
	if status == "none" {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "未签到",
	})
}

func UserGetStudyRec(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	fmt.Println(stuNum)
	pageNum := json["pageNum"].(string)
	pageNumInt, _ := strconv.Atoi(pageNum)
	recordInfo := models.UserGetContent(stuNum, pageNumInt)
	RecData := models.DealRecInfo(recordInfo)
	totalNum := models.GetTimesTotalNum(stuNum)
	totalPage := 0
	if totalNum%10 == 0 {
		totalPage = totalNum / 10
	} else {
		totalPage = totalNum/10 + 1
	}
	//fmt.Println(totalPage)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   RecData,
		"total":  totalPage,
		"msg":    "日志获取成功",
	})
	// TODO 每日学习记录顺序调换 || 超级管理员权限及操作 2022.2.26
}

func UserGetPersonalInfo(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	fmt.Println(token)
	if token == "none" || token == "null" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "未登录",
		})
		return
	}
	stuNum := dao.Redis.Get(token).Val()
	module.ResetTokenTime(token, stuNum)
	fmt.Println(stuNum)
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 205,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	module.ResetTokenTime(token, stuNum)
	userInfo := models.QueryPersonalInfo(stuNum)
	headPic := userInfo["head_pic"].(string)
	// TODO 及时修改地址
	if headPic != "null" {
		headPic2 := "http://" + c.Request.Host + "/uploadPic/headPic/" + headPic
		//headPic2 := "http://5epzud.natappfree.cc" + "/uploadPic/headPic/" + headPic
		userInfo["head_pic"] = headPic2
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   userInfo,
		"msg":    "个人信息获取成功",
	})
}

func GetPersonalInfo(c *gin.Context) {
	username := c.Param("username")
	userInfo := models.QueryPersonalInfo(username)
	headPic := userInfo["head_pic"].(string)
	if headPic != "null" {
		headPic2 := "http://" + c.Request.Host + "/uploadPic/headPic/" + headPic
		//headPic2 := "http://5epzud.natappfree.cc" + "/uploadPic/headPic/" + headPic
		userInfo["head_pic"] = headPic2
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   userInfo,
		"msg":    "个人信息获取成功",
	})
}

func ModifyPersonalInfo(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	realname := json["realname"].(string)
	gender := json["gender"].(string)
	profession := json["profession"].(string)
	email := json["email"].(string)
	nickname := json["nickname"].(string)
	grade := json["grade"].(string)
	err := models.ModifyPersonalInfo(stuNum, realname, gender, profession, email, nickname, grade)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "个人信息提交成功",
	})
}

func UserUploadHeadPic(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	//获取文件头
	file, err := c.FormFile("uploadPicture")
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	fileFormat := file.Filename[len(file.Filename)-4:]
	if file.Filename[len(file.Filename)-3:] == "jpg" || file.Filename[len(file.Filename)-3:] == "png" {
		//获取文件名
		//stuNum := "20062111"
		//保存文件到服务器本地
		//生成独一份的uuid
		fileUuid := module.GetUUID()
		filePath := "./uploadPic/headPic/" + fileUuid + "-" + stuNum + fileFormat
		fmt.Println(filePath)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
			return
		}
		//删除该用户原来的头像
		userInfo := models.QueryPersonalInfo(stuNum)
		lastPic := userInfo["head_pic"].(string)
		if lastPic != "" {
			os.Remove("./uploadPic/headPic/" + lastPic)
		}
		headPicName := fileUuid + "-" + stuNum + fileFormat
		fmt.Println(headPicName)
		err2 := models.ModifyHeadPic(stuNum, headPicName)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 202,
				"msg":    "服务器错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "头像上传成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 208,
		"msg":    "发布失败，不支持上传gif格式图片",
	})
	return

}

func UserGetTodo(c *gin.Context) {
	todoList, err := models.GetTodoList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "服务器错误",
		})
		return
	}
	TodoSli := models.DealTodoList(todoList)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   TodoSli,
		"msg":    "全员待办获取成功",
	})
}

func UserAdminDeleteTodo(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	id := json["id"].(string)
	idInt, _ := strconv.Atoi(id)
	UserRole := models.CheckRole(stuNum)
	if UserRole == "user" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "权限不足，不能删除全员待办",
		})
		return
	}
	if UserRole == "admin" {
		err := models.DeleteTodo(idInt)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 204,
				"msg":    "服务器错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "全员待办删除成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 203,
		"msg":    "服务器错误，不能识别角色",
	})
}

func UserAdminAddTodo(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	content := json["content"].(string)
	UserRole := models.CheckRole(stuNum)
	if UserRole == "user" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "权限不足，不能发表全员待办",
		})
		return
	}
	if UserRole == "admin" {
		err := models.AddTodo(stuNum, content)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 204,
				"msg":    "服务器错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "全员待办发表成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 203,
		"msg":    "服务器错误，不能识别角色",
	})

}

// GET类型
func UserGetThumbsUp(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	//json := map[string]interface{}{}
	//_ = c.BindJSON(&json)
	//id := json["id"].(string)
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	_, err := models.QueryLikeInfo(idInt, stuNum)
	//fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func GetUserHeadPic(c *gin.Context) {
	username := c.Param("username")
	headPic := models.QueryHeadPic(username)
	if headPic != "" {
		headPic = "http://" + c.Request.Host + "/uploadPic/headPic/" + headPic
	} else {
		headPic = "null"
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   headPic,
		"msg":    "头像获取成功",
	})
}

func UserLikeCon(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	id := json["id"].(string)
	beahavior := json["beahavior"].(string)
	fmt.Println(id)
	fmt.Println(beahavior)
	// behavior 只有like和dislike
	idInt, _ := strconv.Atoi(id)
	err := models.ControlLike(idInt, beahavior)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "操作失败，服务器错误",
		})
		return
	}
	if beahavior == "like" {
		err2 := models.AddLikeUserInfo(idInt, stuNum)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 202,
				"msg":    "点赞失败，服务器错误",
			})
			return
		}
	} else {
		errx := models.DeleteLikeInfo(idInt, stuNum)
		fmt.Println(errx)
		if errx != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 209,
				"msg":    "取消点赞失败，服务器错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "取消点赞成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "点赞成功",
	})
}

func GetMoreDis(c *gin.Context) {
	pageNum := c.Param("pageNum")
	artId := c.Param("artId")
	parentId := c.Param("parentId")
	pageNumInt, _ := strconv.Atoi(pageNum)
	dataList := models.GetMoreDis(pageNumInt, artId, parentId)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   dataList,
		"msg":    "分评论获取成功",
	})
}

func GetHeadPicBT(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	headPic := models.QueryHeadPic(stuNum)
	if headPic != "" {
		headPic = "http://" + c.Request.Host + "/uploadPic/headPic/" + headPic
	} else {
		headPic = "null"
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   headPic,
		"msg":    "头像获取成功",
	})
}

func CheckChatStatus(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	json := map[string]interface{}{}
	_ = c.BindJSON(&json)
	toStuNum := json["toStuNum"].(string)
	if stuNum == toStuNum {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "不能与自己发起聊天",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "状态检测成功",
	})
}

func GetContactList(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "未登录",
		})
		return
	}
	ContactList := models.GetContactListIn()
	c.JSON(http.StatusOK, gin.H{
		"status":      200,
		"msg":         "状态检测成功",
		"ContactList": ContactList,
	})
}

func ResetTokenFromMsgList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "token续签成功",
	})
}

//用户删除消息列表的单条聊天
func DeleteOneMsgUser(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	msgId := recJson["msgId"].(string)
	fmt.Println(msgId)
	fmt.Println(stuNum)
	msgIdInt, _ := strconv.Atoi(msgId)
	err := models.DeleteMsg(msgIdInt, stuNum)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "服务器逻辑层错误，该记录不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "删除成功",
	})
}

func DeleteArticle(c *gin.Context) {
	token := module.GetTokenFromHeader(c.Request.Header)
	stuNum := dao.Redis.Get(token).Val()
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	artId := recJson["artId"].(string)
	userRole := models.CheckRole(stuNum)
	artIdInt, _ := strconv.Atoi(artId)
	artInfo := models.GetArticleIn(artIdInt)
	artWriter := artInfo["username"].(string)
	if userRole == "admin" || artWriter == stuNum {
		err := models.DeleteArticle(artIdInt)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 201,
				"msg":    "删除失败，不存在该文章",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "删除成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 202,
		"msg":    "删除失败，身份识别非本人或管理员",
	})
}
