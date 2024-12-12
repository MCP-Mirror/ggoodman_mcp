package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mcp/internal/client"
	"mcp/internal/mcp"
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

func NewServer(ctx context.Context, logger *slog.Logger, r io.ReadCloser, w io.WriteCloser) (io.Closer, error) {
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

			logger.Info("client started", "id", id.String())
		}()
	}

	return server, nil
}

type server struct {
	logger    *slog.Logger
	conn      *jsonrpc2.Conn
	clients   map[string]*client.Client
	clientsMu sync.Mutex
	state     *ServerState
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
		req, err := mcp.MustInitializeRequest(req)
		if err != nil {
			return nil, err
		}
		return s.handleInitializeRequest(ctx, conn, req)
	case "initialized":
		req, err := mcp.MustInitializedNotification(req)
		if err != nil {
			return nil, err
		}
		return nil, s.handleInitializedNotification(ctx, conn, req)
	default:
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method %q not found", req.Method),
		}
	}
}

func (s *server) handleInitializeRequest(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
	// TODO: Negotiate protocol versions

	return &mcp.InitializeResult{
		ProtocolVersion: mcp.MCP_PROTOCOL_VERSION,
		Capabilities: mcp.ServerCapabilities{
			Logging: &mcp.LoggingCapability{},
		},
		ServerInfo: mcp.ImplementationInfo{
			Name:    "mcp",
			Version: "0.1.0",
		},
	}, nil
}

func (s *server) handleInitializedNotification(_ context.Context, _ *jsonrpc2.Conn, _ *mcp.InitializeNotification) error {
	return nil
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
