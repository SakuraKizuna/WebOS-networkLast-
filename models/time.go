package models

import (
	"drop_os_back/dao"
	"fmt"
	"strconv"
	"time"
)

type Time struct {
	Userid    int       `json:"userid" gorm:"primary_key"`
	Username  string    `json:"username"`
	Grade     int       `json:"grade"`
	Begantime time.Time `json:"begantime"`
	Endtime   time.Time `json:"endtime"`
	Content   string    `json:"content"`
	Daytime   string    `json:"daytime"`
	Date      string    `json:"date"`
	Belong    string    `json:"belong"`
}

func reverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//query_student_time 获取<学员考勤>信息
func GetTimeList(adminLevel int) (b []interface{}, err error) {
	var timeList []*Time
	if adminLevel == 0 {
		err = dao.DB.Find(&timeList).Error
		if err != nil {
			fmt.Println("mysql find err:", err)
			return nil, err
		}
	} else {
		err = dao.DB.Where("grade=?", adminLevel).Find(&timeList).Error
		if err != nil {
			fmt.Println("mysql find err:", err)
			return nil, err
		}
	}
	temp_intf := []interface{}{}
	for _, b := range timeList {
		dist := map[string]interface{}{}
		dist["schoolNumber"] = b.Username
		dist["startTime"] = b.Begantime.Format("2006-01-02 15:04:05")
		dist["endTime"] = b.Endtime.Format("2006-01-02 15:04:05")
		dist["content"] = b.Content
		dist["name"] = b.Username
		dist["sex"] = "sex"
		//fmt.Println(dist)
		temp_intf = append(temp_intf, dist)
	}
	temp_intf2 := reverse(temp_intf)
	return temp_intf2, nil
}

//获取最新一条的学员信息
func GetTimeLast(username string) (b map[string]interface{}, err error) {
	var timeList []*Time
	err = dao.DB.Where("username=?", username).Last(&timeList).Error
	if err != nil {
		fmt.Println("GetTimeLast failed,err:", err)
		return nil, err
	}
	dist := map[string]interface{}{}
	for _, b := range timeList {
		dist["schoolNumber"] = b.Username
		dist["startTime"] = b.Begantime.Format("2006-01-02 15:04:05")
		dist["endTime"] = b.Endtime.Format("2006-01-02 15:04:05")
		dist["content"] = b.Content
		dist["name"] = b.Username
		dist["sex"] = "sex"
	}
	return dist, nil
}

// dev/query_unusual 获取<可疑数据>信息
func GetTimeListTime(adminLevel int) (b []interface{}, err error) {
	var timeList []*Time
	if adminLevel == 0 {
		err = dao.DB.Find(&timeList).Error
	} else {
		err = dao.DB.Where("grade=?", adminLevel).Find(&timeList).Error
	}
	temp_intf := []interface{}{}
	for _, b := range timeList {
		dist := map[string]interface{}{}
		dist["schoolNumber"] = b.Username
		dist["startTime"] = b.Begantime.Format("2006-01-02 15:04:05")
		dist["endTime"] = b.Endtime.Format("2006-01-02 15:04:05")
		dist["dayTime"] = b.Daytime
		dist["name"] = b.Username
		//fmt.Println(dist)
		daytime, _ := strconv.Atoi(b.Daytime)
		if daytime >= 300 {
			temp_intf = append(temp_intf, dist)
		}
	}
	temp_intf2 := reverse(temp_intf)
	return temp_intf2, nil
}

//获取个人详细数据 dev/query_student_time_single
func GetTimeListSingle(schoolNum string) (b []interface{}, err error) {
	var timeList []*Time
	err = dao.DB.Where("username=?", schoolNum).Find(&timeList).Error
	temp_intf := []interface{}{}
	for _, b := range timeList {
		dist := map[string]interface{}{}
		dist["schoolNumber"] = b.Username
		dist["startTime"] = b.Begantime.Format("2006-01-02 15:04:05")
		dist["endTime"] = b.Endtime.Format("2006-01-02 15:04:05")
		dist["content"] = b.Content
		dist["name"] = b.Username
		dist["sex"] = "sex"
		//fmt.Println(dist)
		temp_intf = append(temp_intf, dist)
	}
	temp_intf2 := reverse(temp_intf)
	return temp_intf2, nil
}

