package sop

import (
	"sync"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type jaegerLogger struct {
	Logger *zap.Logger

	s *zap.SugaredLogger
	o sync.Once
}

func SetJaegerLogger(logger *zap.Logger) jaegercfg.Option {
	return jaegercfg.Logger(&jaegerLogger{Logger: logger})
}

func (l *jaegerLogger) Error(msg string) {
	l.sugar().Error(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	l.sugar().Infof(msg, args...)
}

func (l *jaegerLogger) Debugf(msg string, args ...interface{}) {
	l.sugar().Debugf(msg, args...)
}

func (l *jaegerLogger) sugar() *zap.SugaredLogger {
	if l != nil {
		l.o.Do(func() {
			if l.Logger != nil {
				l.s = l.Logger.Sugar()
			}
		})
		if l.s != nil {
			return l.s
		}
	}
	return zap.S()
}
