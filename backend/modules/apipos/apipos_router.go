package apipos

import (
	"APIANDORDER/backend/config"
	"APIANDORDER/backend/modules/apipos/setup"
	"net/http"

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

}
