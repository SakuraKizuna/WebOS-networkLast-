package timeTask

import (
	"drop_os_back/models"
	"time"
)

var MONDAYDATE string

//当程序开启时，设置本周周一日期
func SetMondayDateStartProcess() {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekMonday := weekStartDate.Format("2006-01-02")

	MONDAYDATE = weekMonday
}

//每周重置周一日期
func ResetMonday() {
	now := time.Now()
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	MONDAYDATE = weekStartDate.Format("2006-01-02")
}

//API获得周一日期
func GetMondayDate() (MondayDate string) {
	return MONDAYDATE
}

//重置weekTime
func ResetWeekTime(){
	models.ResetTime()
}

