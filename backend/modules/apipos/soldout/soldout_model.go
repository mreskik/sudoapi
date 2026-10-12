package soldout

import (
	"time"

	"github.com/uptrace/bun"
)

// SoldOutModel: baris master_item_sold_out (migration 264 di sudocore2) -- 1 baris = 1 item
// yang LAGI sold out di branch itu. item_id = master_item.id (item FISIK, sama persis
// mr_item.id/flag_soldout di POS), BUKAN item_conversion_detail_id.
type SoldOutModel struct {
	bun.BaseModel `bun:"table:master_item_sold_out"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	BranchID  int64     `bun:"branch_id" json:"branch_id"`
	ItemID    int64     `bun:"item_id" json:"item_id"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`
}
