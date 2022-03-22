package api

import (
	"github.com/gin-gonic/gin"
	"synod/conf"
)

func Run() error {
	app := gin.Default()

	object := app.Group("/objects")
	{
		obj := &Object{}
		object.GET("/*path", obj.loadObject)
		object.PUT("/*path", obj.putObject)
	}

	return app.Run(conf.String("api.address"))
}
