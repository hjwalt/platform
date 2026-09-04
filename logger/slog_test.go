package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --- helpers ---------------------------------------------------------------

func contextValueMap(ctx context.Context) (SlogContextValue, bool) {
	value, ok := ctx.Value(SlogContextKey{}).(SlogContextValue)
	return value, ok
}

// newBufferZapHandler builds a ZapHandler over a zap.Logger writing JSON to an
// in-memory buffer, plus its underlying config (used by Enabled).
func newBufferZapHandler(level zapcore.Level) (*ZapHandler, *bytes.Buffer) {
	buffer := &bytes.Buffer{}
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(buffer), level)
	zapLogger := zap.New(core)

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)

	return &ZapHandler{config: config, logger: zapLogger}, buffer
}

func slogRecord(level slog.Level, message string, attrs ...slog.Attr) slog.Record {
	record := slog.NewRecord(time.Now(), level, message, 0)
	record.AddAttrs(attrs...)
	return record
}

// --- WithContext ------------------------------------------------------------

func TestWithContextStoresSingleKeyValue(t *testing.T) {
	assert := assert.New(t)

	ctx := WithContext(context.Background(), "key", "value")

	value, ok := contextValueMap(ctx)
	assert.True(ok)
	assert.Equal(SlogContextValue{"key": "value"}, value)
}

func TestWithContextChainingMergesAndLaterValueWins(t *testing.T) {
	assert := assert.New(t)

	first := WithContext(context.Background(), "a", "1")
	second := WithContext(first, "b", "2")
	conflict := WithContext(second, "a", "updated")

	value, ok := contextValueMap(second)
	assert.True(ok)
	assert.Equal(SlogContextValue{"a": "1", "b": "2"}, value)

	value, ok = contextValueMap(conflict)
	assert.True(ok)
	assert.Equal(SlogContextValue{"a": "updated", "b": "2"}, value, "later value should win on key conflict")

	// The original context must not have been mutated by chaining.
	value, ok = contextValueMap(first)
	assert.True(ok)
	assert.Equal(SlogContextValue{"a": "1"}, value)
}

func TestWithContextIgnoresNonSlogContextValue(t *testing.T) {
	assert := assert.New(t)

	// A value that is not a SlogContextValue (a plain string) sits under the key.
	foreign := context.WithValue(context.Background(), SlogContextKey{}, "foreign value")
	ctx := WithContext(foreign, "k", "v")

	// The foreign value must not be copied into the new map; it is replaced.
	value, ok := contextValueMap(ctx)
	assert.True(ok)
	assert.Equal(SlogContextValue{"k": "v"}, value)

	// Same when the foreign value is a plain (untyped) map[string]string.
	foreignMap := context.WithValue(context.Background(), SlogContextKey{}, map[string]string{"old": "1"})
	ctxMap := WithContext(foreignMap, "k", "v")

	value, ok = contextValueMap(ctxMap)
	assert.True(ok)
	assert.Equal(SlogContextValue{"k": "v"}, value)
}

// --- WithContexts -----------------------------------------------------------

func TestWithContextsStoresMap(t *testing.T) {
	assert := assert.New(t)

	ctx := WithContexts(context.Background(), map[string]string{"x": "1", "y": "2"})

	value, ok := contextValueMap(ctx)
	assert.True(ok)
	assert.Equal(SlogContextValue{"x": "1", "y": "2"}, value)
}

func TestWithContextsMergesWithPriorValues(t *testing.T) {
	assert := assert.New(t)

	prior := WithContext(context.Background(), "a", "1")
	merged := WithContexts(prior, map[string]string{"b": "2"})

	value, ok := contextValueMap(merged)
	assert.True(ok)
	assert.Equal(SlogContextValue{"a": "1", "b": "2"}, value)

	// keys from the new map win on conflict
	conflict := WithContexts(prior, map[string]string{"a": "updated"})
	value, ok = contextValueMap(conflict)
	assert.True(ok)
	assert.Equal(SlogContextValue{"a": "updated"}, value)
}

func TestWithContextsEmptyMapKeepsPriorValues(t *testing.T) {
	assert := assert.New(t)

	prior := WithContext(context.Background(), "a", "1")
	ctx := WithContexts(prior, map[string]string{})

	value, ok := contextValueMap(ctx)
	assert.True(ok)
	assert.Equal(SlogContextValue{"a": "1"}, value)
}

// --- ZapHandler.Enabled -----------------------------------------------------

func TestZapHandlerEnabledRespectsConfigLevel(t *testing.T) {
	cases := []struct {
		name            string
		configLevel     zapcore.Level
		debug, info     bool
		warn, errorSeen bool
	}{
		{name: "warn config", configLevel: zapcore.WarnLevel, debug: false, info: false, warn: true, errorSeen: true},
		{name: "debug config", configLevel: zapcore.DebugLevel, debug: true, info: true, warn: true, errorSeen: true},
		{name: "error config", configLevel: zapcore.ErrorLevel, debug: false, info: false, warn: false, errorSeen: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			config := zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(tc.configLevel)
			handler := &ZapHandler{config: config, logger: zap.NewNop()}

			assert.Equal(tc.debug, handler.Enabled(context.Background(), slog.LevelDebug), "debug")
			assert.Equal(tc.info, handler.Enabled(context.Background(), slog.LevelInfo), "info")
			assert.Equal(tc.warn, handler.Enabled(context.Background(), slog.LevelWarn), "warn")
			assert.Equal(tc.errorSeen, handler.Enabled(context.Background(), slog.LevelError), "error")
		})
	}
}

