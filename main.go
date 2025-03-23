package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message string    `json:"message"`
	History []Message `json:"history"`
}

type ChatResponse struct {
	Response string    `json:"response"`
	History  []Message `json:"history"`
}

type MCPRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type MCPResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	mcpAPIKey := os.Getenv("MCP_API_KEY")
	if mcpAPIKey == "" {
		log.Fatal("MCP_API_KEY environment variable is required")
	}

	router := gin.Default()

	router.Static("/", "./public")

	router.POST("/api/chat", func(c *gin.Context) {
		var chatReq ChatRequest
		if err := c.ShouldBindJSON(&chatReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		context := chatReq.History
		if context == nil {
			context = []Message{}
		}

		context = append(context, Message{
			Role:    "user",
			Content: chatReq.Message,
		})

		mcpReq := MCPRequest{
			Model:     "gpt-4",
			Messages:  context,
			MaxTokens: 1000,
		}

		jsonData, err := json.Marshal(mcpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		client := &http.Client{}
		httpReq, err := http.NewRequest("POST", "https://api.mcp.ai/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+mcpAPIKey)

		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call MCP API"})
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("MCP API error: %s", string(body))})
			return
		}

		var mcpResp MCPResponse
		if err := json.Unmarshal(body, &mcpResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
			return
		}

		if len(mcpResp.Choices) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from MCP API"})
			return
		}

		aiResponse := mcpResp.Choices[0].Message.Content

		context = append(context, Message{
			Role:    "assistant",
			Content: aiResponse,
		})

		c.JSON(http.StatusOK, ChatResponse{
			Response: aiResponse,
			History:  context,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
