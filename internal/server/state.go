package server

import "mcp/internal/mcp"

type ClientState struct {
	Servers map[string]ServerState
}

type ServerState struct {
	Capabilities mcp.ServerCapabilities
	Tools        []mcp.ToolDefinition
}
