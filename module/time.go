package module

import (
	"fmt"
	"strconv"
	"time"
)

//字符串格式时间转换为时间格式
func TimeStrToTime(timer string) (TrueTime time.Time) {
	local, _:=time.ParseInLocation("2006-01-02 15:04:05", timer, time.Local)
	return local
}

//时间转换为小时（保存两位小数）
func MinutesToH(minutues int) (value float64) {
	a2 := float64(minutues)
	c := a2 / float64(60)
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", c), 64)
	return value
}
