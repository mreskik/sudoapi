package apipos

import (
	"APIANDORDER/backend/config"
	"APIANDORDER/backend/helpers"
	"APIANDORDER/backend/middleware"
	"APIANDORDER/backend/modules/apipos/pushdata"
	"APIANDORDER/backend/modules/apipos/setup"
	"APIANDORDER/backend/modules/apipos/sync"
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

	setupRouter.POST("/get_promo_list/:branch_id", setupHandler.GetPromoList)
	setupRouter.POST("/get_promo_branch/:branch_id", setupHandler.GetPromoBranch)
	setupRouter.POST("/get_promo_visit_purpose/:branch_id", setupHandler.GetPromoVisitPurpose)
	setupRouter.POST("/get_promo_type_member/:branch_id", setupHandler.GetPromoTypeMember)
	setupRouter.POST("/get_promo_category/:branch_id", setupHandler.GetPromoCategory)
	setupRouter.POST("/get_promo_sub_category/:branch_id", setupHandler.GetPromoSubCategory)
	setupRouter.POST("/get_promo_item/:branch_id", setupHandler.GetPromoItem)
	setupRouter.POST("/get_promo_day/:branch_id", setupHandler.GetPromoDay)
	setupRouter.POST("/get_promo_time/:branch_id", setupHandler.GetPromoTime)

	setupRouter.POST("/get_member_type_list/:branch_id", setupHandler.GetMemberTypeList)
	setupRouter.POST("/get_member_list/:branch_id", setupHandler.GetMemberList)


	// sync — token-based auth, no username/password
	syncRouter := router.Group("/sync", middleware.BranchTokenAuth(config.DB))
	syncHandler := sync.NewHandler(config.DB)

	syncRouter.GET("/get_data_branch/:branch_id", syncHandler.GetDataBranch)
	syncRouter.GET("/get_station_list/:branch_id", syncHandler.GetStationList)
	syncRouter.GET("/get_category_list/:branch_id", syncHandler.GetCategoryList)
	syncRouter.GET("/get_subcategory_list/:branch_id", syncHandler.GetSubCategoryList)
	syncRouter.GET("/get_tablesection_list/:branch_id", syncHandler.GetTableSectionList)
	syncRouter.GET("/get_table/:branch_id", syncHandler.GetTable)
	syncRouter.GET("/get_tax/:branch_id", syncHandler.GetMasterTax)
	syncRouter.GET("/get_terminal/:branch_id", syncHandler.GetMasterTerminal)

	syncRouter.GET("/get_item/:branch_id", syncHandler.GetItem)
	syncRouter.GET("/get_item_conv/:branch_id", syncHandler.GetItemConv)

	syncRouter.GET("/get_item_package/:branch_id", syncHandler.GetItemPackage)
	syncRouter.GET("/get_item_package_group/:branch_id", syncHandler.GetItemPackageGroup)
	syncRouter.GET("/get_item_package_detail/:branch_id", syncHandler.GetItemPackageDetail)

	syncRouter.GET("/get_pricelist/:branch_id", syncHandler.GetPriceList)
	syncRouter.GET("/get_pricelist_detail/:branch_id", syncHandler.GetPriceListDetail)

	syncRouter.GET("/get_payment_method/:branch_id", syncHandler.GetMasterPaymentMethod)
	syncRouter.GET("/get_payment_method_group/:branch_id", syncHandler.GetMasterPaymentMethodGroup)
	syncRouter.GET("/get_payment_method_type/:branch_id", syncHandler.GetMasterPaymentMethodType)
	syncRouter.GET("/get_payment_method_visit_purpose/:branch_id", syncHandler.GetMasterPaymentMethodVisitPurposes)
	syncRouter.GET("/get_branch_visit_purpose/:branch_id", syncHandler.GetMasterBranchVisitPurpose)
	syncRouter.GET("/get_visit_purpose/:branch_id", syncHandler.GetMasterVisitPurpose)
	syncRouter.GET("/get_table_section_print_category_setting/:branch_id", syncHandler.GetMasterTableSectionPrintCategorySetting)
	syncRouter.GET("/get_master_user/:branch_id", syncHandler.GetMasterUser)

	syncRouter.GET("/get_master_role_access/:branch_id", syncHandler.GetMasterRoleAccess)
	syncRouter.GET("/get_menu_app/:branch_id", syncHandler.GetMenuApp)

	syncRouter.GET("/get_promo_list/:branch_id", syncHandler.GetPromoList)
	syncRouter.GET("/get_promo_branch/:branch_id", syncHandler.GetPromoBranch)
	syncRouter.GET("/get_promo_visit_purpose/:branch_id", syncHandler.GetPromoVisitPurpose)
	syncRouter.GET("/get_promo_type_member/:branch_id", syncHandler.GetPromoTypeMember)
	syncRouter.GET("/get_promo_category/:branch_id", syncHandler.GetPromoCategory)
	syncRouter.GET("/get_promo_sub_category/:branch_id", syncHandler.GetPromoSubCategory)
	syncRouter.GET("/get_promo_item/:branch_id", syncHandler.GetPromoItem)
	syncRouter.GET("/get_promo_day/:branch_id", syncHandler.GetPromoDay)
	syncRouter.GET("/get_promo_time/:branch_id", syncHandler.GetPromoTime)

	syncRouter.GET("/get_member_type_list/:branch_id", syncHandler.GetMemberTypeList)
	syncRouter.GET("/get_member_list/:branch_id", syncHandler.GetMemberList)

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


		_,err := http.Get(SERVER_SUDOCORE+"/api/pos/endday-jurnal/"+branch_id+"/"+dayshift_ulid)
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


		_,err := http.Get(SERVER_SUDOCORE+"/api/pos/revert-jurnal/"+dayshift_ulid)
		if err != nil {
			ctx.JSON(200, resp.SetCode(100).SetMessage(err.Error()))
			return
		}
			ctx.JSON(200, resp.SetCode(0).SetMessage("ok"))

	});
}
