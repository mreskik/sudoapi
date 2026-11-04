package member

import (
	"APIANDORDER/backend/helpers"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

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

// GetDetail: detail member online buat halaman Member Detail POS -- gabungan header (incl.
// tier), point snapshot, last order LENGKAP (header+items+packages), dan top 5 best seller.
// member_id gak ketemu -- SAMA pola CheckByPhone, balikin data null bukan error (FE yang
// mutusin mau nampilin apa).
func (h *MemberHandler) GetDetail(c *gin.Context) {
	res := helpers.NewResponse()

	memberId, err := strconv.Atoi(c.Param("member_id"))
	if err != nil || memberId <= 0 {
		c.JSON(200, res.GeneralError().SetMessage("member_id tidak valid"))
		return
	}

	header, err := h.getMemberDetailHeader(c, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(200, res.Success().SetData(nil))
			return
		}
		c.JSON(200, res.GeneralError().SetMessage("gagal mengambil data member"))
		return
	}

	point, err := h.getMemberPoint(c, memberId)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal mengambil point member"))
		return
	}
	header.Point = point

	lastOrder, err := h.getMemberLastOrder(c, memberId)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal mengambil riwayat order member"))
		return
	}

	bestSeller, err := h.getMemberBestSeller(c, memberId)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal mengambil best seller member"))
		return
	}

	result := MemberDetailResponse{
		MemberDetailResult: header,
		Tier:               MemberTierResult{Level: header.TierLevel, Name: header.TierName},
		LastOrder:          lastOrder,
		BestSeller:         bestSeller,
	}

	c.JSON(200, res.Success().SetData(result))
}

// CreateMemberRequestDTO: body POST /pos/member/:branch_id -- 6 field yang dibutuhkan kasir POS
// buat daftarin member baru langsung dari kasir (2026-11-04). code/name/phone_number WAJIB,
// gender/birth_date/description opsional. branch_id di path CUMA buat konsistensi URL dengan
// endpoint member lain (/pos/member/:branch_id/...) -- TIDAK dipakai cross-check (beda dari
// membertopup yang pakai helpers.MatchingBranchID), karena master_member TIDAK punya kolom
// branch_id -- member itu entitas GLOBAL lintas branch, bukan terikat 1 branch tertentu.
type CreateMemberRequestDTO struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Gender      *string `json:"gender"`
	PhoneNumber string  `json:"phone_number"`
	BirthDate   *string `json:"birth_date"`
	Description *string `json:"description"`
}

// defaultMemberTypeIDCustomer: id master_member_type "Customer" -- HARDCODE, BUKAN lookup by
// nama (2026-11-04, keputusan eksplisit user: "Customer" itu data DEFAULT bawaan sistem, id-nya
// gak akan berubah, sementara NAMANYA berisiko diganti ke depan mis. jadi "Cust" -- lookup by
// nama malah rawan salah kalau itu terjadi, hardcode by id lebih aman buat kasus spesifik ini).
// Dikonfirmasi lewat GET /backend/data/master/member-type di DB dev: {"id":3,"name":"Customer"}.
const defaultMemberTypeIDCustomer = 3

// CreateMember: INSERT LANGSUNG ke master_member (BUKAN proxy HTTP ke sudocore2, konsisten pola
// CheckByPhone/GetDetail di atas yang semua query DB langsung -- APIANDORDER konek ke Postgres
// yang sama dengan sudocore2). member_type_id DEFAULT ke "Customer" (lihat defaultMemberTypeIDCustomer
// di atas) -- member yang didaftarkan kasir langsung dari POS secara natural adalah customer baru.
// contact_name/email/profile_photo_src/shift_id/attendance_branch_setting/pin TIDAK diisi (di
// luar scope 6 field yang diminta) -- is_active default true, tier_level default 1 (default kolom DB).
//
// gender divalidasi lewat CHECK constraint DB (master_member_gender_check, 'male'/'female') --
// TIDAK direplikasi manual di sini (sama prinsip sudocore2: validasi di DB, bukan app-level),
// cukup translate pesan errornya kalau constraint itu yang nge-reject. birth_date DIVALIDASI
// app-level dulu (format YYYY-MM-DD) SEBELUM insert, biar error-nya jelas (bukan error Postgres
// "invalid input syntax for type date" yang mentah).
func (h *MemberHandler) CreateMember(c *gin.Context) {
	res := helpers.NewResponse()

	req := CreateMemberRequestDTO{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, res.GeneralError().SetMessage("invalid request format"))
		return
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	phoneNumber := strings.TrimSpace(req.PhoneNumber)
	if code == "" {
		c.JSON(200, res.GeneralError().SetMessage("code is required"))
		return
	}
	if name == "" {
		c.JSON(200, res.GeneralError().SetMessage("name is required"))
		return
	}
	if phoneNumber == "" {
		c.JSON(200, res.GeneralError().SetMessage("phone_number is required"))
		return
	}

	var birthDate *string
	if req.BirthDate != nil && strings.TrimSpace(*req.BirthDate) != "" {
		trimmed := strings.TrimSpace(*req.BirthDate)
		if _, err := time.Parse("2006-01-02", trimmed); err != nil {
			c.JSON(200, res.GeneralError().SetMessage("invalid date format, must be YYYY-MM-DD"))
			return
		}
		birthDate = &trimmed
	}

	var gender *string
	if req.Gender != nil && strings.TrimSpace(*req.Gender) != "" {
		trimmed := strings.TrimSpace(*req.Gender)
		gender = &trimmed
	}

	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		description = &trimmed
	}

	var newId int
	err := h.DB.NewRaw(`
	INSERT INTO master_member (code, name, phone_number, gender, birth_date, description, member_type_id, is_active, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, true, now())
	RETURNING id
	`, code, name, phoneNumber, gender, birthDate, description, defaultMemberTypeIDCustomer).Scan(c, &newId)
	if err != nil {
		if strings.Contains(err.Error(), "master_member_gender_check") {
			c.JSON(200, res.GeneralError().SetMessage("gender must be male or female"))
			return
		}
		c.JSON(200, res.GeneralError().SetMessage("failed to create new member"))
		return
	}

	c.JSON(200, res.Success().SetData(gin.H{"id": newId}))
}
