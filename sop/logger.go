package sop

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string) (*zap.Logger, *zap.SugaredLogger) {
	conf := zap.NewProductionConfig()
	if level != "" {
		if i, e := strconv.Atoi(level); e == nil {
			conf.Level.SetLevel(zapcore.Level(i))
		} else {
			_ = conf.Level.UnmarshalText([]byte(level))
		}
	}
	conf.EncoderConfig.MessageKey = "msg"
	conf.EncoderConfig.LevelKey = "level"
	conf.EncoderConfig.TimeKey = "@t"
	conf.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger, logger.Sugar()
}

type LogHelper struct {
	ctx      context.Context
	funcname string
	logtail  string
	start    time.Time
	last     time.Time
	logtime  string
	logger   *zap.Logger
	err      error
}

func NewLogHelperWithContext(ctx context.Context, funcname string, logtail string) *LogHelper {
	logger := ctxzap.Extract(ctx)
	l := NewLogHelperWithLogger(logger, funcname, logtail)
	l.ctx = ctx
	return l
}

func NewLogHelperWithLogger(logger *zap.Logger, funcname string, logtail string) *LogHelper {
	newLogger := logger.WithOptions(zap.AddCallerSkip(1))
	l := &LogHelper{logtail: logtail, funcname: funcname, start: time.Now(),
		last: time.Now(), logtime: "", logger: newLogger, err: nil}
	sugar := l.logger.Sugar()
	sugar.Infof("BEGIN %s! %s", l.funcname, l.logtail)
	return l
}

func (l *LogHelper) AddTimeInterval(note string) *LogHelper {
	elapsed := int(time.Since(l.last) / time.Millisecond)
	l.logtime += fmt.Sprintf(" %s: %vms", note, elapsed)
	l.last = time.Now()
	return l
}

func (l *LogHelper) AddLogParam(k interface{}, v interface{}) *LogHelper {
	l.logtail += fmt.Sprintf(" %v: %v", k, v)
	return l
}

func (l *LogHelper) AppendLogtail(logtail string) *LogHelper {
	l.logtail += logtail
	return l
}

func (l *LogHelper) SetErr(err error) *LogHelper {
	l.err = err
	return l
}

func (l *LogHelper) LogErr(err error, logstr string) *LogHelper {
	l.SetErr(err).AppendLogtail(logstr)
	return l
}

func (l *LogHelper) End() {
	elapsed := int(time.Since(l.start) / time.Millisecond)
	s := fmt.Sprintf("END %s! %s all cost: %vms ", l.funcname, l.logtime, elapsed)
	if l.err != nil {
		s += fmt.Sprintf("err: <%v> ", l.err)
	}
	s += l.logtail
	l.logger.Info(s)
}
