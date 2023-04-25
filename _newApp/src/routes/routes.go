package routes

import (
	"github.com/kukkar/common-golang/pkg/middleware"
	appConf "github.com/kukkar/{{APP_NAME}}/conf"
	"github.com/kukkar/{APP_NAME}}/src/otp"
	controller "github.com/kukkar/{{APP_NAME}}/src/{{APP_NAME}}_controllers"
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {

	gConf, err := appConf.GetGlobalConfig()
	if err != nil {
		panic(err)
	}
	appConfig, err := appConf.GetAppConfig()
	if err != nil {
		panic(err)
	}
	v1 := route.Group(string(gConf.AppName) + "/v1")
	{
		defaultMiddleware := middleware.DefaultMiddleware{}
		v1.GET("/hellworld", defaultMiddleware.MonitorRequest(), controller.HelloWorld)
	}
}
