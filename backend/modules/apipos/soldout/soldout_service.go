package soldout

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	DB *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{DB: db}
}

// PushSoldOut: REPLACE PENUH baris master_item_sold_out milik branch_id itu -- DELETE semua
// baris lama punya branch ini, lalu INSERT ulang item_ids yang dikirim (snapshot kondisi sold
// out POS SEKARANG, bukan diff/incremental). ItemIDs kosong = valid, hasilnya tinggal DELETE
// tanpa INSERT (branch itu lagi gak ada item sold out sama sekali).
func (s *Service) PushSoldOut(ctx context.Context, branchID int64, itemIDs []int64) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	if _, err := tx.NewDelete().Model((*SoldOutModel)(nil)).Where("branch_id = ?", branchID).Exec(ctx); err != nil {
		return err
	}

	if len(itemIDs) > 0 {
		now := time.Now()
		rows := make([]SoldOutModel, 0, len(itemIDs))
		for _, itemID := range itemIDs {
			rows = append(rows, SoldOutModel{BranchID: branchID, ItemID: itemID, CreatedAt: now})
		}
		if _, err := tx.NewInsert().Model(&rows).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}
