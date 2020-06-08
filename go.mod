module github.com/shunjiecloud/account-srv

go 1.14

require (
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/jinzhu/gorm v1.9.12
	github.com/micro/go-micro/v2 v2.8.0
	github.com/o1egl/govatar v0.3.0
	github.com/shunjiecloud-proto/account v0.0.0-20200606193357-d8b622034528
	github.com/shunjiecloud-proto/captcha v0.0.0-20200606113732-01d020f2eff8
	github.com/shunjiecloud-proto/encrypt v0.0.0-20200605191118-52b44ce39445
	github.com/shunjiecloud/errors v1.0.3-0.20200427091440-d2c8251bbc81
	github.com/shunjiecloud/pkg v0.0.0-20200608213205-7936a725a0c8
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
