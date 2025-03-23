# MCP Chat Bot (Golang Version)

This is a chat bot project implemented using Golang and the Gin framework, integrating various AI models through the Model Context Protocol (MCP).

## Features

- Simple web chat interface
- Communication with various AI models using the MCP protocol
- Conversation history tracking
- Efficient server implementation using Golang and Gin framework

## Installation Steps

1. Clone this repository
   ```bash
   git clone https://github.com/yourusername/mcp-chat-bot.git
   cd mcp-chat-bot
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Create a `.env` file and set your MCP API key
   ```
   MCP_API_KEY=your_api_key_here
   PORT=8080
   ```

4. Compile and start the server
   ```bash
   go build -o mcp-chat-bot
   ./mcp-chat-bot
   ```

5. Open `http://localhost:8080` in your browser to start chatting

## Project Structure

- `main.go` - Gin server and MCP integration
- `public/index.html` - Chat interface
- `go.mod` - Go module dependencies

## Customization

You can modify the `Model` parameter in the `main.go` file to use different AI models:

```go
mcpReq := MCPRequest{
    Model:     "gpt-4", // Change to your preferred model
    Messages:  context,
    MaxTokens: 1000,
}
```
