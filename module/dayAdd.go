package module

import (
	"strconv"
)

func DayAdd(date string) (date2 string) {
	ThiryOneDaySli := []string{"01", "03", "05", "07", "08", "10", "12"}
	dayNum := date[8:]
	datNumInt, _ := strconv.Atoi(dayNum)
	newDayNum := datNumInt + 1
	if newDayNum < 10 {
		newDayNumString := strconv.Itoa(newDayNum)
		date2 = date[:8] + "0" + newDayNumString
	}else {
		newDayNumString := strconv.Itoa(newDayNum)
		date2 = date[:8] + newDayNumString
	}
	monthNum := date[5:7]
	if dayNum == "30" {
		flag := 0
		for _, v := range ThiryOneDaySli {
			if monthNum == v {
				flag = 1
			}
		}
		if flag == 0 {
			MonthNum, _ := strconv.Atoi(date[5:7])
			newMonthNum := strconv.Itoa(MonthNum + 1)
			return date[:5] + newMonthNum + "-01"
		}
	}
	if monthNum == "12" && dayNum == "31" {
		yearNum := date[:4]
		yearNumInt, _ := strconv.Atoi(yearNum)
		newYearNum := strconv.Itoa(yearNumInt + 1)
		return newYearNum + "-01-01"
	}
	if dayNum == "31" {
		MonthNum, _ := strconv.Atoi(date[5:7])
		newMonthNum := strconv.Itoa(MonthNum + 1)
		return date[:5] + newMonthNum + "-01"
	}
	return date2
}
