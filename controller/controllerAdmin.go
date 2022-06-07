package controller

import (
	"drop_os_back/dao"
	"drop_os_back/models"
	"drop_os_back/module"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func AdminLogin(c *gin.Context) {
	//把请求拿出来
	var administrator models.Administrator
	c.BindJSON(&administrator)
	//fmt.Println(administrator)
	//fmt.Println(administrator.Adminname)
	//models.Test()

	admin, err := models.GetAdminInfo(administrator.Adminname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "管理员账号不存在",
		})
		return
	}

	if admin.Password == administrator.Password {
		token := dao.RandString(50)
		dao.Redis.Set(token, administrator.Adminname, 14400*time.Second)
		adminLevel, _ := module.TokenToLevelAndBelong(token)
		data := module.GetDashBoard(administrator.Adminname, token, adminLevel)
		c.JSON(http.StatusOK, gin.H{
			"data":   data,
			"status": 200,
			"msg":    "登录成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "密码错误",
		})
	}
}

func AdminLogout(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	token := json["token"].(string)

	dao.Redis.Del(token)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "登出成功",
	})
}

func QueryStudent(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	pageNum := json["pageNum"]

	if json["token"] == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 208,
			"msg":    "token为空",
		})
	}
	adminLevel, adminBelong := module.TokenToLevelAndBelong(json["token"])
	if adminLevel == -1 {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
		return
	}
	//fmt.Println(adminLevel)
	//fmt.Printf("%T\n", adminBelong)
	uL, err := models.GetUserLists(adminLevel, adminBelong)
	//fmt.Println(uL)
	if err != nil {
		fmt.Println("fucking err:", err)
	}
	data := module.MakeSli(pageNum, uL)
	c.JSON(http.StatusOK, gin.H{
		"data":   data,
		"msg":    "查询成功",
		"status": 200,
		"total":  len(uL),
	})
}

func QueryStudentTime(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	if json["token"] == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
	}
	token := json["token"]
	pageNum := json["pageNum"]
	adminLevel, adminBelong := module.TokenToLevelAndBelong(token)
	//fmt.Printf("%T\n", adminLevel)
	stuSli, _ := models.GetUserSli(adminLevel, adminBelong)
	temp_intf := []interface{}{}
	//wg.Add(len(stuSli))
	for _, stuNum := range stuSli {
		//fmt.Println(stuNum)
		//go func(adminLevel int, stuNum string) {
		stuInfo, _ := models.GetUserSingle(stuNum)
		//fmt.Println(stuInfo)
		stuDist, _ := models.GetTimeLast(stuNum)
		//fmt.Println(stuDist)
		stuDist["name"] = stuInfo["name"]
		stuDist["sex"] = stuInfo["sex"]
		temp_intf = append(temp_intf, stuDist)
		//wg.Done()
		//}(adminLevel, stuNum)
	}
	//wg.Wait()
	data := module.MakeSli(pageNum, temp_intf)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "数据获取成功",
		"data":   data,
		"total":  len(temp_intf),
	})
}

func QueryUnusual(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	//fmt.Println(json)
	token := json["token"]
	pageNum := json["pageNum"]
	// TODO 508
	//fmt.Println(token)
	if token == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
	}
	adminLevel, _ := module.TokenToLevelAndBelong(token)
	userMap, _ := models.GetMapStuNumName(adminLevel)
	unTimeData, err := models.GetTimeListTime(adminLevel)
	for _, b := range unTimeData {
		//fmt.Printf("%T\n", b)
		stuNum := b.(map[string]interface{})["schoolNumber"]
		stuNum2 := stuNum.(string)
		b.(map[string]interface{})["name"] = userMap[stuNum2]
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    err,
		})
	}
	data := module.MakeSli(pageNum, unTimeData)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "数据获取成功",
		"data":   data,
		"total":  len(unTimeData),
	})
}

func QueryStuTimeSingle(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	fmt.Println(json["token"])
	schoolNum := json["schoolNum"].(string)
	pageNum := json["pageNum"]
	stuData, _ := models.GetUserSingle(schoolNum)
	SingleTimeData, err := models.GetTimeListSingle(schoolNum)
	for _, items := range SingleTimeData {
		stuTimeData := items.(map[string]interface{})
		//fmt.Println(stuTimeData)
		stuTimeData["name"] = stuData["name"]
		stuTimeData["sex"] = stuData["sex"]
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    err,
		})
	}
	data := module.MakeSli(pageNum, SingleTimeData)
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"msg":      "数据获取成功",
		"data":     data,
		"totalNum": len(SingleTimeData),
	})
}

