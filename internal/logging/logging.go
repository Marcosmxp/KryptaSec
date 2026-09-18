package logging

import (
	"io"
	"log/slog"
	"strings"
)

const redactedValue = "[REDACTED]"

var sensitiveKeys = map[string]struct{}{
	"api_key":       {},
	"apikey":        {},
	"authorization": {},
	"cookie":        {},
	"password":      {},
	"secret":        {},
	"token":         {},
}

func New(w io.Writer, level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: redactAttr,
	})
	return slog.New(handler)
}

func redactAttr(_ []string, attr slog.Attr) slog.Attr {
	if _, sensitive := sensitiveKeys[strings.ToLower(attr.Key)]; sensitive {
		return slog.String(attr.Key, redactedValue)
	}
	return attr
}
