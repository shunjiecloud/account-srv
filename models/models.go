package models

import (
	"github.com/jinzhu/gorm"
)

//  需要初始化的表
func InitTables(db *gorm.DB) {
	//  User用户表
	if db.HasTable(&User{}) == false {
		db.CreateTable(&User{})
	}
}
