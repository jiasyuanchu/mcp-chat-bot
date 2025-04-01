package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"mcp-chat-bot/models"
	"mcp-chat-bot/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	mcpAPIKey := os.Getenv("MCP_API_KEY")
	if mcpAPIKey == "" {
		log.Fatal("MCP_API_KEY environment variable is required")
	}

	router := gin.Default()
	router.Static("/", "./public")

	router.POST("/api/chat", handleChat(mcpAPIKey))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}

func handleChat(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var chatReq models.ChatRequest
		if err := c.ShouldBindJSON(&chatReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure history is initialized
		context := chatReq.History
		if context == nil {
			context = []models.Message{}
		}

		context = services.AppendMessage(context, "user", chatReq.Message)

		// Prepare MCP request
		mcpReq := models.MCPRequest{
			Model:     "gpt-4",
			Messages:  context,
			MaxTokens: 1000,
		}

		// Call MCP API
		mcpResp, err := services.CallMCPAPI(apiKey, mcpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Append AI's response
		aiResponse := mcpResp.Choices[0].Message.Content
		context = services.AppendMessage(context, "assistant", aiResponse)

		c.JSON(http.StatusOK, models.ChatResponse{
			Response: aiResponse,
			History:  context,
		})
	}
}
