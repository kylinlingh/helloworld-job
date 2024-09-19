package job

import (
	"context"
	"github.com/ecordell/optgen/helpers"
	"helloworld/internal/dataflow/pump"
	"helloworld/internal/job/try"
	log "helloworld/pkg/logger"
	"helloworld/pkg/signal"
)

type TryJobParamConfig struct {
	*ParamConfig
}

type completedTryJobConfig struct {
	jobs []try.TryJob
	// 程序终止时的回调函数
	closeFunc func() error
}

func (c *completedTryJobConfig) Run(ctx context.Context) error {
	for _, job := range c.jobs {
		job.RunJob(ctx)
	}
	return nil
}

func (c *TryJobParamConfig) CompleteTryJob(ctx context.Context) (RunnableJob, error) {
	log.Ctx(ctx).Info().Fields(helpers.Flatten(c.DebugMap())).Msg("* Configuration: ")
	closeables := closeableStack{}
	var err error
	defer func() {
		// if an error happens during the execution of CompleteJobs, all resources are cleaned up
		if closeableErr := closeables.CloseIfError(err); closeableErr != nil {
			log.Ctx(ctx).Err(closeableErr).Msg("failed to clean up resources on ParamConfig.CompleteJobs")
		}
	}()

	dbInstance, err := c.getDBInstance()
	scanJobs := try.CreateTryJobs(ctx, dbInstance)

	// 开启数据上报功能
	if c.Upload.Enable {

		if c.Upload.Storage == "memory" {
			// 上报数据
			ups, cache := c.getUploadServiceWithMemoryStorage()
			ups.Start()

			pc := c.NewPumps()
			pumpService := pump.CreatePumpService(c.Download.PurgeDelay, pc, cache)

			// 拉取数据并导出到其他目的地
			stopCh := signal.SetupSignalHandler()
			preparedPumpService := pumpService.PrepareRun()
			go func() {
				preparedPumpService.Run(stopCh)
			}()
		} else if c.Upload.Storage == "redis" {

		}
	}

	return &completedTryJobConfig{
		jobs:      scanJobs,
		closeFunc: closeables.Close,
	}, nil
}

func (c *TryJobParamConfig) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ready to run try job")

	log.Info().Msg("all job finished without errors")
	return nil
}

func NewTryJobConfig(baseConf *ParamConfig) *TryJobParamConfig {
	return &TryJobParamConfig{
		ParamConfig: baseConf,
	}
}
