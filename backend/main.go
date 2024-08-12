package main

import (
	"github.com/glebarez/sqlite"
	"github.com/linxlib/fw"
	"github.com/linxlib/fw/middlewares"
	"github.com/linxlib/fw_openapi"
	"github.com/linxlib/godeploy/controllers"
	"github.com/linxlib/godeploy/controllers/models"
	middlewares2 "github.com/linxlib/godeploy/middlewares"
	"github.com/linxlib/godeploy/middlewares/session"
	"github.com/linxlib/godeploy/middlewares/weblog"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// @Title A Simple Go Deploy Server
// @Version v0.0.1
// @Contact linx https://github.com/linxlib/fw slk1133@qq.com
// @Description <a href="https://github.com/linxlib/fw">Link</a>

//go:generate go run github.com/linxlib/astp/astpg -o gen.json
func main() {
	s := fw.New()
	db, err := gorm.Open(sqlite.Open("deploy.db"))
	if err != nil {
		logrus.WithError(err).Error("connect to db failed")
		return
	}
	err = db.AutoMigrate(&models.User{}, &models.Service{})
	if err != nil {
		logrus.WithError(err).Error("auto migrate failed")
		return
	}
	s.Map(db)

	fw_openapi.NewOpenAPIFromFWServer(s, "openapi.yaml")

	s.Use(middlewares2.NewStaticMiddleware("./static"))
	//s.Use(middlewares.NewWebsocketHubMiddleware())
	s.Use(middlewares.NewLoggerMiddleware(nil))
	s.Use(middlewares.NewDefaultCorsMiddleware())
	s.Use(weblog.NewWebLogMiddleware())

	//s.Use(fwOpenApi.NewOpenApiMiddleware())

	s.Use(middlewares.NewRecoveryMiddleware(&middlewares.RecoveryOptions{
		NiceWeb: true,
		Console: true,
	}, nil))
	s.Use(session.NewSessionMiddleware("redis", map[string]any{
		"Redis_Addr": "10.10.0.16:6379",
	}))

	s.RegisterRoute(new(controllers.HomeController))
	s.RegisterRoute(new(controllers.DeployController))
	s.RegisterRoute(controllers.NewServiceController(db))
	s.RegisterRoute(controllers.NewUserController(db))

	s.Run()
}
