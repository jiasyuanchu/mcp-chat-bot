package models

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request payload for chat API
type ChatRequest struct {
	Message string    `json:"message"`
	History []Message `json:"history"`
}

// ChatResponse represents the response payload for chat API
type ChatResponse struct {
	Response string    `json:"response"`
	History  []Message `json:"history"`
}

// MCPRequest represents the request payload for MCP API
type MCPRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// MCPResponse represents the response payload from MCP API
type MCPResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
