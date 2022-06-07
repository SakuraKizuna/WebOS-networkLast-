package models

import (
	"drop_os_back/dao"
)

//定义获取条数的model
type CountNum struct {
	Total_num int `json:"total_num"`
}

//times 获取用户对应的总条数 times表适用
func GetTimesTotalNum(username string) int {
	var CountNum CountNum
	sql := "SELECT COUNT(*) as total_num FROM times WHERE username = '" + username + "'"
	dao.DB.Raw(sql).Scan(&CountNum)
	//fmt.Println(sql)
	//fmt.Println(CountNum.Total_num)
	return CountNum.Total_num
}

//discussions 获取所有条数
func GetDiscussionsTotalNum(classification string) int {
	var CountNum CountNum
	sql := ""
	if classification == "全部" {
		sql = "SELECT COUNT(*) as total_num FROM `discussions`  WHERE level = 0 "
	} else {
		sql = "SELECT COUNT(*) as total_num FROM `discussions`  WHERE level = 0 and classification = '" + classification + "'"
	}
	dao.DB.Raw(sql).Scan(&CountNum)
	return CountNum.Total_num
}

//discussions 获取某个用户的所有文章条数
func GetPersonalDiscussionsTotalNum(username string) int {
	var CountNum CountNum
	sql := "SELECT COUNT(*) as total_num FROM `discussions`  WHERE level = 0 and username ='"+ username + "'"
	dao.DB.Raw(sql).Scan(&CountNum)
	return CountNum.Total_num
}

//获取评论总数
func GetDiscussionsAreaTotalNum(artId string) int {
	var CountNum CountNum
	sql := "SELECT COUNT(*) as total_num FROM discussions WHERE belong_id = " + artId + " and level = 1"
	dao.DB.Raw(sql).Scan(&CountNum)
	return CountNum.Total_num
}
