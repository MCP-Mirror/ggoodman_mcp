package localbroker

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mcp/internal/client"
	"mcp/internal/mcp"
	serverrunner "mcp/internal/server_runner"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/sourcegraph/jsonrpc2"
)

var servers []client.MCPServerDefinition = []client.MCPServerDefinition{
	{
		Name:        "@modelcontextprotocol/server-filesystem",
		Description: "MCP server for filesystem access",
		Cmd:         "npx",
		Args: []string{
			"-y",
			"@modelcontextprotocol/server-github",
		},
		Env: []string{
			"GITHUB_PERSONAL_ACCESS_TOKEN=yoink",
		},
	},
}

func NewServer(ctx context.Context, logger *slog.Logger, runner serverrunner.ServerRunner, r io.ReadCloser, w io.WriteCloser) (io.Closer, error) {
	server := &server{
		logger:  logger,
		clients: make(map[string]*client.Client),
	}

	handler := jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(server.handleRequest).SuppressErrClosed())

	conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(&stdioCloser{
		r: r,
		w: w,
	}), handler, jsonrpc2.LogMessages(&slogLogger{logger: logger}))

	server.conn = conn

	for _, serverDef := range servers {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		client := &client.Client{
			Id:     id.String(),
			Server: &serverDef,
		}

		go func() {
			if err := client.Start(ctx, id.String(), logger); err != nil {
				logger.Error("failed to start client", "err", err)

				client.Close()
				return
			}

			// Only add bootstrapped clients
			server.clientsMu.Lock()
			server.clients[id.String()] = client
			server.clientsMu.Unlock()

			logger.Info("child MCP server bootstrapped", "id", id.String(), "name", serverDef.Name)
		}()
	}

	return server, nil
}

type server struct {
	logger    *slog.Logger
	conn      *jsonrpc2.Conn
	clients   map[string]*client.Client
	clientsMu sync.Mutex
}

func (s *server) Close() error {
	s.logger.Debug("closing server")
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	for _, client := range s.clients {
		client.Close()
	}

	return s.conn.Close()
}

func (s *server) handleRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	switch req.Method {
	case "initialize":
		req, err := mcp.MustParams[mcp.InitializeRequest](req)
		if err != nil {
			return nil, err
		}
		return s.handleInitializeRequest(ctx, conn, req)
	case "initialized":
		req, err := mcp.MustParams[mcp.InitializedNotification](req)
		if err != nil {
			return nil, err
		}
		return nil, s.handleInitializedNotification(ctx, conn, req)
	case "tools/call":
		req, err := mcp.MustParams[mcp.ToolsCallRequest](req)
		if err != nil {
			return nil, err
		}
		return s.handleToolsCallRequest(ctx, conn, req)
	case "tools/list":
		req, err := mcp.MustParams[mcp.ToolsListRequest](req)
		if err != nil {
			return nil, err
		}
		return s.handleToolsListRequest(ctx, conn, req)
	default:
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method %q not found", req.Method),
		}
	}
}

func (s *server) handleInitializeRequest(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
	// TODO: Negotiate protocol versions

	instructions := strings.TrimSpace(`
# Introduction

The MCP broker is a Model Context Protocol (MCP) server that manages multiple
child MCP servers. It is responsible for starting, stopping, and managing
the lifecycle of child servers.

The MCP broker will list the Prompts, Resources and Tools exposed by child
servers. These capabilities will be prefixed with an identifier that is unique
to each child server.

In addition to the capabilities offered by child servers, the MCP broker also
exposes sever of its own tools. These tools are prefixed with the identifier
` + "`" + `"__mcp__"` + "`" + ` and should ONLY be used as a last-resort when none of the
installed child servers can fulfill the request.

# Instructions

1. NEVER use the tools prefixed with ` + "`" + `"__mcp__"` + "`" + ` unless no other Prompt,
	 Resource or Tool fulfills the request.
2. ONLY use the tool ` + "`" + `"__mcp__search_registry"` + "`" + ` when no child server
   provides the capabilities necessary to fulfill the request. NEVER call
	 this tool if there is any ambiguity. Indicate that you CAN'T fulfill
	 the request if you are not sure. IMPORTANT: This tool will NEVER actually help you
	 fulfill the request. Use this only if the user asks you to search for tools.
3. ALWAYS use the tool ` + "`" + `"__mcp__suggest_tool"` + "`" + ` describing a child
   MCP server that would be generically useful for the request and other
	 anticipated requests. ALWAYS indicate that you CAN'T fullfill the
	 request after calling this tool. This Tool will NEVER change whether
	 the request can be fulfilled or not.

# Tools

- ` + "`" + `"__mcp__search_registry"` + "`" + `: Search the MCP registry for Servers that
	provide the capabilities necessary to fulfill the request.
- ` + "`" + `"__mcp__install_server"` + "`" + `: Install a server from the MCP registry.
- ` + "`" + `"__mcp__suggest_tool"` + "`" + `: Suggest a child MCP server that would be
	generically useful for the request and other anticipated requests.
			`)

	return &mcp.InitializeResult{
		ProtocolVersion: mcp.MCP_PROTOCOL_VERSION,
		Capabilities: mcp.ServerCapabilities{
			Logging: &mcp.LoggingCapability{},
		},
		ServerInfo: mcp.ImplementationInfo{
			Name:    "mcp",
			Version: "0.1.0",
		},
		Instructions: &instructions,
	}, nil
}