func SumDaytimeSingle(username, date string) (allDaytime int, err error) {
	var timeList []*Time
	err = dao.DB.Where("username=?", username).Where("date=?", date).Find(&timeList).Error
	if err != nil {
		fmt.Println("SumDaytime failed,err:", err)
		return -1, err
	}
	allDaytime = 0
	//fmt.Println(timeList)
	for _, k := range timeList {
		daytime, _ := strconv.Atoi(k.Daytime)
		//fmt.Printf("%T\n",daytime)
		allDaytime = allDaytime + daytime
	}
	//fmt.Println(timeList)
	return allDaytime, err
}

func EndSign(username string, endTime time.Time) (err error, dayTime string) {
	newTime := Time{}
	var startTime time.Time
	err1 := dao.DB.Where("username=?", username).Last(&newTime).Error
	if err1 != nil {
		fmt.Println("EndSign err1 failed,err:", err1)
		return err1, "error"
	}
	if newTime.Daytime != "0" {
		return nil, "normal"
	}
	startTime = newTime.Begantime
	daytime := int(endTime.Sub(startTime).Minutes())
	if daytime == 0 {
		daytime = 1
	}
	newTime.Endtime = endTime
	newTime.Daytime = strconv.Itoa(daytime)
	err = dao.DB.Save(&newTime).Error
	if err != nil {
		fmt.Println("EndSign saves failed,err:", err)
		return err, "error"
	}
	return nil, newTime.Daytime
}

//用户
func UserEndSign(username, content string, endTime time.Time) (err error, dayTime string) {
	newTime := Time{}
	var startTime time.Time
	err1 := dao.DB.Where("username=?", username).Last(&newTime).Error
	if err1 != nil {
		fmt.Println("EndSign err1 failed,err:", err1)
		return err1, "error"
	}
	if newTime.Daytime != "0" {
		return nil, "normal"
	}
	startTime = newTime.Begantime
	daytime := int(endTime.Sub(startTime).Minutes())
	if daytime == 0 {
		daytime = 1
	}
	newTime.Content = content
	newTime.Endtime = endTime
	newTime.Daytime = strconv.Itoa(daytime)
	err = dao.DB.Save(&newTime).Error
	if err != nil {
		fmt.Println("EndSign saves failed,err:", err)
		return err, "error"
	}
	return nil, newTime.Daytime
}

//定义补签
func SignSupply(startTime, endTime time.Time, date, username string) (dayTime string) {
	timeList := Time{}
	daytime := int(endTime.Sub(startTime).Minutes())
	timeList.Date = date
	timeList.Grade, _ = strconv.Atoi(username[:2])
	timeList.Username = username
	timeList.Daytime = strconv.Itoa(daytime)
	timeList.Begantime = startTime
	timeList.Endtime = endTime
	dao.DB.Create(&timeList)
	return timeList.Daytime
}

//定义(倒序)查询学生总时长
func GetTimeSli(username string) (timeSli2 []*Time, err error) {
	var timeSli []*Time
	err = dao.DB.Where("username=?", username).Order("userid desc").Find(&timeSli).Error
	if err != nil {
		return nil, err
	}
	return timeSli, nil
}

//定义查询学生总时长(原生sql语句版本)
func GetTimeSli2(username, startDate, endDate string) (timeSli2 []*Time) {
	var timeSli []*Time
	startDate2 := startDate
	endDate2 := endDate
	sql := "select * from times where username = '" + username + "' and date BETWEEN '" + startDate2 + "' AND '" + endDate2 + "'"
	//fmt.Println(sql)
	dao.DB.Raw(sql).Scan(&timeSli)
	return timeSli
}

