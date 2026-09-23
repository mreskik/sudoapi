package membertopup

// CreateTopupRequestDTO -- body buat POST /pos/member-topup/:branch_id.
// PhoneNumber -- caller (POS/Kiosk/Mobile) cuma pegang nomor HP customer, BUKAN member_id --
// resolve ke member_id dilakuin di sini (APIANDORDER, konek langsung ke master_member), sama
// query yang dipakai member.CheckByPhone(), biar caller gak perlu manggil 2 endpoint terpisah
// (by-phone dulu baru topup) buat 1 aksi.
// PaymentMethodID kosong/null = tunai (langsung 'paid' saat itu juga). Keisi = lewat payment
// gateway (status mulai 'pending', nunggu settlement -- lihat CheckStatus()) --
// master_payment_method.id (PUSAT, sudocore2), BUKAN id lokal per-caller. payment_gateway_code
// di-resolve SERVER-SIDE dari sini (migration 221, 2026-09-22 -- SEBELUMNYA caller kirim
// payment_gateway_code langsung, diganti biar payment_method_id yang tersimpan di
// member_topup_online JELAS & GAK AMBIGU, dipakai memberbalancejurnal resolve akun COA).
type CreateTopupRequestDTO struct {
	PhoneNumber     string  `json:"phone_number"`
	Amount          float64 `json:"amount"`
	Source          string  `json:"source"` // 'pos', 'kiosk', 'mobile'
	PaymentMethodID *int64  `json:"payment_method_id"`
	Notes           *string `json:"notes"`
	// TerminalID -- opsional (source 'pos'/'mobile' mungkin gak punya konsep terminal), dari
	// Kiosk selalu keisi (device udah tau ID-nya sendiri, sama pola SaveOrder()). Disimpan ke
	// member_topup_online.terminal_id + member_balance_ledger.terminal_id, di-echo balik di
	// CheckStatus buat dipakai POS resolve receipt_station pas print struk.
	TerminalID *int64 `json:"terminal_id"`
}

// CreateTopupResponseDTO -- jalur tunai balik dengan BalanceAfter langsung keisi (saldo baru
// udah final). Jalur gateway balik dengan data QR, BalanceAfter masih nil (nunggu settlement).
type CreateTopupResponseDTO struct {
	ReferenceNumber string  `json:"reference_number"`
	Status          string  `json:"status"`
	BalanceAfter    *string `json:"balance_after,omitempty"`
	QRString        *string `json:"qr_string,omitempty"`
	QRURL           *string `json:"qr_url,omitempty"`
	ExpiredAt       *string `json:"expired_at,omitempty"`
}

// CheckTopupStatusResponseDTO -- respons buat GET check-status. Field TerminalID/Amount/
// MemberName/PaymentGatewayCode/PaidAt cuma keisi kalau Status == "paid" -- dipakai POS
// (Laravel, stateless, gak nyimpen ulang data topup) buat resolve receipt_station + isi struk,
// gak perlu query balik ke sini pas print.
type CheckTopupStatusResponseDTO struct {
	ReferenceNumber    string  `json:"reference_number"`
	Status             string  `json:"status"`
	BalanceAfter       *string `json:"balance_after,omitempty"`
	TerminalID         *int64  `json:"terminal_id,omitempty"`
	Amount             *string `json:"amount,omitempty"`
	MemberName         *string `json:"member_name,omitempty"`
	MemberPhoneNumber  *string `json:"member_phone_number,omitempty"`
	PaymentGatewayCode *string `json:"payment_gateway_code,omitempty"`
	PaidAt             *string `json:"paid_at,omitempty"`
}

// paymentGatewayCreateQrisResponse -- bentuk response {code, message, data} dari service
// payment POST /payment-gateway/qris (lihat payment/backend/modules/paymentgateway/paymentgateway_dto.go).
type paymentGatewayCreateQrisResponse struct {
	Code int `json:"code"`
	Data struct {
		OrderID        string  `json:"order_id"`
		Status         string  `json:"status"`
		VendorQRString *string `json:"vendor_qr_string"`
		VendorQRURL    *string `json:"vendor_qr_url"`
		ExpiredAt      *string `json:"expired_at"`
	} `json:"data"`
	Message string `json:"message"`
}

// paymentGatewayStatusResponse -- bentuk response GET /payment-gateway/:order_id.
type paymentGatewayStatusResponse struct {
	Code int `json:"code"`
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}
