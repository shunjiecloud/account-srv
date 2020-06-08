package services

import (
	"context"

	"github.com/jinzhu/gorm"
	merr "github.com/micro/go-micro/v2/errors"
	"github.com/shunjiecloud-proto/account/proto"
	proto_captcha "github.com/shunjiecloud-proto/captcha/proto"
	proto_encrypt "github.com/shunjiecloud-proto/encrypt/proto"
	"github.com/shunjiecloud/account-srv/models"
	"github.com/shunjiecloud/account-srv/modules"
)

type AccountService struct{}

func (a *AccountService) CreateUser(ctx context.Context, in *proto.CreateUserRequest, out *proto.CreateUserResponse) error {
	//  TODO : 异步生成头像
	avatarSha1, err := models.GeneralAvatar(in.Name, in.Gender)
	if err != nil {
		return err
	}

	//  写入数据库
	user := models.User{
		Name:     in.Name,
		Password: in.Password,
		Avatar:   avatarSha1,
		Gender:   int8(in.Gender),
		Mail:     in.Mail,
		Phone:    in.Phone,
	}
	userId, err := models.CreateUser(&user)
	if err != nil {
		return err
	}
	out.UserId = int64(userId)
	out.Name = user.Name
	out.Avatar = user.Avatar
	out.Gender = int32(user.Gender)
	out.Mail = user.Mail
	out.Phone = user.Phone
	return nil
}

func (a *AccountService) SignUp(ctx context.Context, in *proto.SignUpRequest, out *proto.SignUpResponse) error {
	//  校验验证码
	_, err := modules.ModuleContext.CaptchaSrvClient.CaptchaVerfify(ctx, &proto_captcha.CaptchaVerfifyRequest{
		CaptchaId: in.CaptchaId,
		Solution:  in.CaptchaSolution,
	})
	if err != nil {
		return err
	}
	//  对密码进行解密
	decryptResp, err := modules.ModuleContext.EncryptSrvClient.Decrypt(ctx, &proto_encrypt.DecryptRequest{
		CipherText: in.Password,
	})
	if err != nil {
		return err
	}

	//  TODO : 异步生成头像
	avatarSha1, err := models.GeneralAvatar(in.Name, in.Gender)
	if err != nil {
		return err
	}

	//  写入数据库
	user := models.User{
		Name:     in.Name,
		Password: decryptResp.Original,
		Avatar:   avatarSha1,
		Gender:   int8(in.Gender),
		Mail:     in.Mail,
		Phone:    in.Phone,
	}
	userId, err := models.CreateUser(&user)
	if err != nil {
		return err
	}
	out.UserId = int64(userId)
	out.Name = user.Name
	out.Avatar = user.Avatar
	out.Gender = int32(user.Gender)
	out.Mail = user.Mail
	out.Phone = user.Phone
	return nil
}

func (a *AccountService) SignIn(ctx context.Context, in *proto.SignInRequest, out *proto.SignInResponse) error {
	//  校验验证码
	// _, err := modules.ModuleContext.CaptchaSrvClient.CaptchaVerfify(ctx, &proto_captcha.CaptchaVerfifyRequest{
	// 	CaptchaId: in.CaptchaId,
	// 	Solution:  in.CaptchaSolution,
	// })
	// if err != nil {
	// 	return err
	// }
	//  获取用户信息
	userInfo, err := models.GetUserInfoByAccount(in.Account)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//  密码错误或用户不存在均返回用户或密码错误
			return merr.Unauthorized("user or password is wrong", "account:%v, password:%v", in.Account, in.Password)
		}
		return err
	}
	//  对密码进行解密
	decryptResp, err := modules.ModuleContext.EncryptSrvClient.Decrypt(ctx, &proto_encrypt.DecryptRequest{
		CipherText: in.Password,
	})
	if err != nil {
		return err
	}
	//  确认密码
	err = models.UserPasswordCheck(userInfo.Id, decryptResp.Original)
	if err != nil {
		return err
	}
	out.Avatar = userInfo.Avatar
	out.Gender = (int32)(userInfo.Gender)
	out.Name = userInfo.Name
	out.Phone = userInfo.Phone
	out.UserId = (int64)(userInfo.Id)
	out.Mail = userInfo.Mail
	return nil
}

func (a *AccountService) UserInfo(ctx context.Context, in *proto.UserInfoRequest, out *proto.UserInfoResponse) error {

	cnt := 0
	idType := models.IDTYPE_USERID
	if in.UserId != 0 {
		idType = models.IDTYPE_USERID
		cnt++
	}
	if len(in.Name) > 0 {
		idType = models.IDTYPE_NAME
		cnt++
	}
	if len(in.Mail) > 0 {
		idType = models.IDTYPE_MAIL
		cnt++
	}
	if len(in.Phone) > 0 {
		idType = models.IDTYPE_PHONE
		cnt++
	}
	if cnt != 1 {
		return merr.BadRequest("Only Need One IdType", "userid:%v, name:%v, mail:%v, phone:%v", in.UserId, in.Name, in.Mail, in.Phone)
	}

	var userInfo *models.UserInfo
	var err error
	switch idType {
	case models.IDTYPE_USERID:
		userInfo, err = models.GetUserInfoByUserId(in.UserId)

	case models.IDTYPE_NAME:
		userInfo, err = models.GetUserInfoByName(in.Name)

	case models.IDTYPE_PHONE:
		userInfo, err = models.GetUserInfoByPhone(in.Phone)

	case models.IDTYPE_MAIL:
		userInfo, err = models.GetUserInfoByMail(in.Mail)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return merr.NotFound("user not existed", "userid:%v, name:%v, phone:%v, mail:%v", in.UserId, in.Name, in.Phone, in.Mail)
		}
		return err
	}
	//  格式化数据返回
	out.UserId = int64(userInfo.Id)
	out.Name = userInfo.Name
	out.Mail = userInfo.Mail
	out.Phone = userInfo.Phone
	out.Avatar = userInfo.Avatar
	out.Gender = int32(userInfo.Gender)
	return nil
}

func (a *AccountService) RegistrationCheck(ctx context.Context, in *proto.RegistrationCheckRequest, out *proto.RegistrationCheckResponse) error {
	err := models.RegistrationCheck(in.Name, in.Mail, in.Phone)
	if err != nil {
		return nil
	}
	return nil
}

func (a *AccountService) UserPasswordCheck(ctx context.Context, in *proto.UserPasswordCheckRequest, out *proto.UserPasswordCheckResponse) error {
	userInfo, err := models.GetUserInfoByUserId(in.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//  密码错误或用户不存在均返回用户或密码错误
			return merr.Unauthorized("user or password is wrong", "userId:%v, password:%v", in.UserId, in.Password)
		}
		return err
	}
	//  密码校验
	err = models.UserPasswordCheck(userInfo.Id, in.Password)
	if err != nil {
		return err
	}
	return nil
}
