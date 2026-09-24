package apipos

import (
	"APIANDORDER/backend/config"
	"APIANDORDER/backend/middleware"
	"APIANDORDER/backend/modules/apipos/endday"
	"APIANDORDER/backend/modules/apipos/heartbeat"
	"APIANDORDER/backend/modules/apipos/member"
	"APIANDORDER/backend/modules/apipos/membertopup"
	"APIANDORDER/backend/modules/apipos/mobilenotify"
	"APIANDORDER/backend/modules/apipos/mobileorder"
	"APIANDORDER/backend/modules/apipos/pushdata"
	"APIANDORDER/backend/modules/apipos/setup"
	"APIANDORDER/backend/modules/apipos/sync"
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
	setupRouter.POST("/get_branch_ops_setting/:branch_id", setupHandler.GetMasterBranchOpsSetting)
	setupRouter.POST("/get_master_image/:branch_id", setupHandler.GetMasterImage)
	setupRouter.POST("/get_master_image_customer_display/:branch_id", setupHandler.GetMasterImageCustomerDisplay)
	setupRouter.POST("/get_master_image_kiosk/:branch_id", setupHandler.GetMasterImageKiosk)
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
	setupRouter.POST("/get_promo_apply_to/:branch_id", setupHandler.GetPromoApplyTo)

	setupRouter.POST("/get_member_type_list/:branch_id", setupHandler.GetMemberTypeList)
	setupRouter.POST("/get_member_list/:branch_id", setupHandler.GetMemberList)

	setupRouter.POST("/get_notes_menu_list/:branch_id", setupHandler.GetNotesMenuList)
	setupRouter.POST("/get_notes_menu_category/:branch_id", setupHandler.GetNotesMenuCategory)
	setupRouter.POST("/get_notes_menu_sub_category/:branch_id", setupHandler.GetNotesMenuSubCategory)
	setupRouter.POST("/get_notes_menu_detail/:branch_id", setupHandler.GetNotesMenuDetail)

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
	syncRouter.GET("/get_item_package_detail_pricelist/:branch_id", syncHandler.GetItemPackageDetailPricelist)

	syncRouter.GET("/get_pricelist/:branch_id", syncHandler.GetPriceList)
	syncRouter.GET("/get_pricelist_detail/:branch_id", syncHandler.GetPriceListDetail)

	syncRouter.GET("/get_payment_method/:branch_id", syncHandler.GetMasterPaymentMethod)
	syncRouter.GET("/get_payment_method_group/:branch_id", syncHandler.GetMasterPaymentMethodGroup)
	syncRouter.GET("/get_payment_method_type/:branch_id", syncHandler.GetMasterPaymentMethodType)
	syncRouter.GET("/get_payment_method_visit_purpose/:branch_id", syncHandler.GetMasterPaymentMethodVisitPurposes)
	syncRouter.GET("/get_branch_visit_purpose/:branch_id", syncHandler.GetMasterBranchVisitPurpose)
	syncRouter.GET("/get_branch_ops_setting/:branch_id", syncHandler.GetMasterBranchOpsSetting)
	syncRouter.GET("/get_master_image/:branch_id", syncHandler.GetMasterImage)
	syncRouter.GET("/get_master_image_customer_display/:branch_id", syncHandler.GetMasterImageCustomerDisplay)
	syncRouter.GET("/get_master_image_kiosk/:branch_id", syncHandler.GetMasterImageKiosk)
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
	syncRouter.GET("/get_promo_apply_to/:branch_id", syncHandler.GetPromoApplyTo)

	syncRouter.GET("/get_member_type_list/:branch_id", syncHandler.GetMemberTypeList)
	syncRouter.GET("/get_member_list/:branch_id", syncHandler.GetMemberList)

	syncRouter.GET("/get_notes_menu_list/:branch_id", syncHandler.GetNotesMenuList)
	syncRouter.GET("/get_notes_menu_category/:branch_id", syncHandler.GetNotesMenuCategory)
	syncRouter.GET("/get_notes_menu_sub_category/:branch_id", syncHandler.GetNotesMenuSubCategory)
	syncRouter.GET("/get_notes_menu_detail/:branch_id", syncHandler.GetNotesMenuDetail)

	////////////////

	pushHandler := pushdata.NewHander(config.DB)
	pushRouter := router.Group("/push")

	pushRouter.POST("/data_order", pushHandler.PushDataPosOrder)
	pushRouter.POST("/data_order_detail", pushHandler.PushDataPosOrderDetail)
	pushRouter.POST("/data_order_detail_package", pushHandler.PushDataPosOrderDetailPackage)
	pushRouter.POST("/data_order_payment", pushHandler.PushDataPosOrderPayment)
	pushRouter.POST("/data_dayshift", pushHandler.PushDataPosDayShift)
	pushRouter.POST("/data_dayshift_detail", pushHandler.PushDataPosDayShiftDetail)
	pushRouter.POST("/data_remove_item_before_save", pushHandler.PushDataPosRemoveItemBeforeSave)
	pushRouter.POST("/data_remove_item_before_save_package", pushHandler.PushDataPosRemoveItemBeforeSavePackage)

	// ENDDAY JURNAL & REVERT -- token-based auth (sama kayak syncRouter), token diteruskan ke
	// sudocore2 biar divalidasi ulang di sana juga (lihat backend/modules/apipos/endday).
	enddayRouter := router.Group("/endday", middleware.BranchTokenAuth(config.DB))
	enddayHandler := endday.NewHandler()
	enddayRouter.GET("/jurnal/:branch_id/:dayshift_ulid", enddayHandler.RequestEndDay)
	enddayRouter.GET("/jurnal-revert/:branch_id/:dayshift_ulid", enddayHandler.RequestEndDayRevert)

	// MEMBER -- live lookup by phone_number (bukan sync pull), token-based auth sama kayak endday.
	memberRouter := router.Group("/member", middleware.BranchTokenAuth(config.DB))
	memberHandler := member.NewHandler(config.DB)
	memberRouter.GET("/:branch_id/by-phone/:phone_number", memberHandler.CheckByPhone)

	// MEMBER TOPUP -- live, real-time, dari channel online (pos/kiosk/mobile), token-based
	// auth sama kayak member/endday. Lihat backend/modules/apipos/membertopup.
	memberTopupRouter := router.Group("/member-topup", middleware.BranchTokenAuth(config.DB))
	memberTopupHandler := membertopup.NewHandler(config.DB)
	memberTopupRouter.POST("/:branch_id", memberTopupHandler.CreateTopup)
	memberTopupRouter.GET("/:branch_id/check-status/:reference_number", memberTopupHandler.CheckStatus)

	// MOBILE NOTIFY -- relay "ada order mobile baru" (sudomobile pg_notify) ke worker POS lewat
	// WebSocket, per branch. Token-based auth SAMA kayak sync/endday/member (branch token yang
	// sama). Listener (LISTEN mb_order_paid) jalan SEKALI di sini, dipegang selama proses hidup
	// -- lihat mobilenotify/listener.go buat detail koneksi terpisah dari config.DB (pool).
	mobileNotifyHub := mobilenotify.NewHub()
	go mobilenotify.RunListener(mobileNotifyHub)
	mobileNotifyHandler := mobilenotify.NewHandler(mobileNotifyHub)
	router.GET("/ws/mobile-order/:branch_id", middleware.BranchTokenAuth(config.DB), mobileNotifyHandler.Serve)

	// MOBILE ORDER PULL -- kandidat mb_order (paid, pulled_at IS NULL) buat ditarik POS +
	// ack setelah beres diproses lokal. Token auth SAMA (branch token), gak dibedain dari
	// mobilenotify di atas. Lihat mobileorder/mobileorder_service.go buat detail query.
	mobileOrderHandler := mobileorder.NewHandler(config.DB)
	mobileOrderRouter := router.Group("/mobile-order", middleware.BranchTokenAuth(config.DB))
	mobileOrderRouter.GET("/get_pending/:branch_id", mobileOrderHandler.GetPending)
	mobileOrderRouter.POST("/ack/:order_number", mobileOrderHandler.Ack)

	// HEARTBEAT -- worker POS kirim tiap 30 detik (command heartbeat:send, belum dibangun),
	// upsert branch_heartbeat.last_ping_at, dibaca sudomobile buat barrier order +
	// flag_status_store_open (belum dibangun juga). Token auth SAMA (branch token).
	heartbeatHandler := heartbeat.NewHandler(config.DB)
	router.POST("/heartbeat/:branch_id", middleware.BranchTokenAuth(config.DB), heartbeatHandler.Ping)
}
