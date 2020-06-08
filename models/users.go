package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shunjiecloud/account-srv/modules"
	"github.com/shunjiecloud/errors"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"index,unique"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	Mail     string `json:"mail" gorm:"index,unique"`
	Phone    string `json:"phone" gorm:"index,unique"`
}

//  创建用户
func CreateUser(user *User) (userId uint, err error) {
	if err = modules.ModuleContext.DefaultDB.Create(&user).Error; err != nil {
		return
	}
	return user.ID, nil
}

//  邮箱是否存在
func IsEmailExisted(email string) (isExisted bool, err error) {
	count := 0
	err = modules.ModuleContext.DefaultDB.Model(User{}).Where("mail = ?", email).Count(&count).Error
	if err != nil {
		return false, errors.Adapt(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//  用户名是否存在
func IsUserNameExisted(name string) (isExisted bool, err error) {
	count := 0
	err = modules.ModuleContext.DefaultDB.Model(User{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, errors.Adapt(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//  电话是否存在
func IsPhoneExisted(phone string) (isExisted bool, err error) {
	count := 0
	err = modules.ModuleContext.DefaultDB.Model(User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, errors.Adapt(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