//定义查询学生每日总时长（用户端）
func GetAllDayTime(username, date string) (Alldaytime int) {
	var timeList []*Time
	sql := "select * from times where username = '" + username + "' AND date = '" + date + "'"
	dao.DB.Raw(sql).Scan(&timeList)
	allDaytime := 0
	for _, v := range timeList {
		daytimeInt, _ := strconv.Atoi(v.Daytime)
		allDaytime += daytimeInt
	}
	return allDaytime
}

//定义通过原生语句拿取时间进行计算
func SumDaytimeBYSDAndED(username, startDate, endDate string) (sumTime int) {
	var timeSli []*Time
	Sum := 0
	startDate2 := startDate
	endDate2 := endDate
	sql := "select * from times where username = '" + username + "' and date BETWEEN '" + startDate2 + "' AND '" + endDate2 + "'"
	dao.DB.Raw(sql).Scan(&timeSli)
	for _, items := range timeSli {
		daytime, _ := strconv.Atoi(items.Daytime)
		Sum += daytime
	}
	//fmt.Println(Sum)
	return Sum
}

//定义通过开始时间查询该条记录并进行删除
func DeleteByBeginTime(stuNum string, beginTime time.Time) (err error) {
	timeInfo := Time{}
	err = dao.DB.Where("username=?", stuNum).Where("begantime=?", beginTime).Unscoped().Delete(&timeInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func StartSign(stuNum, date string, startTime time.Time) {
	timeInfo := Time{}
	timeInfo.Username = stuNum
	timeInfo.Begantime = startTime
	timeInfo.Endtime = startTime
	timeInfo.Daytime = "0"
	timeInfo.Date = date
	// TODO 年级问题
	grade := stuNum[0:2]
	gradeInt, _ := strconv.Atoi(grade)
	timeInfo.Grade = gradeInt
	dao.DB.Save(&timeInfo)
}

func JudgeLast(username string) (status string, err error) {
	timeInfo := Time{}
	err = dao.DB.Where("username=?", username).Last(&timeInfo).Error
	if err != nil {
		fmt.Println("GetTimeLast failed,err:", err)
		return "none", err
	}
	if timeInfo.Begantime != timeInfo.Endtime {
		return "normal", nil
	}
	return "abnormal", nil
}

//sql实例
//sql:SELECT * FROM `times`  WHERE  userid >= (SELECT userid FROM `times` LIMIT 100, 1) LIMIT 10;
//获取所有用户的信息（分页）
func GetContent(pageNum int) {
	var timeSli []*Time
	pageNumStr := strconv.Itoa((pageNum - 1) * 10)
	sql := "SELECT * FROM `times` WHERE userid >= (SELECT userid FROM `times` LIMIT" + pageNumStr + ", 1) LIMIT 10;"
	dao.DB.Raw(sql).Scan(&timeSli)
}

//用户获取自己的日志
func UserGetContent(username string, pageNum int) (Info []*Time) {
	var timeSli []*Time
	pageNumStr := strconv.Itoa((pageNum - 1) * 10)
	sql := "SELECT * FROM `times`  WHERE username = '" + username + "' ORDER BY userid desc LIMIT " + pageNumStr + ",10; "
	//sql := "SELECT * FROM `times`  WHERE username = '" + username + "' LIMIT " + pageNumStr + ",10 order by userid desc; "
	fmt.Println(sql)
	dao.DB.Raw(sql).Scan(&timeSli)
	return timeSli
}

//将日志整理成map形式
func DealRecInfo(Info []*Time) (RecSli []interface{}) {
	RecDataSli := []interface{}{}
	for _, v := range Info {
		RecData := map[string]interface{}{}
		RecData["begantime"] = v.Begantime.Format("2006-01-02 15:04:05")
		RecData["endtime"] = v.Endtime.Format("2006-01-02 15:04:05")
		RecData["content"] = v.Content
		RecDataSli = append(RecDataSli, RecData)
	}
	return RecDataSli
}
