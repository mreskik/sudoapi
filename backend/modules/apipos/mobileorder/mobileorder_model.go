package mobileorder

// PendingOrder: header mb_order (ERP, sudocore2) apa adanya -- field operasional POS
// (terminal/table_section/dayshift/waiter/order_queue) SENGAJA gak ada di sini, itu diisi POS
// sendiri saat processOrder() dari terminal worker yang udah di-resolve (resolveWorkerTerminal()),
// bukan dari payload ini. Lihat migration sudocore2 119_create_table_mb_order.sql.
type PendingOrder struct {
	OrderNumber string `bun:"order_number" json:"order_number"`
	BranchID    int64  `bun:"branch_id" json:"branch_id"`
	// MemberID: NULLABLE (2026-09-17, migration sudocore2 208) -- order QR Order itu TAMU, gak
	// ada member sama sekali. WAJIB *int64 (bukan int64) -- GetPending() narik banyak order
	// dalam 1 Scan, kalau field ini non-pointer, 1 baris QR (member_id NULL) bakal nge-gagalin
	// SELURUH batch pull branch itu, bukan cuma order QR-nya doang.
	MemberID *int64 `bun:"member_id" json:"member_id"`
	// MemberName: JOIN live ke master_member di query (bukan denormalisasi ke mb_order) --
	// selalu fresh (gak basi kalau member ganti nama), dan gak gantung POS udah nyinkron member
	// itu apa belum (lihat GetPending() di mobileorder_service.go). Kosong kalau member_id NULL
	// (order QR/tamu) ATAU member ke-hapus.
	MemberName string `bun:"member_name" json:"member_name"`
	// OrderName: mb_order.order_name apa adanya (migration sudocore2 210+211, 2026-09-17) --
	// SELALU NULL buat order member app (identitasnya dari MemberName di atas), SELALU keisi buat
	// order QR (identitas tamu, wajib diisi Create). Puller POS pakai `OrderName ?? MemberName`
	// buat isi tr_order.order_name -- 2 field ini gak pernah keisi bareng, salah satu pasti NULL.
	OrderName      *string `bun:"order_name" json:"order_name"`
	VisitPurposeID int64   `bun:"visit_purpose_id" json:"visit_purpose_id"`
	OrderType      string  `bun:"order_type" json:"order_type"`
	// OrderSource: channel asal order -- 'mobile' (member app) atau 'qr' (QR Order, belum ada
	// endpoint Create-nya per 2026-09-17, migration sudocore2 209) -- dipakai POS buat nulis
	// tr_order.order_source (MobileOrderPullServices::processOrder()), BUKAN cuma dilewatin.
	OrderSource         string  `bun:"order_source" json:"order_source"`
	Pax                 *int    `bun:"pax" json:"pax"`
	OrderFee            float64 `bun:"order_fee" json:"order_fee"`
	ServiceCharge       float64 `bun:"service_charge" json:"service_charge"`
	PlatformFee         float64 `bun:"platform_fee" json:"platform_fee"`
	DeliveryCost        float64 `bun:"delivery_cost" json:"delivery_cost"`
	SubTotal            float64 `bun:"sub_total" json:"sub_total"`
	TotalDiscount       float64 `bun:"total_discount" json:"total_discount"`
	TotalTax            float64 `bun:"total_tax" json:"total_tax"`
	TotalBilling        float64 `bun:"total_billing" json:"total_billing"`
	FlagInclusiveTax    bool    `bun:"flag_inclusive_tax" json:"flag_inclusive_tax"`
	CustomerPhoneNumber *string `bun:"customer_phone_number" json:"customer_phone_number"`
	PaymentNumber       *string `bun:"payment_number" json:"payment_number"`
	PaymentAt           *string `bun:"payment_at" json:"payment_at"`
	PaymentNotes        *string `bun:"payment_notes" json:"payment_notes"`
	CreatedAt           string  `bun:"created_at" json:"created_at"`

	Detail  []PendingOrderDetail `bun:"-" json:"detail"`
	Payment *PendingOrderPayment `bun:"-" json:"payment"`
}