func (s *server) handleInitializedNotification(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.InitializedNotification) error {
	return nil
}

func (s *server) handleToolsCallRequest(_ context.Context, _ *jsonrpc2.Conn, req *mcp.ToolsCallRequest) (*mcp.ToolsCallResult, error) {
	switch req.ToolName {
	case "__mcp__install_server":
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: "tool not implemented",
		}
	case "__mcp__search_registry":
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: "tool not implemented",
		}
	case "__mcp__suggest_tool":
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: "tool not implemented",
		}
	}

	return nil, &jsonrpc2.Error{
		Code:    jsonrpc2.CodeMethodNotFound,
		Message: fmt.Sprintf("tool %q not found", req.ToolName),
	}
}

func (s *server) handleToolsListRequest(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.ToolsListRequest) (*mcp.ToolsListResult, error) {
	builtInTools := []mcp.ToolDefinition{
		{
			Name: "__mcp__install_server",
			Description: strings.TrimSpace(`
Install a child MCP Server based on a server discovered via ` + "`" + `"__mcp__search_registry"` + "`" + `.

ONLY use this tool if you discover a server that can fulfill the request. If you are not sure
whether the server can fulfill the request, indicate that you CAN'T fulfill the request.
			`),
			InputSchema: mcp.JSONSchemaObject{
				Type: "object",
				Properties: map[string]any{
					"name": map[string]any{
						"type": "string",
					},
					"version": map[string]any{
						"type": "string",
					},
				},
				RequiredProperties: []string{"name", "version"},
			},
		},
		{
			Name: "__mcp__search_registry",
			Description: strings.TrimSpace(`
Search the MCP registry for Servers that provide the capabilities necessary to fulfill the request.

When using this tool, DO NOT simply pass on the user's request. Instead, try to come up with
search terms that would be useful for finding a server that can fulfill the request.

If you 
			`),
			InputSchema: mcp.JSONSchemaObject{
				Type: "object",
				Properties: map[string]any{
					"query": map[string]any{
						"type": "string",
					},
				},
			},
		},
		{
			Name:        "__mcp__suggest_tool",
			Description: "Suggest a child MCP server that would be generically useful for the request and other anticipated requests.",
			InputSchema: mcp.JSONSchemaObject{
				Type: "object",
				Properties: map[string]any{
					"tools": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"name": map[string]any{
									"type": "string",
								},
								"description": map[string]any{
									"type": "string",
								},
								"inputSchema": map[string]any{
									"type":       "object",
									"properties": true,
								},
							},
						},
					},
				},
			},
		},
	}

	return &mcp.ToolsListResult{
		Tools: builtInTools,
	}, nil
}

type slogLogger struct {
	logger *slog.Logger
}

func (s *slogLogger) Printf(format string, v ...interface{}) {
	s.logger.Debug(fmt.Sprintf(format, v...), "role", "server")
}

var _ io.ReadWriteCloser = &stdioCloser{}

type stdioCloser struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (s *stdioCloser) Read(p []byte) (n int, err error) {
	return s.r.Read(p)
}

func (s *stdioCloser) Write(p []byte) (n int, err error) {
	return s.w.Write(p)
}

func (s *stdioCloser) Close() error {
	rCloseErr := s.r.Close()
	wCloseErr := s.w.Close()

	if rCloseErr != nil {
		return rCloseErr
	}

	if wCloseErr != nil {
		return wCloseErr
	}

	return nil
}
