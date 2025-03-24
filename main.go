package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Message structure representing chat messages
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Chat request structure
type ChatRequest struct {
	Message string    `json:"message"`
	History []Message `json:"history"`
}

// Chat response structure
type ChatResponse struct {
	Response string    `json:"response"`
	History  []Message `json:"history"`
}

// MCP API request structure
type MCPRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// MCP API response structure
type MCPResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found") // ⚡ Changed from `log.Fatal` to `log.Println` as this is not a critical error
	}

	mcpAPIKey := os.Getenv("MCP_API_KEY")
	if mcpAPIKey == "" {
		log.Fatal("MCP_API_KEY environment variable is required")
	}

	router := gin.Default()
	router.Static("/", "./public")

	// Define the chat API endpoint
	router.POST("/api/chat", handleChat(mcpAPIKey))

	// Allow port configuration via environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port) // ⚡ Added logging for better debugging
	router.Run(":" + port)
}

// handleChat handles incoming chat requests
func handleChat(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var chatReq ChatRequest
		if err := c.ShouldBindJSON(&chatReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure history is initialized
		context := chatReq.History
		if context == nil {
			context = []Message{}
		}

		// Append the user's message to history
		context = appendMessage(context, "user", chatReq.Message)

		mcpReq := MCPRequest{
			Model:     "gpt-4",
			Messages:  context,
			MaxTokens: 1000,
		}

		// Marshal the request payload
		jsonData, err := json.Marshal(mcpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		// Make the API request
		resp, err := callMCPAPI(apiKey, jsonData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close() // ⚡ Ensure response body is closed

		// Read the API response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
			return
		}

		// Handle non-200 responses from MCP API
		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("MCP API error: %s", string(body))})
			return
		}

		// Parse MCP API response
		var mcpResp MCPResponse
		if err := json.Unmarshal(body, &mcpResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
			return
		}

		// Ensure MCP API returned a valid response
		if len(mcpResp.Choices) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from MCP API"})
			return
		}

		aiResponse := mcpResp.Choices[0].Message.Content

		// Append AI's response to history
		context = appendMessage(context, "assistant", aiResponse)

		c.JSON(http.StatusOK, ChatResponse{
			Response: aiResponse,
			History:  context,
		})
	}
}

// callMCPAPI sends a request to the MCP API
func callMCPAPI(apiKey string, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.mcp.ai/chat/completions", io.NopCloser(bytes.NewReader(jsonData))) // ⚡ Using `io.NopCloser` for better compatibility
	if err != nil {
		return nil, fmt.Errorf("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// ⚡ Using `http.DefaultClient` instead of creating a new instance
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP API")
	}

	return resp, nil
}

// appendMessage appends a new message to the conversation history
func appendMessage(history []Message, role, content string) []Message {
	return append(history, Message{
		Role:    role,
		Content: content,
	})
}