// --- ZapHandler.Handle ------------------------------------------------------

func TestZapHandlerHandleWritesMessageAttrsAndContextValues(t *testing.T) {
	assert := assert.New(t)

	handler, buffer := newBufferZapHandler(zapcore.DebugLevel)
	ctx := WithContext(context.Background(), "ctxkey", "ctxval")

	err := handler.Handle(ctx, slogRecord(slog.LevelInfo, "hello world", slog.String("recattr", "recval")))

	assert.NoError(err)
	output := buffer.String()
	assert.Contains(output, `"msg":"hello world"`)
	assert.Contains(output, `"level":"info"`)
	assert.Contains(output, `"recattr":"recval"`, "record attrs should appear as fields")
	assert.Contains(output, `"ctxkey":"ctxval"`, "context values should appear as fields")
}

func TestZapHandlerHandleIgnoresForeignContextValue(t *testing.T) {
	assert := assert.New(t)

	handler, buffer := newBufferZapHandler(zapcore.DebugLevel)
	ctx := context.WithValue(context.Background(), SlogContextKey{}, "foreign value")

	err := handler.Handle(ctx, slogRecord(slog.LevelInfo, "msg"))

	assert.NoError(err)
	assert.NotContains(buffer.String(), "foreign value")
}

func TestZapHandlerHandleDropsRecordBelowLoggerLevel(t *testing.T) {
	assert := assert.New(t)

	handler, buffer := newBufferZapHandler(zapcore.ErrorLevel)

	err := handler.Handle(context.Background(), slogRecord(slog.LevelInfo, "dropped message"))

	assert.NoError(err)
	assert.Equal("", buffer.String(), "nothing below the logger level should be written")
}

// --- ZapHandler.WithAttrs ---------------------------------------------------

func TestZapHandlerWithAttrsAppendsAndDoesNotMutate(t *testing.T) {
	assert := assert.New(t)

	handler, buffer := newBufferZapHandler(zapcore.DebugLevel)

	withOne := handler.WithAttrs([]slog.Attr{slog.String("one", "1")})
	withTwo := withOne.WithAttrs([]slog.Attr{slog.String("two", "2")})

	// attrs accumulate; the source handler is not mutated.
	assert.Equal([]slog.Attr{slog.String("one", "1")}, withOne.(*ZapHandler).attrs)
	assert.Equal(
		[]slog.Attr{slog.String("one", "1"), slog.String("two", "2")},
		withTwo.(*ZapHandler).attrs,
	)
	assert.Equal(0, len(handler.attrs))

	err := withTwo.Handle(context.Background(), slogRecord(slog.LevelInfo, "with attrs"))
	assert.NoError(err)
	output := buffer.String()
	assert.Contains(output, `"one":"1"`)
	assert.Contains(output, `"two":"2"`)
}

// --- ZapHandler.WithGroup ---------------------------------------------------

func TestZapHandlerWithGroupPrefixesSubsequentAttrs(t *testing.T) {
	assert := assert.New(t)

	handler, buffer := newBufferZapHandler(zapcore.DebugLevel)

	grouped := handler.WithGroup("grp").WithAttrs([]slog.Attr{slog.String("gk", "gv")})

	assert.Equal("grp.gk", grouped.(*ZapHandler).attrs[0].Key)

	err := grouped.Handle(context.Background(), slogRecord(slog.LevelInfo, "grouped"))
	assert.NoError(err)
	assert.Contains(buffer.String(), `"grp.gk":"gv"`)
}

func TestZapHandlerWithGroupDoesNotPrefixPreviouslyAddedAttrs(t *testing.T) {
	assert := assert.New(t)

	handler, _ := newBufferZapHandler(zapcore.DebugLevel)

	withAttr := handler.WithAttrs([]slog.Attr{slog.String("pre", "1")})
	grouped := withAttr.WithGroup("grp").WithAttrs([]slog.Attr{slog.String("post", "2")})

	// Only attrs added after the group is set get the prefix.
	assert.Equal(
		[]slog.Attr{
			slog.String("pre", "1"),
			{Key: "grp.post", Value: slog.StringValue("2")},
		},
		grouped.(*ZapHandler).attrs,
	)
}

func TestZapHandlerWithGroupOverwritesPreviousGroup(t *testing.T) {
	assert := assert.New(t)

	handler, _ := newBufferZapHandler(zapcore.DebugLevel)

	// Oddity: WithGroup replaces the current group rather than nesting it.
	// Standard slog semantics would produce "first.second.k".
	grouped := handler.WithGroup("first").WithGroup("second").WithAttrs([]slog.Attr{slog.String("k", "v")})

	assert.Equal("second.k", grouped.(*ZapHandler).attrs[0].Key)
}

func TestZapHandlerWithGroupKeepsExistingAttrs(t *testing.T) {
	assert := assert.New(t)

	handler, _ := newBufferZapHandler(zapcore.DebugLevel)

	withAttr := handler.WithAttrs([]slog.Attr{slog.String("k", "v")})
	grouped := withAttr.WithGroup("grp")

	// Attrs recorded before WithGroup are carried over untouched, and the
	// handler passed to WithGroup is not mutated.
	assert.Equal(1, len(grouped.(*ZapHandler).attrs))
	assert.Equal("grp", grouped.(*ZapHandler).group)
	assert.Equal("", withAttr.(*ZapHandler).group)
}
