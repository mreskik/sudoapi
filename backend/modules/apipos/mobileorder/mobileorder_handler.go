package mobileorder

import (
	"net/http"
	"strconv"

	"APIANDORDER/backend/helpers"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type Handler struct {
	DB *bun.DB
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{DB: db}
}

// GetPending: GET /pos/mobile-order/get_pending/:branch_id -- kandidat mb_order (paid,
// pulled_at IS NULL) buat ditarik POS. Response tetap Success walau kosong (array kosong,
// bukan error -- "gak ada order baru" itu kondisi normal, sering kejadian tiap polling).
func (h *Handler) GetPending(c *gin.Context) {
	res := helpers.NewResponse()

	branchID, err := strconv.ParseInt(c.Param("branch_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("branch_id salah"))
		return
	}

	orders, err := GetPending(c, h.DB, branchID)
	if err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("gagal ambil data order"))
		return
	}

	c.JSON(http.StatusOK, res.Success().SetData(orders))
}

// Ack: POST /pos/mobile-order/ack/:order_number -- dipanggil POS SETELAH tr_order lokal
// berhasil keinsert (lihat catatan di Ack() service). Selalu balikin Success kalau query jalan
// (idempotent, no-op kalau udah pernah di-ack atau order_number gak ketemu -- POS gak perlu
// nge-handle 2 hasil beda buat retry).
func (h *Handler) Ack(c *gin.Context) {
	res := helpers.NewResponse()

	orderNumber := c.Param("order_number")
	if orderNumber == "" {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("order_number wajib diisi"))
		return
	}

	if err := Ack(c, h.DB, orderNumber); err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("gagal ack order"))
		return
	}

	c.JSON(http.StatusOK, res.Success())
}
