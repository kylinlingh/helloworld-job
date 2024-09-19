package sop

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	AddFields    = ctxzap.AddFields
	TagsToFields = ctxzap.AddFields
	ToContext    = ctxzap.ToContext
	Debug        = ctxzap.Debug
)

func Extract(ctx context.Context) *SpanLogger {
	return &SpanLogger{
		Logger: ctxzap.Extract(ctx).WithOptions(zap.AddCallerSkip(1)),
		span:   opentracing.SpanFromContext(ctx),
	}
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Error(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Warn(msg, fields...)
}

type SpanLogger struct {
	*zap.Logger
	span           opentracing.Span
	spanWithFields []zap.Field
}

func (log *SpanLogger) Sugar() *SpanSugaredLogger {
	return &SpanSugaredLogger{
		SugaredLogger:  log.Logger.Sugar(),
		span:           log.span,
		spanWithFields: log.spanWithFields,
	}
}

func (log *SpanLogger) Named(s string) *SpanLogger {
	return &SpanLogger{
		Logger:         log.Logger.Named(s),
		span:           log.span,
		spanWithFields: log.spanWithFields,
	}
}

func (log *SpanLogger) WithOptions(opts ...zap.Option) *SpanLogger {
	return &SpanLogger{
		Logger:         log.Logger.WithOptions(opts...),
		span:           log.span,
		spanWithFields: log.spanWithFields,
	}
}

func (log *SpanLogger) With(fields ...zap.Field) *SpanLogger {
	spanWithFields := log.spanWithFields
	spanWithFields = append(spanWithFields, fields...)
	return &SpanLogger{
		Logger:         log.Logger.With(fields...),
		span:           log.span,
		spanWithFields: spanWithFields,
	}
}

func (log *SpanLogger) Info(msg string, fields ...zap.Field) {
	log.Logger.Info(msg, fields...)
	logSpan(log.span, msg, log.spanWithFields, fields)
}

func (log *SpanLogger) Warn(msg string, fields ...zap.Field) {
	log.Logger.Warn(msg, fields...)
	logSpan(log.span, msg, log.spanWithFields, fields)
}

func (log *SpanLogger) Error(msg string, fields ...zap.Field) {
	log.Logger.Error(msg, fields...)
	logSpan(log.span, msg, log.spanWithFields, fields)
}

type SpanSugaredLogger struct {
	*zap.SugaredLogger
	span           opentracing.Span
	spanWithFields []zap.Field
}

func (s *SpanSugaredLogger) Desugar() *SpanLogger {
	return &SpanLogger{
		Logger:         s.SugaredLogger.Desugar(),
		span:           s.span,
		spanWithFields: s.spanWithFields,
	}
}

func (s *SpanSugaredLogger) Named(name string) *SpanSugaredLogger {
	return &SpanSugaredLogger{
		SugaredLogger:  s.SugaredLogger.Named(name),
		span:           s.span,
		spanWithFields: s.spanWithFields,
	}
}

func (s *SpanSugaredLogger) With(args ...interface{}) *SpanSugaredLogger {
	spanWithFields := s.spanWithFields
	spanWithFields = append(spanWithFields, sweetenFields(args)...)
	return &SpanSugaredLogger{
		SugaredLogger:  s.SugaredLogger.With(args...),
		span:           s.span,
		spanWithFields: spanWithFields,
	}
}

func (s *SpanSugaredLogger) Info(args ...interface{}) {
	s.SugaredLogger.Info(args...)
	logSpanSugar(s.span, s.spanWithFields, "", args, nil)
}

func (s *SpanSugaredLogger) Warn(args ...interface{}) {
	s.SugaredLogger.Warn(args...)
	logSpanSugar(s.span, s.spanWithFields, "", args, nil)
}

func (s *SpanSugaredLogger) Error(args ...interface{}) {
	s.SugaredLogger.Error(args...)
	logSpanSugar(s.span, s.spanWithFields, "", args, nil)
}

func (s *SpanSugaredLogger) Infof(template string, args ...interface{}) {
	s.SugaredLogger.Infof(template, args...)
	logSpanSugar(s.span, s.spanWithFields, template, args, nil)
}

func (s *SpanSugaredLogger) Warnf(template string, args ...interface{}) {
	s.SugaredLogger.Warnf(template, args...)
	logSpanSugar(s.span, s.spanWithFields, template, args, nil)
}

func (s *SpanSugaredLogger) Errorf(template string, args ...interface{}) {
	s.SugaredLogger.Errorf(template, args...)
	logSpanSugar(s.span, s.spanWithFields, template, args, nil)
}

func (s *SpanSugaredLogger) Infow(msg string, keysAndValues ...interface{}) {
	s.SugaredLogger.Infow(msg, keysAndValues...)
	logSpanSugar(s.span, s.spanWithFields, msg, nil, keysAndValues)
}

func (s *SpanSugaredLogger) Warnw(msg string, keysAndValues ...interface{}) {
	s.SugaredLogger.Warnw(msg, keysAndValues...)
	logSpanSugar(s.span, s.spanWithFields, msg, nil, keysAndValues)
}

func (s *SpanSugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {
	s.SugaredLogger.Errorw(msg, keysAndValues...)
	logSpanSugar(s.span, s.spanWithFields, msg, nil, keysAndValues)
}

func logSpan(span opentracing.Span, msg string, spanWithFields []zap.Field, fields []zap.Field) {
	if span == nil {
		return
	}
	spanFields := make([]opentracinglog.Field, 0, 2+len(fields))
	spanFields = append(spanFields, opentracinglog.Event("log"), opentracinglog.Message(msg))
	for _, f := range spanWithFields {
		if f.Type == zapcore.SkipType {
			continue
		}
		spanFields = append(spanFields, zapFieldToOpentracing(f))
	}
	for _, f := range fields {
		if f.Type == zapcore.SkipType {
			continue
		}
		spanFields = append(spanFields, zapFieldToOpentracing(f))
	}
	span.LogFields(spanFields...)
}

func zapFieldToOpentracing(zapField zapcore.Field) opentracinglog.Field {
	switch zapField.Type {
	case zapcore.BoolType:
		val := false
		if zapField.Integer >= 1 {
			val = true
		}
		return opentracinglog.Bool(zapField.Key, val)
	case zapcore.Float32Type:
		return opentracinglog.Float32(zapField.Key, math.Float32frombits(uint32(zapField.Integer)))
	case zapcore.Float64Type:
		return opentracinglog.Float64(zapField.Key, math.Float64frombits(uint64(zapField.Integer)))
	case zapcore.Int64Type:
		return opentracinglog.Int64(zapField.Key, int64(zapField.Integer))
	case zapcore.Int32Type:
		return opentracinglog.Int32(zapField.Key, int32(zapField.Integer))
	case zapcore.StringType:
		return opentracinglog.String(zapField.Key, zapField.String)
	case zapcore.StringerType:
		return opentracinglog.String(zapField.Key, zapField.Interface.(fmt.Stringer).String())
	case zapcore.Uint64Type:
		return opentracinglog.Uint64(zapField.Key, uint64(zapField.Integer))
	case zapcore.Uint32Type:
		return opentracinglog.Uint32(zapField.Key, uint32(zapField.Integer))
	case zapcore.DurationType:
		return opentracinglog.String(zapField.Key, time.Duration(zapField.Integer).String())
	case zapcore.ErrorType:
		return opentracinglog.Error(zapField.Interface.(error))
	default:
		return opentracinglog.Object(zapField.Key, zapField.Interface)
	}
}

func logSpanSugar(span opentracing.Span, spanWithFields []zap.Field, template string, fmtArgs []interface{}, keysAndValues []interface{}) {
	if span == nil {
		return
	}
	kv := make([]interface{}, 0, 4+len(spanWithFields)*2+len(keysAndValues))
	kv = append(kv, "event", "log", "message", getMessage(template, fmtArgs))

	for _, f := range spanWithFields {
		if f.Type == zapcore.SkipType {
			continue
		}
		of := zapFieldToOpentracing(f)
		kv = append(kv, of.Key(), of.Value())
	}
	kv = append(kv, keysAndValues...)
	span.LogKV(kv...)
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}
	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}
	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}
	fields := make([]zap.Field, 0, len(args))
	for i := 0; i < len(args); {
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}
		if i == len(args)-1 {
			break
		}
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); ok {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}
	return fields
}
