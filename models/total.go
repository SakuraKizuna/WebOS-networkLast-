package models

import (
	"drop_os_back/dao"
	"fmt"
	"strconv"
)

type Total struct {
	Id        int    `json:"id" gorm:"primary_key"`
	Username  string `json:"username"`
	Totaltime string `json:"totaltime"`
}

func TotalTimeAdd(username,totaltime string)(err error){
	totalTimeInfo := Total{}
	err1 := dao.DB.Where("username=?",username).Last(&totalTimeInfo).Error
	fmt.Println(err1)
	if err1 != nil{
		totalTimeInfo.Username = username
		totalTimeInfo.Totaltime = "0"
	}
	totalTimeSource := totalTimeInfo.Totaltime
	tTSInt,_ := strconv.Atoi(totalTimeSource)
	totaltimeInt,_ := strconv.Atoi(totaltime)
	total := tTSInt + totaltimeInt
	totalStr := strconv.Itoa(total)
	//totalTimeInfo.Username = username
	//totalTimeInfo.Totaltime = totaltime
	totalTimeInfo.Totaltime = totalStr
	err = dao.DB.Save(&totalTimeInfo).Error
	if err != nil{
		return err
	}
	return nil
}

func GetPersonalTotal(username string)(totalTime string){
	totalInfo := Total{}
	err := dao.DB.Where("username=?",username).Last(&totalInfo).Error
	if err != nil{
		return "0"
	}
	return totalInfo.Totaltime
}







