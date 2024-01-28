package logger

import (
	"context"
	"log/slog"
	"maps"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[slog.Level]zapcore.Level{
	slog.LevelDebug: zapcore.DebugLevel,
	slog.LevelInfo:  zapcore.InfoLevel,
	slog.LevelWarn:  zapcore.WarnLevel,
	slog.LevelError: zapcore.ErrorLevel,
}

type SlogContextKey struct {
}

type SlogContextValue map[string]string

func WithContext(ctx context.Context, key string, value string) context.Context {
	newValue := make(SlogContextValue)

	if ctxValue := ctx.Value(SlogContextKey{}); ctxValue != nil {
		switch ctxValue := ctxValue.(type) {
		case SlogContextValue:
			maps.Copy(newValue, ctxValue)
			newValue[key] = value
			return context.WithValue(ctx, SlogContextKey{}, newValue)
		}
	}

	newValue[key] = value
	return context.WithValue(ctx, SlogContextKey{}, newValue)
}

func WithContexts(ctx context.Context, kvMap map[string]string) context.Context {
	newValue := make(SlogContextValue)

	if ctxValue := ctx.Value(SlogContextKey{}); ctxValue != nil {
		switch ctxValue := ctxValue.(type) {
		case SlogContextValue:
			maps.Copy(newValue, ctxValue)
			maps.Copy(newValue, kvMap)
			return context.WithValue(ctx, SlogContextKey{}, newValue)
		}
	}

	maps.Copy(newValue, kvMap)
	return context.WithValue(ctx, SlogContextKey{}, newValue)
}

type ZapHandler struct {
	config zap.Config
	logger *zap.Logger

	attrs []slog.Attr
	group string
}

func (h *ZapHandler) Enabled(_ context.Context, level slog.Level) bool {
	return levelMap[level] >= h.config.Level.Level()
}

func (h *ZapHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := []zapcore.Field{}

	for _, attr := range h.attrs {
		fields = append(fields, zap.Any(attr.Key, attr.Value))
	}

	if ctxValue := ctx.Value(SlogContextKey{}); ctxValue != nil {
		switch ctxValue := ctxValue.(type) {
		case SlogContextValue:
			for k, v := range ctxValue {
				fields = append(fields, zap.Any(k, v))
			}
		}
	}

	record.Attrs(func(a slog.Attr) bool {
		fields = append(fields, zap.Any(a.Key, a.Value))
		return true
	})

	if ce := h.logger.Check(levelMap[record.Level], record.Message); ce != nil {
		ce.Write(fields...)
	}

	return nil
}

func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	groupPrepend := ""
	if len(h.group) > 0 {
		groupPrepend = h.group + "."
	}

	newattrs := append([]slog.Attr{}, h.attrs...)
	for _, attr := range attrs {
		newattrs = append(newattrs, slog.Attr{
			Key:   groupPrepend + attr.Key,
			Value: attr.Value,
		})
	}

	return &ZapHandler{
		config: h.config,
		logger: h.logger,

		attrs: newattrs,
		group: h.group,
	}
}

func (h *ZapHandler) WithGroup(name string) slog.Handler {
	return &ZapHandler{
		config: h.config,
		logger: h.logger,

		attrs: h.attrs,
		group: name,
	}
}
