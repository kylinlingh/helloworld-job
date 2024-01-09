package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/hashicorp/go-multierror"
	"helloworld/config"
	"helloworld/internal/dataflow/pump"
	"helloworld/internal/dataflow/storage/uploadto/memory"
	"helloworld/internal/dataflow/upload"
	"helloworld/internal/datastore"
	"helloworld/internal/datastore/mysql"
	"helloworld/internal/datastore/postgres"
	"helloworld/internal/job/scan"
	"helloworld/pkg/db"
	log "helloworld/pkg/logger"
	"helloworld/pkg/signal"
	"io"
	"sync"
)

// ParamConfig 用于描述程序运行时的所有配置项（配置文件 + 命令行传入参数 + 其他自定义参数）
//
//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . ParamConfig
type ParamConfig struct {
	// From config.yml
	App       *config.App       `debugmap:"visible"`
	Log       *config.Log       `debugmap:"visible"`
	Feature   *config.Feature   `debugmap:"visible"`
	Datastore *config.DataStore `debugmap:"visible"`
	Mysql     *config.Mysql     `debugmap:"visible"`
	Postgres  *config.Postgres  `debugmap:"visible"`
	Upload    *config.Upload    `debugmap:"visible"`
	Download  *config.Download  `debugmap:"visible"`
	Backends  *config.Backends  `debugmap:"visible"`

	// From command flags
	//ConfigFile string `debugmap:"visible"`
	//ReportPath string `debugmap:"visible"`

	ScanTargets []string `debugmap:"visible"`
	ScanMode    string   `debugmap:"visible"`
	CodePath    string   `debugmap:"visible"`
}

var jobMap map[string]func(factory datastore.DBFactory) scan.ScanJob

const (
	ANDROID_SOURCECODE = "android_sc"
	ANDROID_ARTIFACT   = "android_af"
	IOS_SCOURCECODE    = "ios_sc"
	IOS_ARTIFACT       = "ios_af"
)

func initJobMap() {
	jobMap = make(map[string]func(factory datastore.DBFactory) scan.ScanJob)
	jobMap[ANDROID_SOURCECODE] = scan.NewAndroidSourceCodeScanJob
	jobMap[ANDROID_ARTIFACT] = scan.NewAndroidArtifactScanJob
	jobMap[IOS_SCOURCECODE] = scan.NewIOSSourceCodeScanJob
	jobMap[IOS_ARTIFACT] = scan.NewIOSArtifactScanJob

}

func (c *ParamConfig) Complete(ctx context.Context) (RunnableJob, error) {
	//log.Ctx(ctx).Info().Fields(helpers.Flatten(c.DebugMap())).Msg("configuration as: ")

	initJobMap()

	closeables := closeableStack{}
	var err error
	defer func() {
		// if an error happens during the execution of Complete, all resources are cleaned up
		if closeableErr := closeables.CloseIfError(err); closeableErr != nil {
			log.Ctx(ctx).Err(closeableErr).Msg("failed to clean up resources on ParamConfig.Complete")
		}
	}()

	dbInstance, err := c.getDBInstance()

	// 创建多个扫描任务
	scanJobs := []scan.ScanJob{}
	for _, st := range c.ScanTargets {
		if initFunc, ok := jobMap[st]; !ok {
			log.Ctx(ctx).Fatal().Err(errors.New(fmt.Sprintf("invalid scan target: %s [%s, %s, %s, %s]", st, ANDROID_SOURCECODE, ANDROID_ARTIFACT, IOS_ARTIFACT, IOS_SCOURCECODE)))
		} else {
			scanJobs = append(scanJobs, initFunc(dbInstance))
		}
	}

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

	sp := &scan.JobParam{CodePath: c.CodePath}

	return &completedJobConfig{
		jobs:       scanJobs,
		ScanParams: sp,
		closeFunc:  closeables.Close,
	}, nil
}

func (c *ParamConfig) getUploadServiceWithMemoryStorage() (*upload.UploadService, *ristretto.Cache) {
	storage := &memory.UploadMemoryStorage{}
	storage.Connect()

	uploadConf := &upload.UploadConfig{
		Enable:                  c.Upload.Enable,
		WorkersNum:              c.Upload.WorkersNum,
		FlushInterval:           c.Upload.FlushInterval,
		RecordsBufferSize:       c.Upload.RecordsBufferSize,
		EnableDetailedRecording: c.Upload.EnableDetailedRecording,
	}

	return upload.CreateUploadService(uploadConf, storage), storage.GetStorage()
}

