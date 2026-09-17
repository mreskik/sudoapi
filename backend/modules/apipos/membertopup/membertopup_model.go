package membertopup

import (
	"time"

	"github.com/uptrace/bun"
)

// MemberTopupOnlineModel: mirror tabel member_topup_online (migration 093 di sudocore2).
// Percobaan top-up saldo member dari channel ONLINE (pos/kiosk/mobile) -- real-time, tanpa
// approval. Baris ini cuma "pemicu" -- begitu status jadi 'paid', baris ke
// member_balance_ledger baru diinsert (saldo member ke-update).
type MemberTopupOnlineModel struct {
	bun.BaseModel `bun:"table:member_topup_online"`

	ID                 int64      `bun:"id,pk,autoincrement"`
	MemberID           int64      `bun:"member_id,notnull"`
	BranchID           *int64     `bun:"branch_id"`
	TerminalID         *int64     `bun:"terminal_id"`
	ReferenceNumber    string     `bun:"reference_number,notnull"`
	Amount             string     `bun:"amount,notnull"`       // NUMERIC(20,2) -- string biar presisi gak keganggu float
	Source             string     `bun:"source,notnull"`       // 'pos', 'kiosk', 'mobile'
	PaymentGatewayCode *string    `bun:"payment_gateway_code"` // null = tunai
	Status             string     `bun:"status,notnull,default:'pending'"`
	ExpiredAt          *time.Time `bun:"expired_at"`
	CancelAt           *time.Time `bun:"cancel_at"`
	PaidAt             *time.Time `bun:"paid_at"`
	Notes              *string    `bun:"notes"`
	CompanyID          *int       `bun:"company_id"`
	CreatedBy          *int64     `bun:"created_by"`
	CreatedAt          time.Time  `bun:"created_at,notnull,default:current_timestamp"`
}

// MemberBalanceLedgerModel: mirror tabel member_balance_ledger (migration 092 di sudocore2).
type MemberBalanceLedgerModel struct {
	bun.BaseModel `bun:"table:member_balance_ledger"`

	ID              int64      `bun:"id,pk,autoincrement"`
	MemberID        int64      `bun:"member_id,notnull"`
	BranchID        *int64     `bun:"branch_id"`
	TerminalID      *int64     `bun:"terminal_id"`
	TransactionDate time.Time  `bun:"transaction_date,notnull,default:current_timestamp"`
	TransactionType string     `bun:"transaction_type,notnull"` // 'topup', 'payment', 'refund', 'adjustment'
	Source          string     `bun:"source,notnull"`           // 'erp', 'pos', 'kiosk'
	ReferenceNumber string     `bun:"reference_number,notnull"`
	BalanceIn       string     `bun:"balance_in,notnull,default:0"`
	BalanceOut      string     `bun:"balance_out,notnull,default:0"`
	BalanceAfter    string     `bun:"balance_after,notnull"`
	Notes           *string    `bun:"notes"`
	CompanyID       *int       `bun:"company_id"`
	IsDeleted       bool       `bun:"is_deleted,notnull,default:false"`
	CreatedBy       *int64     `bun:"created_by"`
	CreatedAt       time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	JurnalAt        *time.Time `bun:"jurnal_at"`
}
