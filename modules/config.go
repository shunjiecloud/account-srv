package modules

import (
	"os"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)

type RedisConfig struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type DBConfig struct {
	Host        string `json:"host"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	DbType      string `json:"dbType"`
	TablePrefix string `json:"tbPrefix"`
}

type OssClientConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
}

func setupConfig() {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		panic("CONFIG_FILE_PATH is error")
	}
	if err := config.Load(file.NewSource(
		file.WithPath(configFilePath),
	)); err != nil {
		panic(err)
	}
}

func init() {
	setupConfig()
}
