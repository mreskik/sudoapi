package member

import (
	"APIANDORDER/backend/helpers"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type MemberHandler struct {
	DB *bun.DB
}

func NewHandler(db *bun.DB) *MemberHandler {
	return &MemberHandler{DB: db}
}

// MemberByPhoneResult hasil live lookup master_member by phone_number -- gak ada company/branch
// scoping (member emang global di sistem ini, sama kayak query pull GetMasterMember), murni
// dicocokin lewat phone_number yang sekarang dijamin unik (lihat migration 086 di sudocore2).
type MemberByPhoneResult struct {
	ID             int    `bun:"id" json:"id"`
	MemberTypeID   *int   `bun:"member_type_id" json:"member_type_id"`
	MemberTypeName string `bun:"member_type_name" json:"member_type_name"`
	Code           string `bun:"code" json:"code"`
	Name           string `bun:"name" json:"name"`
	ContactName    string `bun:"contact_name" json:"contact_name"`
	Email          string `bun:"email" json:"email"`
	PhoneNumber    string `bun:"phone_number" json:"phone_number"`
}

// CheckByPhone: live lookup (bukan baca dari cache/sync manapun) -- dipanggil dari POS Laravel
// (Kiosk), token branch divalidasi middleware.BranchTokenAuth di router. phone_number gak
// ketemu bukan error, cukup balikin data null (frontend yang mutusin mau nampilin apa).
func (h *MemberHandler) CheckByPhone(c *gin.Context) {
	res := helpers.NewResponse()

	phoneNumber := c.Param("phone_number")
	if phoneNumber == "" {
		c.JSON(200, res.GeneralError().SetMessage("phone_number wajib diisi"))
		return
	}

	result := MemberByPhoneResult{}
	err := h.DB.NewRaw(`
	SELECT mm.id, mm.member_type_id, COALESCE(mmt.name, '') as member_type_name,
	mm.code, mm.name, COALESCE(mm.contact_name, '') as contact_name,
	COALESCE(mm.email, '') as email, COALESCE(mm.phone_number, '') as phone_number
	FROM master_member mm
	LEFT JOIN master_member_type mmt ON mmt.id = mm.member_type_id
	WHERE mm.phone_number = ? AND mm.is_active = true
	`, phoneNumber).Scan(c, &result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(200, res.Success().SetData(nil))
			return
		}
		c.JSON(200, res.GeneralError().SetMessage("gagal cek member"))
		return
	}

	c.JSON(200, res.Success().SetData(result))
}
