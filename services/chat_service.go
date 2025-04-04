package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mcp-chat-bot/models"
)

const (
	mcpAPIEndpoint = "https://api.mcp.ai/chat/completions"
	contentType    = "application/json"
)

// CallMCPAPI sends a request to the MCP API
func CallMCPAPI(apiKey string, mcpReq models.MCPRequest) (*models.MCPResponse, error) {
	jsonData, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := createRequest(mcpAPIEndpoint, apiKey, jsonData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP API: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp)
}

// createRequest creates a new HTTP request with appropriate headers
func createRequest(url, apiKey string, jsonData []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewReader(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req, nil
}

// handleResponse processes the API response
func handleResponse(resp *http.Response) (*models.MCPResponse, error) {
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
