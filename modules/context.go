package modules

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	proto_captcha "github.com/shunjiecloud-proto/captcha/proto"
	proto_encrypt "github.com/shunjiecloud-proto/encrypt/proto"
)

type moduleWrapper struct {
	DefaultDB        *gorm.DB
	OssClient        *oss.Client
	ImgBucket        *oss.Bucket
	CaptchaSrvClient proto_captcha.CaptchaService
	EncryptSrvClient proto_encrypt.EncryptService
}

//ModuleContext 模块上下文
var ModuleContext moduleWrapper

//Setup 初始化Modules
func Setup() {
	//  db
	var dbConfig DBConfig
	if err := config.Get("config", "db").Scan(&dbConfig); err != nil {
		panic(err)
	}
	db, err := gorm.Open(dbConfig.DbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Name))
	if err != nil {
		panic(fmt.Sprintf("%v", err.Error()))
	}
	//  连接数配置
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	//  设置表名属性
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return dbConfig.TablePrefix + defaultTableName
	}
	db.SingularTable(true)
	ModuleContext.DefaultDB = db

	//  oss
	var imgOssClientConfig OssClientConfig
	if err := config.Get("config", "imgOssClient").Scan(&imgOssClientConfig); err != nil {
		panic(err)
	}
	ModuleContext.OssClient, err = oss.New(imgOssClientConfig.Endpoint, imgOssClientConfig.AccessKey, imgOssClientConfig.SecretKey)
	if err != nil {
		panic(err)
	}
	ModuleContext.ImgBucket, err = ModuleContext.OssClient.Bucket(imgOssClientConfig.Bucket)
	if err != nil {
		panic(err)
	}

	//  captcha-srv client
	m_service := micro.NewService()
	ModuleContext.CaptchaSrvClient = proto_captcha.NewCaptchaService("go.micro.srv.captcha", m_service.Client())
	//  encrypt-srv client
	ModuleContext.EncryptSrvClient = proto_encrypt.NewEncryptService("go.micro.srv.encrypt", m_service.Client())
}
