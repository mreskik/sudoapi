package mobilenotify

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// worker POS (PHP, textalk/websocket) bukan browser -- gak ada isu CORS/origin beneran di
	// sini, tapi gorilla/websocket tetap wajib CheckOrigin dijawab eksplisit (default-nya nolak
	// origin != host).
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// Serve: GET /pos/ws/mobile-order/:branch_id (di belakang middleware.BranchTokenAuth, sama kayak
// syncRouter/enddayRouter -- worker POS pegang branch token yang SAMA dipakai buat sync biasa,
// gak ada mekanisme auth baru).
//
// Upgrade ke WS, daftarin diri ke hub, TAHAN koneksi (ReadMessage loop, isi pesannya gak
// dipedulikan) sampai client disconnect atau proses mati -- broadcast dari listener.go yang
// nulis ke koneksi ini dari goroutine LAIN (lihat Hub.Broadcast).
func (h *Handler) Serve(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 100, "message": "branch_id salah"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.hub.Register(branchID, conn)
	defer h.hub.Unregister(branchID, conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