// PendingOrderDetail: mb_order_detail, 1 baris per item (termasuk induk package -- baris
// package group sendiri juga masuk sini, isi-isinya di Package).
type PendingOrderDetail struct {
	ULID              string   `bun:"ulid" json:"ulid"`
	OrderNumber       string   `bun:"order_number" json:"-"`
	PricelistDetailID int64    `bun:"pricelist_detail_id" json:"pricelist_detail_id"`
	MenuID            int64    `bun:"menu_id" json:"menu_id"`
	CategoryID        *int64   `bun:"category_id" json:"category_id"`
	SubcategoryID     *int64   `bun:"subcategory_id" json:"subcategory_id"`
	Qty               int64    `bun:"qty" json:"qty"`
	FlagInclusiveTax  bool     `bun:"flag_inclusive_tax" json:"flag_inclusive_tax"`
	Price             float64  `bun:"price" json:"price"`
	TaxID             *int64   `bun:"tax_id" json:"tax_id"`
	TaxType           *string  `bun:"tax_type" json:"tax_type"`
	TaxRate           float64  `bun:"tax_rate" json:"tax_rate"`
	TaxAmount         float64  `bun:"tax_amount" json:"tax_amount"`
	Dpp               *float64 `bun:"dpp" json:"dpp"`
	NetDpp            *float64 `bun:"net_dpp" json:"net_dpp"`
	PromoID           *int64   `bun:"promo_id" json:"promo_id"`
	DiscountPercent   float64  `bun:"discount_percent" json:"discount_percent"`
	DiscountAmount    float64  `bun:"discount_amount" json:"discount_amount"`
	Total             float64  `bun:"total" json:"total"`
	Notes             *string  `bun:"notes" json:"notes"`

	Package []PendingOrderDetailPackage `bun:"-" json:"package"`
}

// PendingOrderDetailPackage: mb_order_detail_package, nested di dalam detail induknya lewat
// mb_order_detail_ulid.
type PendingOrderDetailPackage struct {
	ULID              string   `bun:"ulid" json:"ulid"`
	MbOrderDetailULID string   `bun:"mb_order_detail_ulid" json:"-"`
	MenuPackageID     int64    `bun:"menu_package_id" json:"menu_package_id"`
	MenuID            int64    `bun:"menu_id" json:"menu_id"`
	CategoryID        *int64   `bun:"category_id" json:"category_id"`
	SubcategoryID     *int64   `bun:"subcategory_id" json:"subcategory_id"`
	Qty               int64    `bun:"qty" json:"qty"`
	FlagInclusiveTax  bool     `bun:"flag_inclusive_tax" json:"flag_inclusive_tax"`
	Price             float64  `bun:"price" json:"price"`
	TaxID             *int64   `bun:"tax_id" json:"tax_id"`
	TaxType           *string  `bun:"tax_type" json:"tax_type"`
	TaxRate           float64  `bun:"tax_rate" json:"tax_rate"`
	TaxAmount         float64  `bun:"tax_amount" json:"tax_amount"`
	Dpp               *float64 `bun:"dpp" json:"dpp"`
	NetDpp            *float64 `bun:"net_dpp" json:"net_dpp"`
	PromoID           *int64   `bun:"promo_id" json:"promo_id"`
	DiscountPercent   float64  `bun:"discount_percent" json:"discount_percent"`
	DiscountAmount    float64  `bun:"discount_amount" json:"discount_amount"`
	Total             float64  `bun:"total" json:"total"`
	Notes             *string  `bun:"notes" json:"notes"`
}

// PendingOrderPayment: mb_order_payment, 1 baris FINAL per order (settled payment, gak ada
// split payment saat ini -- lihat catatan migration 119). FK-nya ke mb_order.payment_number
// (BUKAN order_number -- direname migration 120), makanya join query pakai payment_number.
type PendingOrderPayment struct {
	ULID                  string  `bun:"ulid" json:"ulid"`
	PaymentNumber         string  `bun:"payment_number" json:"-"`
	PaymentMethodID       int64   `bun:"payment_method_id" json:"payment_method_id"`
	PaymentAmount         float64 `bun:"payment_amount" json:"payment_amount"`
	VoucherCode           *string `bun:"voucher_code" json:"voucher_code"`
	PercentFromBilling    float64 `bun:"percent_from_billing" json:"percent_from_billing"`
	PaymentGatewayOrderID *string `bun:"payment_gateway_order_id" json:"payment_gateway_order_id"`
}
