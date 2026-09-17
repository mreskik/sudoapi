package heartbeat

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

// Ping: POST /pos/heartbeat/:branch_id -- dipanggil worker POS tiap 30 detik (command
// `heartbeat:send`, belum dibangun). Upsert branch_heartbeat.last_ping_at = now(), dibaca
// sudomobile buat barrier order + flag_status_store_open (belum dibangun juga -- item ini
// baru nulis sisi APIANDORDER-nya).
func (h *Handler) Ping(c *gin.Context) {
	res := helpers.NewResponse()

	branchID, err := strconv.ParseInt(c.Param("branch_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("branch_id salah"))
		return
	}

	_, err = h.DB.NewRaw(`
		INSERT INTO branch_heartbeat (branch_id, last_ping_at) VALUES (?, now())
		ON CONFLICT (branch_id) DO UPDATE SET last_ping_at = now()
	`, branchID).Exec(c)
	if err != nil {
		c.JSON(http.StatusOK, res.GeneralError().SetMessage("gagal simpan heartbeat"))
		return
	}

	c.JSON(http.StatusOK, res.Success())
}