func TestAPI(c *gin.Context) {
	//fmt.Println(c.Request.Host)
	//a := models.QueryPersonalInfo("20061123")
	//userInfo := models.QueryPersonalInfo("20062111")
	//headPic := userInfo["head_pic"].(string)
	//headPic2 := "http://" + c.Request.Host + "/uploadPic/headPic/" + headPic
	//userInfo["head_pic"] = headPic2
	//fmt.Println(headPic2)
	//a := models.QueryPersonalInfo("20062112")
	//fmt.Println(a["head_pic"] == "")
	//userInfo := models.QueryPersonalInfo("20062111")
	//lastPic := userInfo["head_pic"].(string)
	//if lastPic != ""{
	//	os.Remove("./uploadPic/headPic/"+lastPic)
	//}
	//a := models.QueryPersonalInfo("200")
	//fmt.Println(a)
	//fmt.Println(a["username"] == "")
	//_ = models.AddTodo("20062111", "asdasd")
	//likeInfo,err := models.QueryLikeInfo(1,"20062111")
	//fmt.Println(likeInfo)
	//fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "测试成功，小老弟！！！",
	})
}

//用原生语句重写日期相加
func QueryAllTime(c *gin.Context) {
	var wg sync.WaitGroup
	json := make(map[string]interface{})
	c.BindJSON(&json)
	startDate := json["start_date"].(string)
	endDate := json["end_date"].(string)
	token := json["token"]
	stuNameSli := []interface{}{}
	stuTimeSli := []interface{}{}
	// TODO 508
	adminLevel, adminBelong := module.TokenToLevelAndBelong(token)
	StuNameAName, _ := models.GetStuNumAName(adminLevel)
	stuSli, _ := models.GetUserSli(adminLevel, adminBelong)
	wg.Add(len(stuSli))
	//fmt.Println(stuSli)
	for _, stuNum := range stuSli {
		go func(stuNum, startDate, endDate string) {
			daySum := models.SumDaytimeBYSDAndED(stuNum, startDate, endDate)
			stuName := StuNameAName[stuNum]
			stuNameSli = append(stuNameSli, stuName)
			allTime2 := module.MinutesToH(daySum)
			stuTimeSli = append(stuTimeSli, allTime2)
			wg.Done()
		}(stuNum, startDate, endDate)
	}
	wg.Wait()
	data := []interface{}{stuNameSli, stuTimeSli}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   data,
		"msg":    "数据获取成功",
	})
}

//
//func QueryAllTime2(c *gin.Context) {
//	var wg sync.WaitGroup
//	json := make(map[string]interface{})
//	c.BindJSON(&json)
//	fmt.Println(json)
//	startDate := json["start_date"].(string)
//	endDate := json["end_date"].(string)
//	fmt.Println(startDate)
//	fmt.Printf("%T\n", startDate)
//	token := json["token"]
//	stuNameSli := []interface{}{}
//	stuTimeSli := []interface{}{}
//	adminLevel := module.TokenToLevel(token)
//	StuNameAName, _ := models.GetStuNumAName(adminLevel)
//	//fmt.Println(StuNameAName)
//	stuSli, _ := models.GetUserSli(adminLevel)
//	wg.Add(len(stuSli))
//	fmt.Println(stuSli)
//	for _, stuNum := range stuSli {
//		go func(stuNum, startDate, endDate string) {
//			allTime := 0
//			date2 := startDate
//			if startDate == endDate {
//				//fmt.Println("the way 1")
//				allDayTime, _ := models.SumDaytimeSingle(stuNum, startDate)
//				allTime = allTime + allDayTime
//			} else {
//				//fmt.Println("the way 2")
//				for {
//					//fmt.Println(date2, endDate)
//					allDayTime, _ := models.SumDaytimeSingle(stuNum, date2)
//					//fmt.Println("the way 2,alldaytime:", allDayTime)
//					allTime = allTime + allDayTime
//					if date2 == endDate {
//						break
//					}
//					date2 = module.DayAdd(date2)
//				}
//			}
//			stuName := StuNameAName[stuNum]
//			stuNameSli = append(stuNameSli, stuName)
//			allTime2 := module.MinutesToH(allTime)
//			stuTimeSli = append(stuTimeSli, allTime2)
//			wg.Done()
//		}(stuNum, startDate, endDate)
//	}
//	wg.Wait()
//	data := []interface{}{stuNameSli, stuTimeSli}
//	c.JSON(http.StatusOK, gin.H{
//		"status": 200,
//		"data":   data,
//		"msg":    "数据获取成功",
//	})
//}

