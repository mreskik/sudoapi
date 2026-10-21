package pushdata

import (
	"APIANDORDER/backend/helpers"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type PushDataHandler struct {
	PushDataService *PushDataService
}

func NewHander(db *bun.DB) *PushDataHandler {
	return &PushDataHandler{PushDataService: NewPushDataService(db)}
}

// requireMatchingBranch (2026-10-16): cross-check branch_id SETIAP baris body sama branch_id
// hasil lookup token (c.Get("pos_branch_id"), di-set middleware.BranchTokenAuth) -- token valid
// doang gak cukup, device branch A gak boleh nyuntik data branch B meski somehow tau/pakai
// token A. Cuma dipasang di endpoint yang body-nya punya branch_id LANGSUNG per baris
// (data_order/data_dayshift/data_remove_item_before_save) -- endpoint detail/package/payment
// gak punya branch_id sendiri (cuma referensi ULID ke baris induk yang udah dicek di sini),
// jadi gak dicek ulang row-level di situ (perlindungannya tetap ada lewat token-nya sendiri).
func requireMatchingBranch(c *gin.Context, branchIds []int64) error {
	posBranchId, exists := c.Get("pos_branch_id")
	if !exists {
		return errors.New("branch id tidak ditemukan dari token")
	}
	tokenBranchId, ok := posBranchId.(int64)
	if !ok {
		return errors.New("branch id dari token tidak valid")
	}
	for _, id := range branchIds {
		if id != tokenBranchId {
			return errors.New("branch_id tidak cocok dengan token yang dipakai")
		}
	}
	return nil
}

func (this *PushDataHandler) PushDataPosOrder(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_order := PosOrderDTO{}

	err := c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	branchIds := make([]int64, 0, len(list_data_order.ListOrder))
	for _, v := range list_data_order.ListOrder {
		branchIds = append(branchIds, v.BranchID)
	}
	if err := requireMatchingBranch(c, branchIds); err != nil {
		c.JSON(200, res.SetCode(100).SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSOrder(c, list_data_order)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetCode(0).SetMessage("success push data data order!"))

}

func (this *PushDataHandler) PushDataPosOrderDetail(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_order := PosOrderDetailDTO{}

	err := c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSOrderDetail(c, list_data_order)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetCode(0).SetMessage("success push data!"))

}

func (this *PushDataHandler) PushDataPosOrderDetailPackage(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_order := PosOrderDetailPackageDTO{}

	err := c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSOrderDetailPackage(c, list_data_order)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetMessage("success push data!").SetCode(0))

}

func (this *PushDataHandler) PushDataPosOrderPayment(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_order := PosOrderPaymentDTO{}

	err := c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSOrderPayment(c, list_data_order)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
}

func (this *PushDataHandler) PushDataPosDayShift(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_dayshift := PosDayShiftDTO{}

	err := c.ShouldBindJSON(&list_data_dayshift)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	branchIds := make([]int64, 0, len(list_data_dayshift.ListDayShift))
	for _, v := range list_data_dayshift.ListDayShift {
		branchIds = append(branchIds, v.BranchID)
	}
	if err := requireMatchingBranch(c, branchIds); err != nil {
		c.JSON(200, res.SetCode(100).SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSDayShift(c, list_data_dayshift)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
}

func (this *PushDataHandler) PushDataPosDayShiftDetail(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_dayshift_detail := PosDayShiftDetailDTO{}

	err := c.ShouldBindJSON(&list_data_dayshift_detail)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSDayShiftDetail(c, list_data_dayshift_detail)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
}

func (this *PushDataHandler) PushDataPosRemoveItemBeforeSave(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_request := PosRemoveItemBeforeSaveDTO{}

	err := c.ShouldBindJSON(&list_data_request)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	branchIds := make([]int64, 0, len(list_data_request.ListRemoveItemBeforeSave))
	for _, v := range list_data_request.ListRemoveItemBeforeSave {
		branchIds = append(branchIds, v.BranchID)
	}
	if err := requireMatchingBranch(c, branchIds); err != nil {
		c.JSON(200, res.SetCode(100).SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSRemoveItemBeforeSave(c, list_data_request)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
}

func (this *PushDataHandler) PushDataPosRemoveItemBeforeSavePackage(c *gin.Context) {
	res := helpers.NewResponse()
	list_data_request := PosRemoveItemBeforeSavePackageDTO{}

	err := c.ShouldBindJSON(&list_data_request)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSRemoveItemBeforeSavePackage(c, list_data_request)

	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
}