func (c *ParamConfig) NewPumps() map[string]pump.PumpConfig {
	m := make(map[string]pump.PumpConfig)
	for _, name := range c.Download.Backends {
		switch name {
		case "csv":
			m["csv"] = pump.PumpConfig{
				Type: "csv",
				Meta: map[string]interface{}{
					"csv_dir": c.Backends.CSV.CSVDIR,
				},
			}
		}
	}
	return m
}

func (c *ParamConfig) getDBInstance() (datastore.DBFactory, error) {
	var factory datastore.DBFactory
	var err error
	switch c.Datastore.Engine {
	case "mysql":
		mysqlOptions := db.MysqlOptions{
			RunMode:               c.App.RunMode,
			Host:                  c.Mysql.Host,
			Port:                  c.Mysql.Port,
			Username:              c.Mysql.Username,
			Password:              c.Mysql.Password,
			Database:              c.Mysql.DBName,
			MaxIdleConnections:    c.Mysql.MaxIdleConnections,
			MaxOpenConnections:    c.Mysql.MaxOpenConnections,
			MaxConnectionLifeTime: c.Mysql.MaxConnectionLifeTime,
			//LogLevel:              0,
		}
		factory, err = mysql.GetMysqlFactoryOr(&mysqlOptions)
		if err != nil {
			return nil, err
		}
	case "psql":
		psqlOpts := db.PsqlOptions{
			RunMode:               c.App.RunMode,
			Host:                  c.Postgres.PHost,
			Port:                  c.Postgres.PPort,
			Username:              c.Postgres.DBUser,
			Password:              c.Postgres.DBPassword,
			DBName:                c.Postgres.DBNAME,
			MaxIdleConnections:    c.Postgres.PMaxIdleConnections,
			MaxOpenConnections:    c.Postgres.PMaxOpenConnections,
			MaxConnectionLifeTime: c.Postgres.PMaxConnectionLifeTime,
		}
		factory, err = postgres.GetPsqlFactoryOr(&psqlOpts)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("not invalid datastore engine: %s", c.Datastore.Engine))
	}
	return factory, nil
}

type closeableStack struct {
	closers []func() error
}

func (c *closeableStack) AddWithError(closer func() error) {
	c.closers = append(c.closers, closer)
}

// AddCloser try to call Close() for closer
func (c *closeableStack) AddCloser(closer io.Closer) {
	if closer != nil {
		c.closers = append(c.closers, closer.Close)
	}
}

func (c *closeableStack) AddWithoutError(closer func()) {
	c.closers = append(c.closers, func() error {
		closer()
		return nil
	})
}

func (c *closeableStack) Close() error {
	var err error
	// closer in reverse order how it's expected in deferred funcs
	for i := len(c.closers) - 1; i >= 0; i-- {
		if closerErr := c.closers[i](); closerErr != nil {
			err = multierror.Append(err, closerErr)
		}
	}
	return err
}

func (c *closeableStack) CloseIfError(err error) error {
	if err != nil {
		return c.Close()
	}
	return nil
}

// RunnableJob is a job ready to run
type RunnableJob interface {
	Run(ctx context.Context) error
}

type completedJobConfig struct {
	jobs       []scan.ScanJob
	ScanParams *scan.JobParam

	// 程序终止时的回调函数
	closeFunc func() error
}

func (c *completedJobConfig) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ready to run scan job")
	wg := sync.WaitGroup{}
	finishChan := make(chan struct{})
	var multiErr error
	jobCount := len(c.jobs)

	// 创建新的 context 用于传递参数
	newContext := context.WithValue(ctx, scan.KEY_JOB_PARAM, c.ScanParams)

	wg.Add(jobCount)
	for _, job := range c.jobs {
		go func(scanJob scan.ScanJob) {
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

	log.Info().Msg("all job finished without errors")
	return nil
}

func NewRunConfig() *ParamConfig {
	readConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal().Msgf("Initialization failed: %s", err)
	}
	return &ParamConfig{
		App:       &readConfig.App,
		Log:       &readConfig.Log,
		Feature:   &readConfig.Feature,
		Datastore: &readConfig.DataStore,
		Mysql:     &readConfig.Mysql,
		Postgres:  &readConfig.Postgres,
		Upload:    &readConfig.Upload,
		Download:  &readConfig.Download,
		Backends:  &readConfig.Backends,
	}
}
