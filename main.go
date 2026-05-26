package main

import (
	"APIANDORDER/backend"
	"APIANDORDER/backend/config"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env tidak ditemukan")
	}
	config.InitDB()
	defer config.CloseDB()

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	backend.Router(app)
	app.Run("0.0.0.0:7001")

}
