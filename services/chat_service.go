package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mcp-chat-bot/models"
)

// CallMCPAPI sends a request to the MCP API
func CallMCPAPI(apiKey string, mcpReq models.MCPRequest) (*models.MCPResponse, error) {
	jsonData, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.mcp.ai/chat/completions", io.NopCloser(bytes.NewReader(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("MCP API error: %s", string(body))
	}

	var mcpResp models.MCPResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(mcpResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from MCP API")
	}

	return &mcpResp, nil
}

// AppendMessage adds a new message to the conversation history
func AppendMessage(history []models.Message, role, content string) []models.Message {
	return append(history, models.Message{
		Role:    role,
		Content: content,
	})
}
