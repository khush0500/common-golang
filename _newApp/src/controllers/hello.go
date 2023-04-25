package controllers

import (
	"fmt"

	"github.com/kukkar/common-golang/globalconst"

	"github.com/kukkar/common-golang/pkg/responsewriter"
	"github.com/kukkar/common-golang/pkg/utils"
	"github.com/kukkar/common-golang/pkg/utils/rError"
	appConf "github.com/kukkar/{{APP_NAME}}/conf"

	"github.com/gin-gonic/gin"
)

//
//HelloWorld service
func HelloWorld(c *gin.Context) {

	var rc utils.RequestContext
	if requestContext, ok := c.Get(globalconst.RequestContext); ok {
		rc = requestContext.(utils.RequestContext)
	}
	conf, err := appConf.GetAppConfig()
	if err != nil {
		err = rError.MiscError(c, err, "Unable to get appconfig")
		responsewriter.BuildResponse(c, "", err)
		return
	}
	gConfig, err := appConf.GetGlobalConfig()
	if err != nil {
		err = rError.MiscError(c, err, "Unable to get appconfig")
		responsewriter.BuildResponse(c, "", err)
		return
	}
	fmt.Printf("%v %v %v %v", rc, conf, gConfig)

	responsewriter.BuildResponse(c, "", nil)
	return
}
