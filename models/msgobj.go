package models

import "drop_os_back/dao"

type Msgobj struct {
	Id           int    `json:"id" gorm:"primary_key"`
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
}

//信息列表
func GetMsgList(fromUsername string) []Msgobj {
	var msgObjs []Msgobj
	dao.DB.Where("from_username=?", fromUsername).Order("id desc").Find(&msgObjs)
	return msgObjs
}

//删除指定id的记录
func DeleteMsg(msgId int, username string) error {
	var msgObj Msgobj
	err := dao.DB.Where("id=?", msgId).Where("from_username=?", username).Find(&msgObj).Error
	if err != nil {
		return err
	}
	dao.DB.Delete(&msgObj)
	return nil
}
