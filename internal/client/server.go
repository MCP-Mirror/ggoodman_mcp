package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mcp/internal/mcp"
	"os"
	"os/exec"

	"github.com/sourcegraph/jsonrpc2"
)

type MCPServerDefinition struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
	Env         []string
}

type Client struct {
	Id     string
	Server *MCPServerDefinition

	cmd    *exec.Cmd
	conn   *jsonrpc2.Conn
	cancel context.CancelFunc
}

func (c *Client) Start(ctx context.Context, id string, logger *slog.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	args := []string{"run", "--rm", "-i", "-u", "node"}

	for _, envDec := range c.Server.Env {
		args = append(args, "-e", envDec)
	}

	args = append(args, "node:20", c.Server.Cmd)
	args = append(args, c.Server.Args...)

	c.cmd = exec.CommandContext(ctx, "docker", args...)
	c.cmd.Env = append(os.Environ(), c.Server.Env...)

	w, err := c.cmd.StdinPipe()
	if err != nil {
		logger.Error("failed to get stdin pipe", "err", err)
		cancel()
		return err
	}

	r, err := c.cmd.StdoutPipe()
	if err != nil {
		logger.Error("failed to get stdout pipe", "err", err)
		cancel()
		return err
	}

	handler := jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(c.handleRequest).SuppressErrClosed())

	c.conn = jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(&stdioCloser{
		r: r,
		w: w,
	}), handler, jsonrpc2.LogMessages(&slogLogger{logger: logger, id: id}))

	c.cmd.Stderr = os.Stderr

	c.cmd.Start()

	var initializeResult mcp.InitializeResult
	if err := c.conn.Call(ctx, "initialize", &mcp.InitializeRequest{}, &initializeResult); err != nil {
		logger.Error("failed to initialize client", "err", err)
		cancel()
		return err
	}

	logger.Info("client initialized", "id", id, "result", initializeResult)

	if err := c.conn.Notify(ctx, "initialized", &mcp.InitializeNotification{}); err != nil {
		logger.Error("failed to notify client initialization", "err", err)
		cancel()
		return err
	}

	var toolsListResult mcp.ToolsListResult
	if err := c.conn.Call(ctx, "tools/list", &mcp.ToolsListRequest{}, &toolsListResult); err != nil {
		logger.Error("failed to list tools", "err", err)
		cancel()
		return err
	}

	logger.Info("client tools listed", "id", id, "tools", toolsListResult.Tools)

	return nil
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) handleRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	return nil, &jsonrpc2.Error{
		Code:    jsonrpc2.CodeMethodNotFound,
		Message: fmt.Sprintf("method %q not found", req.Method),
	}
}

type slogLogger struct {
	logger *slog.Logger
	id     string
}

func (s *slogLogger) Printf(format string, v ...interface{}) {
	s.logger.Debug(fmt.Sprintf(format, v...), "client", s.id)
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
