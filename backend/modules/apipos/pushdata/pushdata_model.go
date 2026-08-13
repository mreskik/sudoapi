package pushdata

import (
	"time"

	"github.com/uptrace/bun"
)

// transaksi
type PosOrderModel struct {
	bun.BaseModel `bun:"table:pos_order"`

	OrderNumber   string  `bun:"order_number,pk" json:"order_number"`
	PaymentNumber *string `bun:"payment_number" json:"payment_number"`

	BranchID   int64 `bun:"branch_id" json:"branch_id"`
	TerminalID int32 `bun:"terminal_id" json:"terminal_id"`

	OrderName   string `bun:"order_name" json:"order_name"`
	OrderType   string `bun:"order_type" json:"order_type"`
	OrderSource string `bun:"order_source" json:"order_source"`

	TableSectionID *int32 `bun:"table_section_id" json:"table_section_id"`
	TableID        *int32 `bun:"table_id" json:"table_id"`

	TotalBatch *int32 `bun:"total_batch" json:"total_batch"`

	OrderDate time.Time `bun:"order_date" json:"order_date"`

	OrderQueue *int32 `bun:"order_queue" json:"order_queue"`

	OrderIn  *time.Time `bun:"order_in" json:"order_in"`
	OrderOut *time.Time `bun:"order_out" json:"order_out"`

	MemberID       *int32 `bun:"member_id" json:"member_id"`
	VisitPurposeID *int32 `bun:"visit_purpose_id" json:"visit_purpose_id"`

	Pax *int32 `bun:"pax" json:"pax"`

	Status string `bun:"status" json:"status"`

	WaiterName string `bun:"waiter_name" json:"waiter_name"`
	SenderName string `bun:"sender_name" json:"sender_name"`

	CustomerPhoneNumber *string `bun:"customer_phone_number" json:"customer_phone_number"`
	ChasierName         *string `bun:"chasier_name" json:"chasier_name"`

	DeliveryCost  string `bun:"delivery_cost,type:numeric(20,2)" json:"delivery_cost"`
	OrderFee      string `bun:"order_fee,type:numeric(20,2)" json:"order_fee"`
	ServiceCharge string `bun:"service_charge,type:numeric(20,2)" json:"service_charge"`
	PlatformFee   string `bun:"platform_fee,type:numeric(20,2)" json:"platform_fee"`

	TotalItem *int64 `bun:"total_item" json:"total_item"`

	SubTotal      string `bun:"sub_total,type:numeric(20,2)" json:"sub_total"`
	TotalDiscount string `bun:"total_discount,type:numeric(20,2)" json:"total_discount"`
	TotalTax      string `bun:"total_tax,type:numeric(20,2)" json:"total_tax"`
	TotalBilling  string `bun:"total_billing,type:numeric(20,2)" json:"total_billing"`

	FlagInclusiveTax *bool `bun:"flag_inclusive_tax" json:"flag_inclusive_tax"`

	CancelNotes *string    `bun:"cancel_notes" json:"cancel_notes"`
	CancelAt    *time.Time `bun:"cancel_at" json:"cancel_at"`
	CancelBy    *int32     `bun:"cancel_by" json:"cancel_by"`

	CancelPrintKe *int32 `bun:"cancel_print_ke" json:"cancel_print_ke"`

	VoidNotes *string    `bun:"void_notes" json:"void_notes"`
	VoidAt    *time.Time `bun:"void_at" json:"void_at"`
	VoidBy    *int32     `bun:"void_by" json:"void_by"`

	VoidPrintKe *int32 `bun:"void_print_ke" json:"void_print_ke"`

	PaymentAt *time.Time `bun:"payment_at" json:"payment_at"`
	PaymentBy *int32     `bun:"payment_by" json:"payment_by"`

	PrintKe *int32 `bun:"print_ke" json:"print_ke"`

	PaymentNotes *string `bun:"payment_notes" json:"payment_notes"`

	MovedAt *time.Time `bun:"moved_at" json:"moved_at"`
	MovedBy *int32     `bun:"moved_by" json:"moved_by"`

	CreatedAt *time.Time `bun:"created_at" json:"created_at"`
	CreatedBy *int32     `bun:"created_by" json:"created_by"`

	UpdatedAt *time.Time `bun:"updated_at" json:"updated_at"`
	UpdatedBy *int32     `bun:"updated_by" json:"updated_by"`

	SyncAt *time.Time `bun:"sync_at" json:"sync_at"`

	DayshiftULID *string `bun:"dayshift_ulid" json:"dayshift_ulid"`

	EndDayAt  *time.Time `bun:"endday_at"`
	CompanyId *int       `bun:"company_id"`
}

