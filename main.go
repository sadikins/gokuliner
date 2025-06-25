package main

import (
	"backend_kalkuliner/config"
	"backend_kalkuliner/database"
	"backend_kalkuliner/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	database.ConnectDatabase(cfg)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	{
		// Routes Bahan Baku
		api.GET("/bahan-bakus", handlers.GetBahanBakus)
		api.POST("/bahan-bakus", handlers.CreateBahanBaku)
		api.GET("/bahan-bakus/:id", handlers.GetBahanBakuByID)
		api.PUT("/bahan-bakus/:id", handlers.UpdateBahanBaku)
		api.DELETE("/bahan-bakus/:id", handlers.DeleteBahanBaku)


		// Routes Resep (Updated)
		api.GET("/reseps", handlers.GetReseps)
		api.POST("/reseps", handlers.CreateResep)
		api.GET("/reseps/:id", handlers.GetResepByID)   // <-- Tambahkan ini
		api.PUT("/reseps/:id", handlers.UpdateResep)    // <-- Tambahkan ini
		api.DELETE("/reseps/:id", handlers.DeleteResep) // <-- Tambahkan ini

		// Routes Perhitungan HPP
		api.GET("/hpp/:resep_id", handlers.GetHPPForResep)
	}

	router.Run(":" + cfg.AppPort)
}