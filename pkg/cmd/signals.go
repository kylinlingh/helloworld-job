package cmd

import (
	"context"
	log "helloworld/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SignalContextWithGracePeriod creates a new context that will be cancelled
// when an interrupt/SIGTERM signal is received and the provided grace period
// subsequently finishes.
func SignalContextWithGracePeriod(ctx context.Context, gracePeriod time.Duration) context.Context {
	childCtx, cancelFn := context.WithCancel(ctx)
	go func() {
		signalCtx, _ := signal.NotifyContext(childCtx, os.Interrupt, syscall.SIGTERM)
		// 阻塞程序，等待终止信号到达
		<-signalCtx.Done()
		log.Ctx(ctx).Info().Msg("received interrupt signal")

		// 接收到终止信号后，等待宽限时间结束
		if gracePeriod > 0 {
			interruptGrace, _ := signal.NotifyContext(context.Background(), os.Interrupt)
			graceTimer := time.NewTimer(gracePeriod)
			log.Ctx(ctx).Info().Stringer("timeout", gracePeriod).Msg("starting shutdown grace period")

			select {
			case <-graceTimer.C:
			case <-interruptGrace.Done():
				log.Ctx(ctx).Warn().Msg("interrupted shutdown grace period")
			}
		}
		log.Ctx(ctx).Info().Msg("shutting down")
		cancelFn()
	}()
	return childCtx
}
