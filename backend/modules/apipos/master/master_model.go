package master

import (
	"time"

	"github.com/uptrace/bun"
)

// BranchPosSessionModel: session token dinamis multi-device per branch POS (sudocore2 migration
// 284) -- pengganti pola lama (master_branch.token, 1 token statis seumur hidup tanpa scope
// device). 1 branch boleh punya BANYAK baris aktif sekaligus (sengaja TIDAK unique per
// branch_id -- multi-device disepakati eksplisit). Token digenerate & di-insert dari sini
// (APIANDORDER) tiap kali POS install (GetDataBranch), BUKAN dari sudocore2.
type BranchPosSessionModel struct {
	bun.BaseModel `bun:"table:master_branch_pos_session"`

	Id             int64      `bun:"id,pk,autoincrement"`
	BranchId       int64      `bun:"branch_id"`
	Token          string     `bun:"token"`
	DeviceHostname *string    `bun:"device_hostname"`
	FlagActive     bool       `bun:"flag_active"`
	CreatedAt      time.Time  `bun:"created_at,default:current_timestamp"`
	CreatedBy      *int64     `bun:"created_by"`
	UpdatedAt      *time.Time `bun:"updated_at"`
	UpdatedBy      *int64     `bun:"updated_by"`
}
