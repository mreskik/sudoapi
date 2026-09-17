package mobilenotify

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// RunListener: goroutine background, jalan SEKALI dari main.go (`go mobilenotify.RunListener(hub)`),
// hidup selama proses APIANDORDER hidup. Pakai koneksi RAW pgx TERPISAH dari config.DB (yang
// pool, buat query CRUD biasa) -- LISTEN butuh 1 koneksi yang dipegang TERUS, gak boleh
// dibalikin ke pool/dipinjem query lain (beda dari pola query biasa di codebase ini).
//
// Auto-reconnect: kalau koneksi putus (network/DB restart), loop luar bikin ulang koneksi +
// LISTEN lagi dari nol, jeda 5 detik biar gak spam retry.
func RunListener(hub *Hub) {
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	for {
		if err := listenOnce(hub, dsn); err != nil {
			log.Println("mobilenotify: listener error, reconnect 5 detik lagi:", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func listenOnce(hub *Hub, dsn string) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, "LISTEN mb_order_paid"); err != nil {
		return err
	}
	log.Println("mobilenotify: LISTEN mb_order_paid aktif")

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			return err
		}

		branchID, orderNumber, ok := parsePayload(notification.Payload)
		if !ok {
			log.Println("mobilenotify: payload gak valid, dilewatin:", notification.Payload)
			continue
		}

		log.Println("mobilenotify: notifikasi masuk -- branch", branchID, "order", orderNumber)
		hub.Broadcast(branchID, orderNumber)
	}
}

// parsePayload: "{branch_id}:{order_number}" -- format yang disepakati sama sudomobile
// (finalizeSettledPayment()). order_number sendiri gak pernah ngandung ":" (format
// "NO"+branch_code+timestamp+random, lihat generators.go sudomobile), jadi SplitN 2 aman.
func parsePayload(payload string) (branchID int, orderNumber string, ok bool) {
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		return 0, "", false
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", false
	}
	return id, parts[1], true
}