func EndSign(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	time2 := json["time"]
	fmt.Println(time2)
	username := json["username"].(string)
	endTime := module.TimeStrToTime(time2.(string))
	fmt.Println(endTime)
	err, daytime := models.EndSign(username, endTime)
	err2 := models.WeekTimeAdd(username, daytime)
	if err != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "结束签到成功",
	})
}

func SignSupply(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	username := json["username"].(string)
	start_time := json["start_time"].(string)
	end_time := json["end_time"].(string)
	date := start_time[:10]
	startTime := module.TimeStrToTime(start_time)
	endTime := module.TimeStrToTime(end_time)
	daytime := models.SignSupply(startTime, endTime, date, username)
	err := models.WeekTimeAdd(username, daytime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "补签成功",
	})
}

//定义<向可疑数据所属的用户发送提醒邮件>API
func SendUnusualEmail(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	stu_data := json["stu_data"]
	stuData := stu_data.([]interface{})
	emailSli, _ := models.GetEmail(stuData)
	fmt.Println(emailSli)
	c.JSON(http.StatusOK, gin.H{
		"status": 202,
		"msg":    "该API未完成，敬清期待",
	})
}

func ReplyLab(c *gin.Context) {
	ApplyList, _ := models.GetApplyLists()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   ApplyList,
		"msg":    "获取成功",
	})
}

func ApplyOK(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	belong := json["belong"].(string)
	username := json["schoolNumber"].(string)
	err := models.ChangeStatus(username, belong)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "操作成功",
	})
}

func ShowAdministrators(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	token := json["token"]
	if token == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
	}
	adInfo, _ := models.GetAdList()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   adInfo,
		"msg":    "获取成功",
	})

}

func AddAdministrator(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	belong := json["belong"].(string)
	level := json["level"].(string)
	password := json["password"].(string)
	username := json["username"].(string)
	err := models.AddAdministrator(belong, level, password, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "管理员添加成功",
	})
}

//定义<查询学生每日详细情况	>API
func QueryStuDateTime(c *gin.Context) {
	json := make(map[string]interface{})
	_ = c.BindJSON(&json)
	timeNum := 0
	timeDate := ""
	startDate := json["start_date"].(string)
	endDate := json["end_date"].(string)
	username := json["username"].(string)
	stuDateTimeSli := []interface{}{}
	timeSli := models.GetTimeSli2(username, startDate, endDate)
	for _, b := range timeSli {
		dataTemp := map[string]interface{}{}
		dataTemp["name"] = "签到"
		timeSli := []interface{}{}
		timeSli = append(timeSli, timeNum)
		if timeDate != b.Date {
			timeNum++
			timeDate = b.Date
		}
		timeSli = append(timeSli, b.Begantime.Format("2006-01-02 15:04:05"))
		timeSli = append(timeSli, b.Endtime.Format("2006-01-02 15:04:05"))
		dataTemp["value"] = timeSli
		stuDateTimeSli = append(stuDateTimeSli, dataTemp)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   stuDateTimeSli,
		"msg":    "查询成功",
	})
}

//定义<重置密码&软删除用户>API
func RemakeORDelete(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	stuNum := json["schoolNumber"].(string)
	behave := json["activity"].(string)
	err := models.ResetPasswordORDelete(stuNum, behave)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "密码重置成功，重置密码为123456",
	})

}

//定义删除<用户->详细>时间
func DeleteSignData(c *gin.Context) {
	json := map[string]interface{}{}
	c.BindJSON(&json)
	beginTime := json["begantime"].(string)
	stuNum := json["schoolNumber"].(string)
	BeginTime := module.TimeStrToTime(beginTime)
	err := models.DeleteByBeginTime(stuNum, BeginTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "服务器错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "操作成功",
	})
}

//升级
func UpIdentity(c *gin.Context) {
	json := map[string]interface{}{}
	c.BindJSON(&json)
	username := json["username"].(string)
	models.UpUserRole(username)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "操作成功",
	})
}