type PosOrderDetailModel struct {
	bun.BaseModel `bun:"table:pos_order_detail"`

	ULID string `bun:"ulid,pk" json:"ulid"`

	OrderNumber string `bun:"order_number" json:"order_number"`

	PricelistDetailID int64 `bun:"pricelist_detail_id" json:"pricelist_detail_id"`
	MenuID            int64 `bun:"menu_id" json:"menu_id"`

	CategoryID    *int64 `bun:"category_id" json:"category_id"`
	SubcategoryID *int64 `bun:"subcategory_id" json:"subcategory_id"`

	Qty int64 `bun:"qty" json:"qty"`

	FlagInclusiveTax bool `bun:"flag_inclusive_tax" json:"flag_inclusive_tax"`

	BasePrice string `bun:"base_price,type:numeric(20,2)" json:"base_price"`

	TaxID   *int64  `bun:"tax_id" json:"tax_id"`
	TaxType *string `bun:"tax_type" json:"tax_type"`

	TaxRate  string `bun:"tax_rate,type:numeric(20,2)" json:"tax_rate"`
	TaxValue string `bun:"tax_value,type:numeric(20,2)" json:"tax_value"`

	PromoID         *int64 `bun:"promo_id" json:"promo_id"`
	IsFreeItemPromo bool   `bun:"is_free_item_promo" json:"is_free_item_promo"`

	DiscountRate  string `bun:"discount_rate,type:numeric(20,2)" json:"discount_rate"`
	DiscountValue string `bun:"discount_value,type:numeric(20,2)" json:"discount_value"`
	AfterDiscount string `bun:"after_discount,type:numeric(20,2)" json:"after_discount"`
	Dpp           string `bun:"dpp,type:numeric(20,2)" json:"dpp"`

	Total string `bun:"total,type:numeric(20,2)" json:"total"`

	Notes *string `bun:"notes" json:"notes"`

	Batch int32 `bun:"batch" json:"batch"`

	DonePrint bool `bun:"done_print" json:"done_print"`

	PrintKe int32 `bun:"print_ke" json:"print_ke"`

	CancelNotes *string    `bun:"cancel_notes" json:"cancel_notes"`
	CancelAt    *time.Time `bun:"cancel_at" json:"cancel_at"`

	CreatedAt *time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at" json:"updated_at"`

	CreatedBy *int64 `bun:"created_by" json:"created_by"`
	UpdatedBy *int64 `bun:"updated_by" json:"updated_by"`

	SyncAt *time.Time `bun:"sync_at" json:"sync_at"`
}

