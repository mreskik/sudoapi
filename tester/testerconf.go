package tester

import (
	"APIANDORDER/backend/config"
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

type Tete struct {
	DB *bun.DB
}

func New() *Tete {

	_, filename, _, _ := runtime.Caller(0)
	base := filepath.Dir(filename)
	env := filepath.Join(base, "../.env")

	godotenv.Load(env)
	config.InitDB()
	return &Tete{DB: config.DB}
}

func (data *Tete) Close() {

	err := data.DB.Close()
	if err != nil {
		log.Print("error close database")
	}
	log.Println("Database Closed !")
}
