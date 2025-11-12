package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"school-assistant-wh/internal/config"
	"school-assistant-wh/internal/handlers"
	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/services/facebook"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fbCfg := config.LoadFacebookConfig()

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.Message{})

	fbSvc := facebook.NewService(fbCfg)

	h := handlers.NewHandler(db, fbSvc)

	if err := h.SetupMessengerProfile(); err != nil {
		log.Printf("Warning: Failed to set up Messenger profile: %v", err)
	} else {
		log.Println("Successfully set up Messenger profile")
	}

	r := setupRouter(h)

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	go func() {
		log.Printf("Server starting on port %s...\n", port)
		if err := r.Run(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

func initDB() (*gorm.DB, error) {
	cfg := config.LoadDBConfig()
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func setupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	webhook := r.Group("/webhook")
	{
		webhook.GET("", h.VerifyWebhook)
		webhook.POST("", h.HandleWebhook)
	}

	return r
}
