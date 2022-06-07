package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var(
	DB *gorm.DB
)


func InitMySQL()(err error){
	dsn := 
	DB, err = gorm.Open("mysql", dsn)
	if err != nil{
		return
	}
	return DB.DB().Ping()
}


func Close(){
	DB.Close()
}
