package main

import (
	"log"

	"github.com/micro/go-micro/v2"
	account_proto "github.com/shunjiecloud-proto/account/proto"
	"github.com/shunjiecloud/account-srv/models"
	"github.com/shunjiecloud/account-srv/modules"
	"github.com/shunjiecloud/account-srv/services"
)

func main() {
	//  Create srv
	service := micro.NewService(
		micro.Name("go.micro.srv.account"),
	)

	//  init modules
	modules.Setup()

	//  init db, tables init must after modules initd
	models.InitTables(modules.ModuleContext.DefaultDB)

	//  init service
	service.Init()

	//  register Handlers
	account_proto.RegisterAccountHandler(service.Server(), new(services.AccountService))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
