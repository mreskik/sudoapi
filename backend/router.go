package backend

import (
	"APIANDORDER/backend/modules/apipos"

	"github.com/gin-gonic/gin"
)

func Router(app *gin.Engine) {
	//apipos router
	apipos.Register(app)

	//index
	app.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "")
	})

}
