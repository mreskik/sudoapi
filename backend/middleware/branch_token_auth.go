package middleware

import (
	"APIANDORDER/backend/helpers"
	"APIANDORDER/backend/modules/apipos/sessiontracker"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type sessionLookup struct {
	Id       int64 `bun:"id"`
	BranchId int64 `bun:"branch_id"`
}

// BranchTokenAuth (2026-10-16, sistem session POS dinamis multi-device): validasi SEKARANG ke
// master_branch_pos_session (bukan master_branch.token statis lagi) -- cek token ADA dan
// flag_active = true. branch_id hasil lookup disimpen ke context (key "pos_branch_id") biar
// handler lain bisa cross-check terhadap :branch_id di path URL kalau perlu diperketat lagi
// nanti (belum dipakai di titik ini, baru infrastructure-nya).
//
// 2026-10-19: tiap request yang lolos sini juga nyatet "dipakai sekarang" (in-memory,
// sessiontracker.TouchLastUsed) dan metadata device (X-App-Version/X-OS-Version header +
// IP request) -- DB cuma kesentuh kalau metadata-nya BEDA dari yang terakhir diketahui
// (sessiontracker.ShouldUpdateMetadata), biar gak nambah 1 UPDATE query di SETIAP request.
// Lihat DOKUMENTASI API/MASTER/MASTER BRANCH POS SESSION.md (sudocore2) buat konteks lengkap.
func BranchTokenAuth(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := helpers.NewResponse()

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, res.SetCode(401).SetMessage("token tidak ditemukan"))
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		var session sessionLookup
		err := db.NewRaw("SELECT id, branch_id FROM master_branch_pos_session WHERE token = ? AND flag_active = true", token).Scan(c, &session)
		if err != nil {
			c.JSON(401, res.SetCode(401).SetMessage("token tidak valid"))
			c.Abort()
			return
		}

		c.Set("pos_branch_id", session.BranchId)

		sessiontracker.TouchLastUsed(session.Id)

		incoming := sessiontracker.Metadata{
			AppVersion: c.GetHeader("X-App-Version"),
			OSVersion:  c.GetHeader("X-OS-Version"),
			IPAddress:  c.ClientIP(),
		}
		if sessiontracker.ShouldUpdateMetadata(session.Id, incoming) {
			go func(id int64, meta sessiontracker.Metadata) {
				_, _ = db.NewUpdate().Table("master_branch_pos_session").
					Set("app_version = ?", meta.AppVersion).
					Set("os_version = ?", meta.OSVersion).
					Set("ip_address = ?", meta.IPAddress).
					Where("id = ?", id).
					Exec(context.Background())
			}(session.Id, incoming)
		}

		c.Next()
	}
}
