module github.com/PenguinCats/Unison-Elastic-Compute

go 1.15

require (
	github.com/PenguinCats/Unison-Docker-Controller v0.0.0
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/astaxie/beego v1.12.3
	github.com/gin-gonic/gin v1.7.2
	github.com/go-ini/ini v1.62.0
	github.com/go-playground/validator/v10 v10.8.0 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomodule/redigo v2.0.1-0.20180401191855-9352ab68be13+incompatible
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)

replace github.com/PenguinCats/Unison-Docker-Controller => ../Unison-Docker-Controller
