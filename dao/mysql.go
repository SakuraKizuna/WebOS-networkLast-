package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var(
	DB *gorm.DB
)


func InitMySQL()(err error){
	dsn := "root:020804@(127.0.0.1:3306)/sakura?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:lijianmin@(127.0.0.1:3306)/signWeb?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:lijianmin@(192.144.235.64:3306)/signWeb?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "rooth:Huawei&519@(172.17.137.39:3306)/519cloud?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil{
		return
	}
	return DB.DB().Ping()
}


func Close(){
	DB.Close()
}