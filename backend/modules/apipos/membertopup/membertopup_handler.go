package membertopup

import (
	"APIANDORDER/backend/helpers"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type MemberTopupHandler struct {
	Service *MemberTopupService
}

func NewHandler(db *bun.DB) *MemberTopupHandler {
	return &MemberTopupHandler{Service: NewService(db)}
}

// CreateTopup: POST /pos/member-topup/:branch_id
func (h *MemberTopupHandler) CreateTopup(c *gin.Context) {
	res := helpers.NewResponse()

	branchID, err := strconv.ParseInt(c.Param("branch_id"), 10, 64)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch_id tidak valid"))
		return
	}

	req := CreateTopupRequestDTO{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, res.GeneralError().SetMessage(err.Error()))
		return
	}

	data, err := h.Service.CreateTopup(c, branchID, req)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.Success().SetData(data))
}

// CheckStatus: GET /pos/member-topup/:branch_id/check-status/:reference_number
func (h *MemberTopupHandler) CheckStatus(c *gin.Context) {
	res := helpers.NewResponse()

	referenceNumber := c.Param("reference_number")
	if referenceNumber == "" {
		c.JSON(200, res.GeneralError().SetMessage("reference_number wajib diisi"))
		return
	}

	data, err := h.Service.CheckStatus(c, referenceNumber)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage(err.Error()))
		return
	}

	c.JSON(200, res.Success().SetData(data))
}
