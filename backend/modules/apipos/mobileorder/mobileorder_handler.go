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
	// cross-check branch_id (path) vs branch_id hasil lookup token (2026-10-16, audit keamanan)
	// -- BranchTokenAuth cuma mastiin token VALID/AKTIF, TIDAK mastiin buat branch yang
	// direquest. Tanpa ini, device branch A bisa ganti :branch_id jadi branch manapun dan baca
	// order pending (data member, harga, no HP customer) milik branch lain.
	if err := helpers.MatchingBranchID(c, branchID); err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage(err.Error()))
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
//
// Cross-check branch_id (2026-10-16, audit keamanan): order_number sendiri gak kebawa
// branch_id di URL, jadi query balik dulu branch_id punya order ini (GetBranchIDByOrderNumber)
// sebelum ack -- tanpa ini device branch A bisa ack (silently hilangkan dari antrian) order
// milik branch manapun cuma dengan tau/nebak order_number-nya. order_number gak ketemu TETAP
// Success (konsisten idempotent yang udah ada, BUKAN dianggap percobaan sabotase) -- yang
// ditolak cuma kalau order-nya KETEMU dan branch_id-nya beda dari token.
func (h *Handler) Ack(c *gin.Context) {
	res := helpers.NewResponse()

	orderNumber := c.Param("order_number")
	if orderNumber == "" {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("order_number wajib diisi"))
		return
	}

	orderBranchID, err := GetBranchIDByOrderNumber(c, h.DB, orderNumber)
	if err == nil {
		if matchErr := helpers.MatchingBranchID(c, orderBranchID); matchErr != nil {
			c.JSON(http.StatusOK, res.GeneralError().SetMessage(matchErr.Error()))
			return
		}
	}
	// err != nil (order_number gak ketemu) -- sengaja DILEWATIN, lanjut ke Ack() di bawah yang
	// emang udah idempotent buat kasus itu (WHERE pulled_at IS NULL, 0 row affected = no-op).

	if err := Ack(c, h.DB, orderNumber); err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("gagal ack order"))
		return
	}

	c.JSON(http.StatusOK, res.Success())
}
