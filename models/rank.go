package models

import (
	"drop_os_back/dao"
	"strconv"
)

type Rank struct {
	Userid   int    `json:"userid" gorm:"primary_key"`
	Username string `json:"username"`
	Weektime string `json:"weektime"`
	Grade    string `json:"grade"`
	Realname string `json:"realname"`
}

func WeekTimeAdd(username, daytime string) (err error) {
	userRank := Rank{}
	err1 := dao.DB.Where("username=?", username).First(&userRank).Error
	if err1 != nil {
		userRank.Username = username
		userInfo := User{}
		dao.DB.Where("username=?", username).First(&userInfo)
		userRank.Grade = strconv.Itoa(userInfo.Grade)
		userRank.Realname = userInfo.Realname
		userRank.Weektime = "0"
	}
	weektime := userRank.Weektime
	weektimeInt, _ := strconv.Atoi(weektime)
	daytimeInt, _ := strconv.Atoi(daytime)
	NewWeektime := weektimeInt + daytimeInt
	userRank.Weektime = strconv.Itoa(NewWeektime)
	err2 := dao.DB.Save(&userRank).Error
	if err2 != nil {
		return err2
	}
	return nil
}

func GetWeekTime(username string) (weekTime string) {
	WeekTimeInfo := Rank{}
	dao.DB.Where("username=?", username).First(&WeekTimeInfo)
	return WeekTimeInfo.Weektime
}

func GetWeekRank(username, grade string) (rank int,weektime string) {
	var weekTimeInfo []*Rank
	userWeekTime := Rank{}
	dao.DB.Where("grade=?", grade).Find(&weekTimeInfo)
	dao.DB.Where("username=?", username).First(&userWeekTime)
	userTime, _ := strconv.Atoi(userWeekTime.Weektime)
	rank = 1
	for _, v := range weekTimeInfo {
		userTime2, _ := strconv.Atoi(v.Weektime)
		//fmt.Println(userTime2)
		if userTime2 > userTime {
			rank++
		}
	}
	return rank,userWeekTime.Weektime
}

func ResetTime() {
	dao.DB.Model(Rank{}).Update("weektime", "0")
}
