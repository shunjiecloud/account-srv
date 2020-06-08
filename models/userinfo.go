package models

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shunjiecloud/account-srv/modules"
	"github.com/shunjiecloud/account-srv/utility"
)

type UserInfo struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int8   `json:"gender"`
	Mail   string `json:"mail"`
	Phone  string `json:"phone"`
}

func fomatUser2UserInfo(user *User) (*UserInfo, error) {
	var info UserInfo
	info.Id = user.ID
	info.Name = user.Name

	info.Gender = user.Gender
	info.Mail = user.Mail
	info.Phone = user.Phone

	avatarUrl, err := modules.ModuleContext.ImgBucket.SignURL(user.Avatar, oss.HTTPGet, 600)
	if err != nil {
		return nil, err
	}
	info.Avatar = avatarUrl
	return &info, nil
}

//  用户名，用户电话，用户邮箱唯一
//  用户id，获取用户信息
func GetUserInfoByUserId(userId int64) (*UserInfo, error) {
	var user User
	err := modules.ModuleContext.DefaultDB.Model(User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	userInfo, err := fomatUser2UserInfo(&user)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

//  用户名，获取用户信息
func GetUserInfoByName(name string) (*UserInfo, error) {
	var user User
	err := modules.ModuleContext.DefaultDB.Model(User{}).Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	userInfo, err := fomatUser2UserInfo(&user)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

//  手机号，获取用户信息
func GetUserInfoByPhone(phone string) (*UserInfo, error) {
	var user User
	err := modules.ModuleContext.DefaultDB.Model(User{}).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	userInfo, err := fomatUser2UserInfo(&user)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

//  邮箱，获取用户信息
func GetUserInfoByMail(mail string) (*UserInfo, error) {
	var user User
	err := modules.ModuleContext.DefaultDB.Model(User{}).Where("mail = ?", mail).First(&user).Error
	if err != nil {
		return nil, err
	}
	userInfo, err := fomatUser2UserInfo(&user)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

//  账号，获取用户信息
func GetUserInfoByAccount(account string) (*UserInfo, error) {
	var err error
	var userInfo *UserInfo
	idType := 0
	isOk := utility.VerifyEmailFormat(account)
	if isOk == true {
		idType = IDTYPE_MAIL
	} else {
		isOk = utility.VerifyPhoneFormat(account)
		if isOk == true {
			idType = IDTYPE_PHONE
		} else {
			idType = IDTYPE_NAME
		}
	}
	switch idType {
	case IDTYPE_MAIL:
		userInfo, err = GetUserInfoByMail(account)
	case IDTYPE_PHONE:
		userInfo, err = GetUserInfoByPhone(account)
	case IDTYPE_NAME:
		userInfo, err = GetUserInfoByName(account)
	}
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
