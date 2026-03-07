package logger

import (
	"context"

	"github.com/parxyws/nego-gin/pkg/contextkey"
	"github.com/sirupsen/logrus"
)

// WithCtx provides a pre-configured logger entry with context values injected
func WithCtx(ctx context.Context, layer, function string) *logrus.Entry {
	entry := Log.WithFields(logrus.Fields{
		"layer":    layer,
		"function": function,
	})
	if reqID, ok := ctx.Value(contextkey.RequestID).(string); ok && reqID != "" {
		entry = entry.WithField("request_id", reqID)
	}
	return entry
}
