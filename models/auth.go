package models

import (
	merr "github.com/micro/go-micro/v2/errors"
	"github.com/shunjiecloud/account-srv/modules"
)

//  密码检查
func UserPasswordCheck(userId uint, password string) error {
	var user User
	err := modules.ModuleContext.DefaultDB.Model(User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return err
	}
	if user.Password != password {
		return merr.BadRequest("password mismatch", "")
	}
	return nil
}
