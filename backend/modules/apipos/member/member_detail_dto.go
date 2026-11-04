package member

import "time"

// MemberDetailResult: header detail member -- kolom SAMA PERSIS query View() di sudocore2
// (backend/modules/master/member/member_services.go) ditambah tier_name hasil JOIN ke
// master_member_tier_setting_detail (level loyalty, bukan cuma angka tier_level mentah). Point
// TIDAK dari query ini -- diisi manual di handler lewat query snapshot ledger terpisah (sama
// alasan sudocore2 misahin View()/GetPoint() jadi 2 fungsi, concern beda). TierLevel/TierName
// JANGAN dikasih json tag di sini -- dibungkus nested lewat MemberTierResult (key "tier") di
// handler, bukan flat campur sama field header lain (2026-10-30, permintaan eksplisit).
type MemberDetailResult struct {
	ID          int     `bun:"id" json:"id"`
	Code        string  `bun:"code" json:"code"`
	Name        string  `bun:"name" json:"name"`
	ContactName string  `bun:"contact_name" json:"contact_name"`
	Email       string  `bun:"email" json:"email"`
	PhoneNumber string  `bun:"phone_number" json:"phone_number"`
	BirthDate   *string `bun:"birth_date" json:"birth_date"`
	JoinDate    string  `bun:"join_date" json:"join_date"`
	Note        string  `bun:"note" json:"note"`
	TierLevel   int     `bun:"tier_level" json:"-"`
	TierName    *string `bun:"tier_name" json:"-"`
	Point       string  `bun:"-" json:"point"`
}

// MemberTierResult: level loyalty member, dibungkus nested (key "tier") -- dipisah dari field
// header biar gak campur flat (2026-10-30, permintaan eksplisit, sama alasan last_visit).
type MemberTierResult struct {
	Level int     `json:"level"`
	Name  *string `json:"name"`
}

// MemberLastOrderItemResult/MemberLastOrderPackageResult: breakdown item + nested package per
// item -- MIRROR PERSIS query items/packages di GetOrderHistoryDetail sudocore2 (ULID dipakai
// cuma buat mapping package ke item yang bener di tahap query, TIDAK diekspos ke JSON akhir).
type MemberLastOrderPackageResult struct {
	TrOrderDetailULID string  `bun:"tr_order_detail_ulid" json:"-"`
	MenuName          string  `bun:"menu_name" json:"menu_name"`
	PromoName         *string `bun:"promo_name" json:"promo_name"`
	Qty               string  `bun:"qty" json:"qty"`
	Dpp               string  `bun:"dpp" json:"dpp"`
	TaxRate           string  `bun:"tax_rate" json:"tax_rate"`
	DiscountAmount    string  `bun:"discount_amount" json:"discount_amount"`
	Total             string  `bun:"total" json:"total"`
	Notes             *string `bun:"notes" json:"notes"`
}

type MemberLastOrderItemResult struct {
	ULID           string                         `bun:"ulid" json:"-"`
	MenuName       string                         `bun:"menu_name" json:"menu_name"`
	PromoName      *string                        `bun:"promo_name" json:"promo_name"`
	Qty            string                         `bun:"qty" json:"qty"`
	Dpp            string                         `bun:"dpp" json:"dpp"`
	TaxRate        string                         `bun:"tax_rate" json:"tax_rate"`
	DiscountAmount string                         `bun:"discount_amount" json:"discount_amount"`
	Total          string                         `bun:"total" json:"total"`
	Notes          *string                        `bun:"notes" json:"notes"`
	Packages       []MemberLastOrderPackageResult `json:"packages"`
}

// MemberLastOrderPaymentResult: 1 baris metode bayar yang dipakai order ini (SPLIT PAYMENT bisa
// lebih dari 1 baris) -- adopsi query payment di posprint_service.go::ResolveReceiptData
// (pos_order.payment_number -> pos_order_payment -> master_payment_method).
type MemberLastOrderPaymentResult struct {
	PaymentMethodName string `bun:"payment_method_name" json:"payment_method_name"`
	PaymentAmount     string `bun:"payment_amount" json:"payment_amount"`
}

// MemberLastOrderResult: bentuk akhir "last_order" -- field header order FLAT langsung di sini
// (BUKAN dibungkus sub-object "header", 2026-10-30 permintaan eksplisit: keluarin field header,
// ga usah pakai wrapper) + item_list[] (nested packages[] per item) + payment_method_list[].
// Semua field header pointer -- null total kalau member belum pernah order 'paid' sama sekali
// (ItemList/PaymentMethodList tetap slice kosong, bukan nil, di kondisi itu).
type MemberLastOrderResult struct {
	OrderNumber       *string                         `bun:"order_number" json:"order_number"`
	BranchID          *int                            `bun:"branch_id" json:"branch_id"`
	BranchName        *string                         `bun:"branch_name" json:"branch_name"`
	OrderOut          *time.Time                      `bun:"order_out" json:"order_out"`
	Status            *string                         `bun:"status" json:"status"`
	OrderType         *string                         `bun:"order_type" json:"order_type"`
	OrderSource       *string                         `bun:"order_source" json:"order_source"`
	SubTotal          *string                         `bun:"sub_total" json:"sub_total"`
	TotalDiscount     *string                         `bun:"total_discount" json:"total_discount"`
	TotalTax          *string                         `bun:"total_tax" json:"total_tax"`
	TotalBilling      *string                         `bun:"total_billing" json:"total_billing"`
	ItemList          []MemberLastOrderItemResult     `json:"item_list"`
	PaymentMethodList []MemberLastOrderPaymentResult  `json:"payment_method_list"`
}

// MemberBestSellerResult: 1 baris agregasi menu -- adopsi query GetBestSeller() sudocore2
// (member_services.go) + LIMIT 5 ditambahkan di SQL (beda dari sudocore2 yang sengaja tanpa
// limit -- di sini FE cuma butuh top 5 pasti).
type MemberBestSellerResult struct {
	ItemName        string    `bun:"item_name" json:"item_name"`
	Category        *string   `bun:"category" json:"category"`
	SubCategory     *string   `bun:"sub_category" json:"sub_category"`
	Qty             string    `bun:"qty" json:"qty"`
	LastPurchasedAt time.Time `bun:"last_purchased_at" json:"last_purchased_at"`
}

// MemberDetailResponse: bentuk akhir yang dikirim ke POS -- gabungan query (header+point, tier,
// last order lengkap, best seller) jadi 1 object, dibangun di handler. Tier/LastOrder SENGAJA
// nested (bukan embed flat) -- lihat komentar MemberTierResult/MemberLastOrderResult.
type MemberDetailResponse struct {
	MemberDetailResult
	Tier       MemberTierResult         `json:"tier"`
	LastOrder  MemberLastOrderResult    `json:"last_order"`
	BestSeller []MemberBestSellerResult `json:"best_seller"`
}
