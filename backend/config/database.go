package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func InitDB() {
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println("Error open koneksi ! : ", err)
	}

	DB = bun.NewDB(sqlDB, pgdialect.New())

	if os.Getenv("APP_ENV") == "development" {
		DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	log.Println("database connected!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()

		log.Println("DB Closed!")
	}
}
