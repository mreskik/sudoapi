package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken bikin token random 256-bit (32 byte), di-encode hex jadi string 64 karakter --
// sama persis pola sudocore2 (backend/helpers/generatetoken.go), dipakai buat generate token
// baru per device pas install POS (master_branch_pos_session.token).
func GenerateToken() (string, error) {
	b := make([]byte, 32) // 256-bit
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
