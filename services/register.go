package services

import (
	merr "github.com/micro/go-micro/v2/errors"
	"github.com/shunjiecloud/account-srv/models"
	"github.com/shunjiecloud/account-srv/utility"
	"github.com/shunjiecloud/errors"
)

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
	isExisted, err := models.IsEmailExisted(email)
	if err != nil {
		return errors.Adapt(err)
	}
	if isExisted == true {
		return merr.BadRequest("email is existed", "%s", email)
	}
	isExisted, err = models.IsUserNameExisted(username)
	if err != nil {
		return errors.Adapt(err)
	}
	if isExisted == true {
		return merr.BadRequest("username is existed", "%s", username)
	}
	isExisted, err = models.IsPhoneExisted(phone)
	if err != nil {
		return errors.Adapt(err)
	}
	if isExisted == true {
		return merr.BadRequest("phone is existed", "%s", phone)
	}
	return nil
}
