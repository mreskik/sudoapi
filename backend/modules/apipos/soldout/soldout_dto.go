package soldout

// PushSoldOutRequestDTO: payload tombol push sold out di POS -- SNAPSHOT PENUH kondisi sold
// out branch itu SEKARANG (bukan diff/incremental). ItemIDs kosong/[] berarti "gak ada item
// yang sold out saat ini", valid -- bukan error, hasilnya semua baris branch itu di ERP
// ke-hapus (lihat soldout_service.go, PushSoldOut: DELETE dulu baru INSERT ulang).
type PushSoldOutRequestDTO struct {
	BranchID int64   `json:"branch_id"`
	ItemIDs  []int64 `json:"item_ids"`
}
