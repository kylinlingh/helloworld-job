package job

import (
	"context"
	"github.com/ecordell/optgen/helpers"
	"golang.org/x/sync/errgroup"
	"helloworld/config"
	"helloworld/internal/scan"
	log "helloworld/pkg/logger"
)

// ParamConfig 用于描述程序运行时的所有配置项（配置文件 + 命令行传入参数 + 其他自定义参数）
//
//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . ParamConfig
type ParamConfig struct {
	// From config.yml
	App        config.App     `debugmap:"visible"`
	Log        config.Log     `debugmap:"visible"`
	PostgreSQL config.PG      `debugmap:"visible"`
	Feature    config.Feature `debugmap:"visible"`

	// From command flags
	ConfigFile string `debugmap:"visible"`
	ReportPath string `debugmap:"visible"`
}

func (c *ParamConfig) Complete(ctx context.Context) (RunnableJob, error) {
	log.Ctx(ctx).Info().Fields(helpers.Flatten(c.DebugMap())).Msg("configuration")

	return &completedJobConfig{}, nil
}

// RunnableJob is a job ready to run
type RunnableJob interface {
	Run(ctx context.Context) error
}

type completedJobConfig struct {
	AndroidJob scan.AndroidScanJob
	IOSJob     scan.IOSScanJob

	// 程序终止时的回调函数
	closeFunc func() error
}

func (c *completedJobConfig) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ready to run scan job")
	g, ctx := errgroup.WithContext(ctx)

	stopOnCancelWithErr := func(stopFn func() error) func() error {
		return func() error {
			// ctx 被关闭时，回调 stopFn
			<-ctx.Done()
			return stopFn()
		}
	}

	g.Go(func() error { return c.IOSJob.RunJob(ctx) })

	g.Go(func() error { return c.AndroidJob.RunJob(ctx) })

	g.Go(stopOnCancelWithErr(c.closeFunc))

	if err := g.Wait(); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("error shutting down job")
		return err
	}

	return nil
}

func NewRunConfig(config *config.Config) *ParamConfig {
	return &ParamConfig{
		App:        config.App,
		Log:        config.Log,
		PostgreSQL: config.PG,
		Feature:    config.Feature,
	}
}
