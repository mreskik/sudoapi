package soldout

import (
	"net/http"

	"APIANDORDER/backend/helpers"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type Handler struct {
	Service *Service
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{Service: NewService(db)}
}

// PushSoldOut: POST /pos/push/sold_out -- tombol manual di POS (BUKAN auto-terjadwal kayak
// SyncPush.php), kirim snapshot PENUH item_id (master_item.id, item fisik -- sama persis
// mr_item.id/flag_soldout di POS) yang LAGI sold out di branch itu sekarang. ERP replace
// total baris branch itu (DELETE+INSERT, lihat soldout_service.go) -- jadi kalau ada
// perubahan di POS (item di-unsoldout-kan), tinggal push ulang, baris lama otomatis kehapus.
func (h *Handler) PushSoldOut(c *gin.Context) {
	res := helpers.NewResponse()

	var req PushSoldOutRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage(err.Error()))
		return
	}
	if req.BranchID == 0 {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("branch_id wajib diisi"))
		return
	}

	// cross-check branch_id body sama branch_id hasil lookup token (2026-10-16) -- token valid
	// doang gak cukup, device branch A gak boleh nyuntik sold out buat branch B.
	posBranchId, exists := c.Get("pos_branch_id")
	if !exists {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("branch id tidak ditemukan dari token"))
		return
	}
	if tokenBranchId, ok := posBranchId.(int64); !ok || tokenBranchId != req.BranchID {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("branch_id tidak cocok dengan token yang dipakai"))
		return
	}

	if err := h.Service.PushSoldOut(c, req.BranchID, req.ItemIDs); err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, res.Success().SetMessage("success push data sold out!"))
}
