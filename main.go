package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Message represents a message in our system
type Message struct {
	gorm.Model
	SenderID string `json:"sender_id"`
	Text     string `json:"text"`
}

var db *gorm.DB

func main() {
	// Initialize database connection
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/messenger?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
db.AutoMigrate(&Message{})

	// Initialize Gin router
	r := gin.Default()

	// Webhook verification endpoint (for Facebook Messenger)
	r.GET("/webhook", verifyWebhook)

	// Webhook handler endpoint
	r.POST("/webhook", handleWebhook)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Webhook verification handler
func verifyWebhook(c *gin.Context) {
	verifyToken := c.Query("hub.verify_token")
	if verifyToken == "YOUR_VERIFY_TOKEN" {
		challenge := c.Query("hub.challenge")
		c.String(http.StatusOK, challenge)
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification token"})
}

// Webhook message handler
func handleWebhook(c *gin.Context) {
	var request struct {
		Object string `json:"object"`
		Entry  []struct {
			Messaging []struct {
				Sender struct {
					ID string `json:"id"`
				} `json:"sender"`
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"messaging"`
		} `json:"entry"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process each message
	for _, entry := range request.Entry {
		for _, messaging := range entry.Messaging {
			senderID := messaging.Sender.ID
			messageText := messaging.Message.Text

			// Save message to database
			message := Message{
				SenderID: senderID,
				Text:     messageText,
			}
			db.Create(&message)

			// Echo the message back (simple response)
			if messageText == "hi" || messageText == "Hi" {
				sendTextMessage(senderID, "Hello! How can I help you today?")
			} else {
				sendTextMessage(senderID, "You said: "+messageText)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Helper function to send text message back to user
func sendTextMessage(recipientID, text string) {
	// In a real implementation, you would call the Facebook Messenger API here
	// This is a placeholder for the actual implementation
	log.Printf("Sending message to %s: %s", recipientID, text)
}
