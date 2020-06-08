package services

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"image/png"

	merr "github.com/micro/go-micro/v2/errors"
	"github.com/o1egl/govatar"
	"github.com/shunjiecloud/account-srv/models"
	"github.com/shunjiecloud/account-srv/modules"
)

func GeneralAvatar(name string, gender int32) (avatarSha1 string, err error) {
	//  生成avatar
	govatarGender := govatar.MALE
	if gender == models.GENDER_FEMALE {
		govatarGender = govatar.FEMALE
	}
	img, err := govatar.GenerateForUsername(govatarGender, name)
	if err != nil {
		return "", merr.InternalServerError("avatar gen failed", "gender:%v, username:%v", gender, name)
	}
	//  头像转为byte
	imgByte := make([]byte, 0)
	buff := bytes.NewBuffer(imgByte)
	err = png.Encode(buff, img)
	if err != nil {
		return "", merr.InternalServerError("avatar gen failed", "%s", err.Error())
	}
	//  计算sha1
	h := sha1.New()
	h.Write(imgByte)
	avatarSha1 = fmt.Sprintf("%x", h.Sum(nil)) + ".png"
	//  上传到oss
	err = modules.ModuleContext.ImgBucket.PutObject(avatarSha1, buff)
	if err != nil {
		return "", merr.InternalServerError("put avatar failed", "%s", err.Error())
	}
	return
}
