package models

import (
	"drop_os_back/dao"
	"fmt"
	"strconv"
)

// Administrator Model
type Administrator struct {
	ID        int    `json:"id" gorm:"primary_key"`
	Adminname string `json:"adminname"`
	Password  string `json:"password"`
	Level     int    `json:"level"`
	Belong    string `json:"belong"`
}

/*
Administrator CRUD
*/

func GetAdminInfo(adminname string) (administrator *Administrator, err error) {
	administrator = new(Administrator)
	if err = dao.DB.Where("adminname=?", adminname).First(&administrator).Error; err != nil {
		return nil, err
	}
	fmt.Println(administrator)
	return
}

func GetAdList() (adInfo []interface{}, err error) {
	var adList []*Administrator
	err = dao.DB.Find(&adList).Error
	if err != nil {
		return nil, err
	}
	adInfo = []interface{}{}
	for _, b := range adList {
		dist := map[string]interface{}{}
		dist["username"] = b.Adminname
		dist["password"] = b.Password
		dist["level"] = b.Level
		dist["belong"] = b.Belong
		if b.Level == 0 {
			continue
		}
		adInfo = append(adInfo, dist)
	}
	return adInfo, nil
}

func AddAdministrator(belong, level, password, username string) (err error) {
	admin := Administrator{}
	admin.Belong = belong
	admin.Level, _ = strconv.Atoi(level)
	admin.Password = password
	admin.Adminname = username
	err = dao.DB.Create(&admin).Error
	if err != nil {
		return err
	}
	return
}
