package backend

import (
	"APIANDORDER/backend/modules/apipos"
	"os"

	"github.com/gin-gonic/gin"
)

func Router(app *gin.Engine) {
	// serve sudocore2 storage (images, files) — path configured via SUDOCORE_STORAGE_PATH
	storagePath := os.Getenv("SUDOCORE_STORAGE_PATH")
	if storagePath != "" {
		app.Static("/storage", storagePath)
	}

	//apipos router
	apipos.Register(app)

	//index
	app.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "")
	})

}
