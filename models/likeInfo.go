package models

import (
	"drop_os_back/dao"
	"fmt"
	"strconv"
)

type Like_info struct {
	Id            int    `json:"id" gorm:"primary_key"`
	Article_id    int    `json:"article_id"`
	Like_username string `json:"like_username"`
}

func AddLikeUserInfo(artId int, username string) error {
	LikeInfo := Like_info{}
	LikeInfo.Article_id = artId
	LikeInfo.Like_username = username
	err := dao.DB.Save(&LikeInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryLikeInfo(id int,username string)(likeInfo Like_info,err error){
	LikeInfo := Like_info{}
	idStr := strconv.Itoa(id)
	sql := "SELECT * FROM `like_infos`  WHERE article_id = " + idStr + " and  like_username = '" + username + "'"
	fmt.Println(sql)
	err = dao.DB.Raw(sql).Scan(&LikeInfo).Error
	if err != nil{
		return LikeInfo,err
	}
	return LikeInfo,nil
}

func DeleteLikeInfo(id int,username string)error{
	LikeInfo := Like_info{}
	//sql := "DELETE FROM `like_infos`  WHERE article_id = " + idStr + " and  like_username = " + username
	err := dao.DB.Where("like_username=?",username).Where("article_id=?",id).Unscoped().Delete(&LikeInfo).Error
	if err !=nil{
		return err
	}
	return nil
}

