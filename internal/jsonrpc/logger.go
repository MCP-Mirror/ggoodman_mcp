package jsonrpc

import (
	"fmt"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
)

var _ jsonrpc2.Logger = &slogLogger{}

func NewJSONRPCLogger(logger *slog.Logger, extraFields ...any) jsonrpc2.Logger {

	return &slogLogger{
		extraFields: extraFields,
		logger:      logger,
	}
}

type slogLogger struct {
	extraFields []any
	logger      *slog.Logger
}

func (s *slogLogger) Printf(format string, v ...interface{}) {
	s.logger.Debug(fmt.Sprintf(format, v...), s.extraFields...)
}
