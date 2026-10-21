package sessiontracker

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterInternalRoutes -- endpoint ANTAR-SERVICE (sudocore2 -> APIANDORDER), BUKAN buat
// device POS. Dilindungi shared secret (header X-Internal-Secret, env INTERNAL_API_SECRET --
// harus sama nilainya di .env APIANDORDER & sudocore2), bukan BranchTokenAuth (device POS gak
// pernah perlu manggil ini). Dipakai sudocore2 buat enrich last_used_at di endpoint monitoring
// /master/branch-pos-session (lihat DOKUMENTASI API/MASTER/MASTER BRANCH POS SESSION.md).
func RegisterInternalRoutes(router gin.IRouter) {
	internal := router.Group("/internal", requireInternalSecret)
	internal.GET("/branch-pos-session/last-used", getLastUsedBatchHandler)
}

func requireInternalSecret(c *gin.Context) {
	secret := os.Getenv("INTERNAL_API_SECRET")
	if secret == "" || c.GetHeader("X-Internal-Secret") != secret {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

// getLastUsedBatchHandler -- GET /internal/branch-pos-session/last-used?ids=3,9,12 -- batch
// biar sudocore2 List (banyak baris) gak perlu panggil endpoint ini 1 per 1 row.
func getLastUsedBatchHandler(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "ids wajib diisi"})
		return
	}

	parts := strings.Split(idsParam, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	result := GetLastUsedBatch(ids)

	data := make(map[string]*string, len(result))
	for id, t := range result {
		key := strconv.FormatInt(id, 10)
		if t == nil {
			data[key] = nil
			continue
		}
		formatted := t.Format("2006-01-02T15:04:05Z07:00")
		data[key] = &formatted
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}