type PosOrderDetailPackageModel struct {
	bun.BaseModel `bun:"table:pos_order_detail_package"`

	ULID string `bun:"ulid,pk" json:"ulid"`

	TrOrderDetailULID string `bun:"tr_order_detail_ulid" json:"tr_order_detail_ulid"`

	MenuPackageID int64 `bun:"menu_package_id" json:"menu_package_id"`

	MenuID int64 `bun:"menu_id" json:"menu_id"`

	CategoryID    *int64 `bun:"category_id" json:"category_id"`
	SubcategoryID *int64 `bun:"subcategory_id" json:"subcategory_id"`

	Qty int64 `bun:"qty" json:"qty"`

	FlagInclusiveTax bool `bun:"flag_inclusive_tax" json:"flag_inclusive_tax"`

	BasePrice string `bun:"base_price,type:numeric(20,2)" json:"base_price"`

	TaxID   *int64  `bun:"tax_id" json:"tax_id"`
	TaxType *string `bun:"tax_type" json:"tax_type"`

	TaxRate  string `bun:"tax_rate,type:numeric(20,2)" json:"tax_rate"`
	TaxValue string `bun:"tax_value,type:numeric(20,2)" json:"tax_value"`

	PromoID *int64 `bun:"promo_id" json:"promo_id"`

	DiscountRate  string `bun:"discount_rate,type:numeric(20,2)" json:"discount_rate"`
	DiscountValue string `bun:"discount_value,type:numeric(20,2)" json:"discount_value"`
	AfterDiscount string `bun:"after_discount,type:numeric(20,2)" json:"after_discount"`
	Dpp           string `bun:"dpp,type:numeric(20,2)" json:"dpp"`

	Total string `bun:"total,type:numeric(20,2)" json:"total"`

	Notes string `bun:"notes" json:"notes"`

	SyncAt *time.Time `bun:"sync_at" json:"sync_at"`
}

type PosOrderPaymentModel struct {
	bun.BaseModel `bun:"table:pos_order_payment"`

	ULID string `bun:"ulid,pk" json:"ulid"`

	PaymentNumber string `bun:"payment_number" json:"payment_number"`

	PaymentMethodID int64 `bun:"payment_method_id" json:"payment_method_id"`

	PaymentAmount string `bun:"payment_amount,type:numeric(20,2)" json:"payment_amount"`

	VoucherCode *string `bun:"voucher_code" json:"voucher_code"`

	CardNumber *string `bun:"card_number" json:"card_number"`
	BankName   *string `bun:"bank_name" json:"bank_name"`

	AccountName      *string `bun:"account_name" json:"account_name"`
	VerificationCode *string `bun:"verification_code" json:"verification_code"`

	CreatedAt *time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at" json:"updated_at"`

	CreatedBy *int64 `bun:"created_by" json:"created_by"`
	UpdatedBy *int64 `bun:"updated_by" json:"updated_by"`

	SyncAt *time.Time `bun:"sync_at" json:"sync_at"`
}

// dayshift
type PosDayShiftModel struct {
	bun.BaseModel `bun:"table:pos_dayshift"`

	ULID string `bun:"ulid,pk" json:"ulid"`

	BranchID int64 `bun:"branch_id" json:"branch_id"`

	DayinTime   *time.Time `bun:"dayin_time" json:"dayin_time"`
	DayinTotal  string     `bun:"dayin_total,type:numeric(20,2)" json:"dayin_total"`
	DayinUserID *int64     `bun:"dayin_user_id" json:"dayin_user_id"`

	SystemCashReceived *string `bun:"system_cash_received,type:numeric(20,2)" json:"system_cash_received"`

	DayoutTime   *time.Time `bun:"dayout_time" json:"dayout_time"`
	DayoutTotal  *string    `bun:"dayout_total,type:numeric(20,2)" json:"dayout_total"`
	DayoutUserID *int64     `bun:"dayout_user_id" json:"dayout_user_id"`

	DayoutNotes *string `bun:"dayout_notes" json:"dayout_notes"`

	SyncAt *time.Time `bun:"sync_at" json:"sync_at"`
}

type PosDayShiftDetailModel struct {
	bun.BaseModel `bun:"table:pos_dayshift_detail"`

	ULID         string     `bun:"ulid,pk" json:"ulid"`
	DayShiftULID string     `bun:"dayshift_ulid" json:"dayshift_ulid"`
	ShiftTime    time.Time  `bun:"shift_time" json:"shift_time"`
	ShiftNumber  int64      `bun:"shift_number" json:"shift_number"`
	ShiftUserID  int64      `bun:"shift_user_id" json:"shift_user_id"`
	SyncAt       *time.Time `bun:"sync_at" json:"sync_at"`
}
