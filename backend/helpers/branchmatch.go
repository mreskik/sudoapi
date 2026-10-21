package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// MatchingBranchID (2026-10-16, audit keamanan jalur token branch): ambil branch_id hasil
// lookup token (c.Get("pos_branch_id"), di-set middleware.BranchTokenAuth) dan pastikan cocok
// dengan branchID yang direquest (path/body) -- BranchTokenAuth cuma mastiin token itu VALID
// dan AKTIF, TIDAK mastiin token itu buat branch yang direquest. Tanpa cross-check ini, device
// branch A yang punya token sah bisa ganti branch_id di request (path atau body) jadi branch
// manapun dan baca/tulis data branch itu -- celah yang ditemukan di hampir semua endpoint
// /pos/sync/*, /pos/mobile-order/*, /pos/member-topup/*, /pos/ws/mobile-order/*, /pos/heartbeat/*
// (push/* dan sold_out udah dibenerin duluan, 2026-10-16 juga, pola yang sama).
func MatchingBranchID(c *gin.Context, branchID int64) error {
	posBranchId, exists := c.Get("pos_branch_id")
	if !exists {
		return errors.New("branch id tidak ditemukan dari token")
	}
	tokenBranchId, ok := posBranchId.(int64)
	if !ok || tokenBranchId != branchID {
		return errors.New("branch_id tidak cocok dengan token yang dipakai")
	}
	return nil
}
