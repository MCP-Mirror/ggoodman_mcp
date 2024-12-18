package localbroker

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mcp/internal/integrations"
	"mcp/internal/jsonrpc"
	"mcp/internal/mcp"
	serverrunner "mcp/internal/server_runner"
	"mcp/internal/util"
	"strings"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

const (
	DEFAULT_START_TIMEOUT_SECONDS = 30
)

var ErrConnectionClosed = fmt.Errorf("connection closed")

type LocalBroker interface {
	Close() error
	Run(ctx context.Context) error
}

var _ LocalBroker = &localBroker{}

type localBroker struct {
	integRepo   integrations.IntegrationsRepository
	integRunner serverrunner.ServerStarter
	logger      *slog.Logger
	conn        *jsonrpc2.Conn

	integrationStartTimeout time.Duration
}

func NewLocalBroker(
	ctx context.Context,
	logger *slog.Logger,
	integRepo integrations.IntegrationsRepository,
	runner serverrunner.ServerStarter,
	r io.ReadCloser,
	w io.WriteCloser,
) LocalBroker {
	lb := &localBroker{
		integRepo:   integRepo,
		integRunner: runner,
		logger:      logger,

		integrationStartTimeout: time.Duration(DEFAULT_START_TIMEOUT_SECONDS) * time.Second,
	}

	handler := jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(lb.handleRequest).SuppressErrClosed())
	stream := jsonrpc2.NewPlainObjectStream(util.NewReaderWriterCloser(r, w))

	lb.conn = jsonrpc2.NewConn(ctx, stream, handler, jsonrpc2.LogMessages(jsonrpc.NewJSONRPCLogger(lb.logger)))

	return lb
}

func (lb *localBroker) Close() error {
	lb.conn.Close()
	return nil
}

func (lb *localBroker) Run(ctx context.Context) error {
	defer lb.Close()

	lb.integRepo.OnIntegrationsChanged(func(e *integrations.IntegrationsChangedEvent) {
		switch e.Type {
		case integrations.IntegrationsChangedEventTypeAdded:
			go lb.startIntegration(ctx, e.Integration)
		case integrations.IntegrationsChangedEventTypeRemoved:
			go lb.stopIntegration(ctx, e.Integration)
		}
	})

	installed, err := lb.integRepo.ListIntegrations(ctx)
	if err != nil {
		return fmt.Errorf("error listing integrations: %w", err)
	}

	lb.logger.Debug("bootstrapping integrations", "count", len(installed))

	for _, integration := range installed {
		lb.logger.Info("bootstrapping integration", "id", integration.Id)
		go lb.startIntegration(ctx, *integration)
	}

	select {
	case <-ctx.Done():
	case <-lb.conn.DisconnectNotify():
		return ErrConnectionClosed
	}

	return nil
}

func (lb *localBroker) startIntegration(ctx context.Context, integration integrations.InstalledIntegration) {
	ctx, cancel := context.WithTimeout(ctx, lb.integrationStartTimeout)
	defer cancel()

	lb.logger.Info("starting integration", "id", integration.Id)

	srv, err := lb.integRunner.Create(ctx, serverrunner.ServerDescription{
		Runtime: integration.Manifest.Runtime,
		Command: integration.Manifest.Command,
		Args:    integration.Manifest.Args,
		Env:     integration.Env,
	})
	if err != nil {
		lb.logger.Error("error creating integration", "id", integration.Id, "err", err)
		return
	}

	if err := srv.Run(ctx); err != nil {
		lb.logger.Error("error running integration", "id", integration.Id, "err", err)
	}

	lb.removeIntegrationById(ctx, integration.Id)

}

func (lb *localBroker) stopIntegration(_ context.Context, integration integrations.InstalledIntegration) {
	lb.logger.Info("stopping integration", "id", integration.Id)
}

func (lb *localBroker) removeIntegrationById(_ context.Context, integrationId string) {
	lb.logger.Info("removing integration", "id", integrationId)
}

func (lb *localBroker) handleRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	lb.logger.Debug("handling request", "method", req.Method)
	switch req.Method {
	case "initialize":
		req, err := mcp.MustParams[mcp.InitializeRequest](req)
		if err != nil {
			return nil, err
		}
		return lb.handleInitializeRequest(ctx, conn, req)
	case "initialized":
		req, err := mcp.MustParams[mcp.InitializedNotification](req)
		if err != nil {
			return nil, err
		}
		return nil, lb.handleInitializedNotification(ctx, conn, req)
	case "tools/call":
		req, err := mcp.MustParams[mcp.ToolsCallRequest](req)
		if err != nil {
			return nil, err
		}
		return lb.handleToolsCallRequest(ctx, conn, req)
	case "tools/list":
		req, err := mcp.MustParams[mcp.ToolsListRequest](req)
		if err != nil {
			return nil, err
		}
		return lb.handleToolsListRequest(ctx, conn, req)
	default:
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method %q not found", req.Method),
		}
	}
}

func (lb *localBroker) handleInitializeRequest(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
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

func (lb *localBroker) handleInitializedNotification(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.InitializedNotification) error {
	return nil
}

func (lb *localBroker) handleToolsCallRequest(_ context.Context, _ *jsonrpc2.Conn, req *mcp.ToolsCallRequest) (*mcp.ToolsCallResult, error) {
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

func (lb *localBroker) handleToolsListRequest(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.ToolsListRequest) (*mcp.ToolsListResult, error) {
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
