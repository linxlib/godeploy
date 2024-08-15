package main

import (
	"github.com/linxlib/fw"
	"github.com/linxlib/fw/middlewares"
	"github.com/linxlib/fw_openapi"
	"github.com/linxlib/godeploy/controllers"
	middlewares2 "github.com/linxlib/godeploy/middlewares"
	"github.com/linxlib/godeploy/middlewares/session"
	"github.com/linxlib/godeploy/middlewares/weblog"
	"github.com/linxlib/godeploy/services/db"
)

// @Title A Simple Go Deploy Server
// @Version v0.0.1
// @Contact linx https://github.com/linxlib/fw slk1133@qq.com
// @Description <a href="https://github.com/linxlib/fw">Link</a>

//go:generate go run github.com/linxlib/astp/astpg -o gen.json
func main() {
	s := fw.New()
	s.UseMapper(new(db.DBMapper))

	fw_openapi.NewOpenAPIFromFWServer(s, "openapi.yaml")

	s.Use(middlewares2.NewStaticMiddleware())
	s.Use(middlewares.NewWebsocketMiddleware())
	//s.Use(middlewares.NewWebsocketHubMiddleware())
	s.Use(middlewares.NewLoggerMiddleware(nil))
	s.Use(middlewares.NewDefaultCorsMiddleware())
	s.Use(weblog.NewWebLogMiddleware())
	s.Use(middlewares.NewServerDownMiddleware())

	//s.Use(fwOpenApi.NewOpenApiMiddleware())

	s.Use(middlewares.NewRecoveryMiddleware(&middlewares.RecoveryOptions{
		NiceWeb: true,
		Console: true,
	}, nil))
	s.Use(session.NewSessionMiddleware())
	s.Use(middlewares2.NewTestMiddleware())

	s.RegisterRoute(new(controllers.HomeController))
	s.RegisterRoute(new(controllers.Home2Controller))
	s.RegisterRoute(new(controllers.DeployController))
	s.RegisterRoute(controllers.NewServiceController())
	s.RegisterRoute(controllers.NewUserController())

	s.Start()
}
