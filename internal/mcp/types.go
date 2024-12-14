package mcp

import (
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

const MCP_PROTOCOL_VERSION = "2024-11-05"

type ImplementationInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ListChangesCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type SubscribeAndListChangesCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type LoggingCapability struct{}

type SamplingCapability struct{}

type ClientCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Roots        *ListChangesCapability `json:"roots,omitempty"`
	Sampling     *SamplingCapability    `json:"sampling,omitempty"`
}

type InitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ImplementationInfo `json:"clientInfo"`
}

type ServerCapabilities struct {
	Experimental map[string]interface{}             `json:"experimental,omitempty"`
	Logging      *LoggingCapability                 `json:"logging,omitempty"`
	Prompts      *ListChangesCapability             `json:"prompts,omitempty"`
	Resources    *SubscribeAndListChangesCapability `json:"resources,omitempty"`
	Tools        *ListChangesCapability             `json:"tools,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ImplementationInfo `json:"serverInfo"`
	Instructions    *string            `json:"instructions,omitempty"`
}

type InitializedNotification struct{}

type Meta map[string]any

type ToolsCallRequest struct {
	ToolName  string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolsCallResult struct {
	Meta    Meta  `json:"_meta,omitempty"`
	Content []any `json:"content"`
	IsError bool  `json:"isError"`
}

type ToolDefinition struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	InputSchema JSONSchemaObject `json:"inputSchema"`
}

type JSONSchemaObject struct {
	Type               string         `json:"type"`
	Properties         map[string]any `json:"properties"`
	RequiredProperties []string       `json:"required"`
}

type JSONSchemaUnknown any

type ToolsListRequest struct{}

type ToolsListResult struct {
	Tools []ToolDefinition `json:"tools"`
}

func MustParams[T any](req *jsonrpc2.Request) (*T, error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: "missing params",
		}
	}

	var r T
	if err := json.Unmarshal(*req.Params, &r); err != nil {
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		}
	}
	return &r, nil
}
