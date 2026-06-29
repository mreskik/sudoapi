package apipos

import (
	"APIANDORDER/backend/config"
	"APIANDORDER/backend/helpers"
	"APIANDORDER/backend/modules/apipos/pushdata"
	"APIANDORDER/backend/modules/apipos/setup"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine) {
	router := app.Group("/pos")

	router.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusForbidden, "Forbidden")
	})

	// miu := New()
	// router.GET("/:angka", miu.Init)
	// router.GET("/cek", miu.Cek)
	//setup
	setupRouter := router.Group("/setup")

	setupHandler := setup.NewHandler(config.DB)
	setupRouter.POST("/get_branch_list", setupHandler.GetBranchList)

	setupRouter.POST("/get_data_branch/:branch_id", setupHandler.GetDataBranch)
	setupRouter.POST("/get_station_list/:branch_id", setupHandler.GetStationList)
	setupRouter.POST("/get_category_list/:branch_id", setupHandler.GetCategoryList)
	setupRouter.POST("/get_subcategory_list/:branch_id", setupHandler.GetSubCategoryList)
	setupRouter.POST("/get_tablesection_list/:branch_id", setupHandler.GetTableSectionList)
	setupRouter.POST("/get_table/:branch_id", setupHandler.GetTable)
	setupRouter.POST("/get_tax/:branch_id", setupHandler.GetMasterTax)
	setupRouter.POST("/get_terminal/:branch_id", setupHandler.GetMasterTerminal)

	setupRouter.POST("/get_item/:branch_id", setupHandler.GetItem)
	setupRouter.POST("/get_item_conv/:branch_id", setupHandler.GetItemConv)

	setupRouter.POST("/get_item_package/:branch_id", setupHandler.GetItemPackage)
	setupRouter.POST("/get_item_package_group/:branch_id", setupHandler.GetItemPackageGroup)
	setupRouter.POST("/get_item_package_detail/:branch_id", setupHandler.GetItemPackageDetail)

	setupRouter.POST("/get_pricelist/:branch_id", setupHandler.GetPriceList)
	setupRouter.POST("/get_pricelist_detail/:branch_id", setupHandler.GetPriceListDetail)

	setupRouter.POST("/get_payment_method/:branch_id", setupHandler.GetMasterPaymentMethod)
	setupRouter.POST("/get_payment_method_group/:branch_id", setupHandler.GetMasterPaymentMethodGroup)
	setupRouter.POST("/get_payment_method_type/:branch_id", setupHandler.GetMasterPaymentMethodType)
	setupRouter.POST("/get_payment_method_visit_purpose/:branch_id", setupHandler.GetMasterPaymentMethodVisitPurposes)
	setupRouter.POST("/get_branch_visit_purpose/:branch_id", setupHandler.GetMasterBranchVisitPurpose)
	setupRouter.POST("/get_visit_purpose/:branch_id", setupHandler.GetMasterVisitPurpose)
	setupRouter.POST("/get_table_section_print_category_setting/:branch_id", setupHandler.GetMasterTableSectionPrintCategorySetting)
	setupRouter.POST("/get_master_user/:branch_id", setupHandler.GetMasterUser)

	setupRouter.POST("/get_master_role_access/:branch_id", setupHandler.GetMasterRoleAccess)
	setupRouter.POST("/get_menu_app/:branch_id", setupHandler.GetMenuApp)


	////////////////

	pushHandler := pushdata.NewHander(config.DB)
	pushRouter := router.Group("/push")

	pushRouter.POST("/data_order", pushHandler.PushDataPosOrder)
	pushRouter.POST("/data_order_detail", pushHandler.PushDataPosOrderDetail)
	pushRouter.POST("/data_order_detail_package", pushHandler.PushDataPosOrderDetailPackage)
	pushRouter.POST("/data_order_payment", pushHandler.PushDataPosOrderPayment)

	// ENDDAY JURNAL & REVERT

	enddayRouter := router.Group("/endday")
	enddayRouter.GET("/jurnal/:branch_id/:dayshift_ulid", func(ctx *gin.Context) {
		resp := helpers.NewResponse()

		branch_id := ctx.Param("branch_id")
		dayshift_ulid := ctx.Param("dayshift_ulid")
		SERVER_SUDOCORE := "http://"+os.Getenv("APP_SUDOCORE_HOST")+":"+os.Getenv("APP_SUDOCORE_PORT")


		_,err := http.Get(SERVER_SUDOCORE+"/pos/endday-jurnal/"+branch_id+"/"+dayshift_ulid)
		if err != nil {

			log.Println(err)
			ctx.JSON(200, resp.SetCode(100).SetMessage(err.Error()))
			return
		}
			ctx.JSON(200, resp.SetCode(0).SetMessage("ok"))

	});

	enddayRouter.GET("/jurnal-revert/:dayshift_ulid", func(ctx *gin.Context) {
		resp := helpers.NewResponse()

		dayshift_ulid := ctx.Param("dayshift_ulid")
		SERVER_SUDOCORE := "http://"+os.Getenv("APP_SUDOCORE_HOST")+":"+os.Getenv("APP_SUDOCORE_PORT")


		_,err := http.Get(SERVER_SUDOCORE+"/pos/revert-jurnal/"+dayshift_ulid)
		if err != nil {
			ctx.JSON(200, resp.SetCode(100).SetMessage(err.Error()))
			return
		}
			ctx.JSON(200, resp.SetCode(0).SetMessage("ok"))

	});
}
