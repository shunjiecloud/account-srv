package models

import (
	"github.com/shunjiecloud/account-srv/modules"
)

//  密码检查
func UserPasswordCheck(userId uint, password string) (bool, error) {
	var user User
	err := modules.ModuleContext.DefaultDB.Model(User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return false, err
	}
	if user.Password != password {
		return false, nil
	}
	return true, nil
}
