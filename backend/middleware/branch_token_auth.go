package middleware

import (
	"APIANDORDER/backend/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

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

		var count int
		err := db.NewRaw("SELECT COUNT(*) FROM master_branch WHERE token = ?", token).Scan(c, &count)
		if err != nil || count == 0 {
			c.JSON(401, res.SetCode(401).SetMessage("token tidak valid"))
			c.Abort()
			return
		}

		c.Next()
	}
}
