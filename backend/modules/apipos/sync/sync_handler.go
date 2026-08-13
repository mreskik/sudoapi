package sync

import (
	"APIANDORDER/backend/helpers"
	"APIANDORDER/backend/modules/apipos/master"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type SyncHandler struct {
	masterService *master.MasterService
}

func NewHandler(db *bun.DB) *SyncHandler {
	return &SyncHandler{masterService: master.NewMasterService(db)}
}

func parseBranchID(c *gin.Context) (int, error) {
	return strconv.Atoi(c.Param("branch_id"))
}

func (h *SyncHandler) GetDataBranch(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetDataBranch(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("error get data branch!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetStationList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetStationList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal tarik daftar station"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetCategoryList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterCategory(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data category!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetSubCategoryList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterSubCategory(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data sub category!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetTableSectionList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetTableSection(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data table section!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetTable(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetTable(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data table!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterTax(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterTax(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data master tax!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterTerminal(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterTerminal(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data terminal!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetItem(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetItem(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data item!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetItemConv(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetItemConv(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data item conv!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetItemPackage(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetItemPackage(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data item package!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetItemPackageGroup(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetItemPackageGroup(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data item package group!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetItemPackageDetail(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetItemPackageDetail(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data item package detail!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPriceList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetPriceList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data pricelist!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPriceListDetail(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetPriceListDetail(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data pricelist detail!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterPaymentMethod(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPaymentMethod(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data payment method!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterPaymentMethodGroup(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPaymentMethodGroup(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data payment method group!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterPaymentMethodType(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPaymentMethodType(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data payment method type!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterPaymentMethodVisitPurposes(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPaymentMethodVisitPurposes(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data payment method visit purpose!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterBranchVisitPurpose(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterBranchVisitPurpose(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data branch visit purpose!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterBranchOpsSetting(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterBranchOpsSetting(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data branch ops setting!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterImage(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterImage(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data master image!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterImageList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterImageList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data master image list!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterImageListApplyFor(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterImageListApplyFor(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data master image list apply for!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterVisitPurpose(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterVisitPurpose(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data visit purpose!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterTableSectionPrintCategorySetting(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetTableSectionPrintCategorySetting(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data table section print category setting!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterUser(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterUserList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data master user!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMasterRoleAccess(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterRoleAccess(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data master role access!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMenuApp(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterMenuApp(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data menu app!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromo(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoBranch(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoBranches(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo branch!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoVisitPurpose(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoVisitPurposes(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo visit purpose!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoTypeMember(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoTypeMembers(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo type member!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoCategory(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoCategories(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo category!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoSubCategory(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoSubCategories(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo sub category!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoItem(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoItems(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo item!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoDay(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoDays(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo day!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetPromoTime(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterPromoTimes(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data promo time!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMemberTypeList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterMemberType(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data member type!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}

func (h *SyncHandler) GetMemberList(c *gin.Context) {
	res := helpers.NewResponse()
	branch_id, err := parseBranchID(c)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}
	data, err := h.masterService.GetMasterMember(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data member!"))
		return
	}
	c.JSON(200, res.Success().SetData(data))
}
