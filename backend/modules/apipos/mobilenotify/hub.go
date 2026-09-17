// Package mobilenotify: relay "ada order mobile baru" dari sudomobile ke worker POS
// (worker_mobile_customer) per branch, lewat WebSocket. sudomobile cuma nembak pg_notify()
// (fire-and-forget, gak nunggu siapa pun) -- APIANDORDER yang LISTEN (listener.go) & nerusin ke
// worker POS yang lagi konek buat branch itu (hub.go, handler.go). Disepakati 2026-08-26, lihat
// CATATAN INTERNAL.md sudomobile.
package mobilenotify

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Hub: registry koneksi WS AKTIF, dikelompokkan per branch_id. 1 hub dipakai bareng oleh
// listener.go (nulis/broadcast) dan handler.go (daftar/bubar) -- makanya butuh mutex, 2-2nya
// jalan di goroutine beda.
type Hub struct {
	mu    sync.Mutex
	conns map[int][]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{conns: make(map[int][]*websocket.Conn)}
}

func (h *Hub) Register(branchID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[branchID] = append(h.conns[branchID], conn)
}

func (h *Hub) Unregister(branchID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	list := h.conns[branchID]
	for i, c := range list {
		if c == conn {
			h.conns[branchID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

// Broadcast: kirim message ke SEMUA koneksi yang lagi terdaftar buat branchID ini. Copy slice-nya
// dulu SEBELUM nulis ke koneksi (di luar lock) -- WriteMessage bisa lambat/blocking (network),
// gak boleh nahan lock hub selama itu (nge-block Register/Unregister dari goroutine lain).
func (h *Hub) Broadcast(branchID int, message string) {
	h.mu.Lock()
	conns := append([]*websocket.Conn{}, h.conns[branchID]...)
	h.mu.Unlock()

	for _, c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
