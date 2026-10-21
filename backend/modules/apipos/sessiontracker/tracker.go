// Package sessiontracker -- state in-memory per session token POS (master_branch_pos_session),
// dipakai middleware.BranchTokenAuth tiap kali token lolos validasi. 2 hal disimpan beda cara:
//
//   - LastUsedAt -- SELALU in-memory, TIDAK PERNAH ditulis ke DB. Sengaja "hilang" (balik nil)
//     tiap APIANDORDER restart -- trade-off disetujui biar gak nambah 1 UPDATE query per
//     request cuma buat timestamp yang nilainya nyaris gak penting presisinya. Diekspos ke
//     sudocore2 lewat endpoint internal (lihat internal_handler.go).
//   - AppVersion/OSVersion/IPAddress -- metadata yang PERLU persisten (device beneran ganti
//     versi/lokasi, bukan angka yang terus jalan kayak timestamp), jadi tetap ditulis ke DB
//     (master_branch_pos_session), tapi di-THROTTLE lewat cache di sini -- UPDATE DB cuma
//     kejadian kalau nilainya beda dari yang udah diketahui, bukan tiap request.
package sessiontracker

import (
	"sync"
	"time"
)

type Metadata struct {
	AppVersion string
	OSVersion  string
	IPAddress  string
}

type tracker struct {
	mu        sync.RWMutex
	lastUsed  map[int64]time.Time
	lastMeta  map[int64]Metadata
}

var global = &tracker{
	lastUsed: make(map[int64]time.Time),
	lastMeta: make(map[int64]Metadata),
}

// TouchLastUsed -- catat "dipakai sekarang" buat session id ini. Dipanggil tiap request yang
// lolos BranchTokenAuth, murah (in-memory doang).
func TouchLastUsed(sessionId int64) {
	global.mu.Lock()
	global.lastUsed[sessionId] = time.Now()
	global.mu.Unlock()
}

// GetLastUsed -- buat dibaca endpoint internal (sudocore2 nanya lewat sini). nil kalau belum
// pernah dipakai sejak APIANDORDER ini nyala (termasuk sesaat setelah restart).
func GetLastUsed(sessionId int64) *time.Time {
	global.mu.RLock()
	defer global.mu.RUnlock()
	t, ok := global.lastUsed[sessionId]
	if !ok {
		return nil
	}
	return &t
}

// GetLastUsedBatch -- versi batch buat endpoint internal (List monitoring bisa punya banyak
// baris sekaligus, jangan sampai sudocore2 kudu panggil endpoint ini 1 per 1 row).
func GetLastUsedBatch(sessionIds []int64) map[int64]*time.Time {
	global.mu.RLock()
	defer global.mu.RUnlock()
	result := make(map[int64]*time.Time, len(sessionIds))
	for _, id := range sessionIds {
		if t, ok := global.lastUsed[id]; ok {
			tCopy := t
			result[id] = &tCopy
		} else {
			result[id] = nil
		}
	}
	return result
}

// ShouldUpdateMetadata -- bandingkan metadata baru vs cache terakhir yang diketahui. Balikin
// true (dan update cache-nya) kalau ADA yang beda -- caller baru UPDATE DB kalau true, jadi
// request yang metadata-nya gak berubah sama sekali gak nyentuh DB.
func ShouldUpdateMetadata(sessionId int64, incoming Metadata) bool {
	global.mu.Lock()
	defer global.mu.Unlock()

	cached, ok := global.lastMeta[sessionId]
	if ok && cached == incoming {
		return false
	}
	global.lastMeta[sessionId] = incoming
	return true
}
