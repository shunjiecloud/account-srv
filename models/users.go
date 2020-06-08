package models

import (
	"github.com/jinzhu/gorm"
	merr "github.com/micro/go-micro/v2/errors"
	"github.com/shunjiecloud/account-srv/modules"
	"github.com/shunjiecloud/account-srv/utility"
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
	//  用户名，用户电话，用户邮箱唯一
	err = RegistrationCheck(user.Name, user.Mail, user.Phone)
	if err != nil {
		return
	}
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

//  用户是否可注册，用户名，用户电话，用户邮箱唯一
func RegistrationCheck(username, email, phone string) error {
	//  用户名大于3个字符，并且不能是邮箱或手机格式
	if len(username) < 3 {
		return merr.BadRequest("username is less than 3", "%s", username)
	}
	if utility.VerifyEmailFormat(username) == true {
		return merr.BadRequest("can not use mail as username", "%s", email)
	}
	if utility.VerifyPhoneFormat(username) == true {
		return merr.BadRequest("can not use phone as username", "%s", phone)
	}
	//  邮箱格式校验
	if utility.VerifyEmailFormat(email) == false {
		return merr.BadRequest("email is invalid", "%s", email)
	}
	//  手机号格式校验
	if utility.VerifyPhoneFormat(phone) == false {
		return merr.BadRequest("phone is invalid", "%s", phone)
	}

	//  唯一性校验
	isExisted, err := IsEmailExisted(email)
	if err != nil {
		return errors.Adapt(err)
	}
	if isExisted == true {
		return merr.BadRequest("email is existed", "%s", email)
	}
	isExisted, err = IsUserNameExisted(username)
	if err != nil {
		return errors.Adapt(err)
	}
	if isExisted == true {
		return merr.BadRequest("username is existed", "%s", username)
	}
	isExisted, err = IsPhoneExisted(phone)
	if err != nil {
		return errors.Adapt(err)
	}
	if isExisted == true {
		return merr.BadRequest("phone is existed", "%s", phone)
	}
	return nil
}
