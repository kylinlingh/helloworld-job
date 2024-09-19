package job

import (
	"context"
	"github.com/ecordell/optgen/helpers"
	"github.com/hashicorp/go-multierror"
	"helloworld/internal/dataflow/pump"
	"helloworld/internal/job/scan"
	log "helloworld/pkg/logger"
	"helloworld/pkg/signal"
	"sync"
	"time"
)

// ScanJobParamConfig 用于专门的扫描任务
//
// //go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . ScanJobParamConfig
type ScanJobParamConfig struct {
	*ParamConfig
	ScanTargets []string `debugmap:"visible"`
	ScanMode    string   `debugmap:"visible"`
	CodePath    string   `debugmap:"visible"`
}

func NewScanJobConfig(baseConf *ParamConfig) *ScanJobParamConfig {
	return &ScanJobParamConfig{
		ParamConfig: baseConf,
	}
}

func (c *ScanJobParamConfig) CompleteJobs(ctx context.Context) (RunnableJob, error) {
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
	scanJobs := scan.CreateScanJobs(ctx, dbInstance, c.ScanTargets)

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

	sp := &scan.JobParam{CodePath: c.CodePath, TaskID: c.TaskID}

	return &completedScanJobConfig{
		jobs:       scanJobs,
		ScanParams: sp,
		closeFunc:  closeables.Close,
	}, nil
}

type completedScanJobConfig struct {
	jobs       []scan.ApplicationScanJob
	ScanParams *scan.JobParam

	// 程序终止时的回调函数
	closeFunc func() error
}

func (c *completedScanJobConfig) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ready to run scan job")
	wg := sync.WaitGroup{}
	finishChan := make(chan struct{})
	var multiErr error
	jobCount := len(c.jobs)

	// 创建新的 context 用于传递参数
	newContext := context.WithValue(ctx, scan.KEY_JOB_PARAM, c.ScanParams)

	wg.Add(jobCount)
	for _, job := range c.jobs {
		go func(scanJob scan.ApplicationScanJob) {
			defer wg.Done()
			err := scanJob.RunJob(newContext)
			if err != nil {
				multiErr = multierror.Append(multiErr, err)
			}
		}(job)
	}

	// 在后台等待任务结束
	go func() {
		wg.Wait()
		close(finishChan)
	}()

	select {
	case <-finishChan:
		log.Ctx(ctx).Info().Msg("job finished")
		if multiErr != nil {
			log.Ctx(ctx).Info().Msg("error detected")
			return multiErr
		}
	case <-ctx.Done():
		log.Ctx(ctx).Warn().Msg("interrupt signal caught, closing resources")
		c.closeFunc()
	}

	log.Info().Msg("sleep 10s to wait for data flushing")
	time.Sleep(10 * time.Second)
	log.Info().Msg("all job finished without errors")
	return nil
}
