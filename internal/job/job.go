package job

import (
	"context"
	"github.com/ecordell/optgen/helpers"
	"github.com/hashicorp/go-multierror"
	"helloworld/config"
	"helloworld/internal/datastore"
	"helloworld/internal/datastore/mysql"
	scan2 "helloworld/internal/job/scan"
	"helloworld/pkg/db"
	log "helloworld/pkg/logger"
	"io"
	"sync"
)

// ParamConfig 用于描述程序运行时的所有配置项（配置文件 + 命令行传入参数 + 其他自定义参数）
//
//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . ParamConfig
type ParamConfig struct {
	// From config.yml
	App       config.App       `debugmap:"visible"`
	Log       config.Log       `debugmap:"visible"`
	Feature   config.Feature   `debugmap:"visible"`
	Datastore config.DataStore `debugmap:"visible"`
	Mysql     config.Mysql     `debugmap:"visible"`
	Postgres  config.Postgres  `debugmap:"visible"`

	// From command flags
	ConfigFile string `debugmap:"visible"`
	ReportPath string `debugmap:"visible"`
}

func (c *ParamConfig) Complete(ctx context.Context) (RunnableJob, error) {
	log.Ctx(ctx).Info().Fields(helpers.Flatten(c.DebugMap())).Msg("configuration")

	closeables := closeableStack{}
	var err error
	defer func() {
		// if an error happens during the execution of Complete, all resources are cleaned up
		if closeableErr := closeables.CloseIfError(err); closeableErr != nil {
			log.Ctx(ctx).Err(closeableErr).Msg("failed to clean up resources on ParamConfig.Complete")
		}
	}()

	dbInstance, err := c.getDBInstance()
	// 通过 Newxxxx 函数实现依赖注入，将数据库实例注入到对象中
	iosJob := scan2.NewIOSScanJob(dbInstance)
	androidJob := scan2.NewAndroidScanJob(dbInstance)

	return &completedJobConfig{
		IOSJob:     iosJob,
		AndroidJob: androidJob,
		closeFunc:  closeables.Close,
	}, nil
}

func (c *ParamConfig) getDBInstance() (datastore.DBFactory, error) {
	var factory datastore.DBFactory
	var err error
	switch c.Datastore.Engine {
	case "mysql":
		mysqlOptions := db.MysqlOptions{
			RunMode:               c.App.RunMode,
			Host:                  c.Mysql.Host,
			Username:              c.Mysql.Username,
			Password:              c.Mysql.Password,
			Database:              c.Mysql.Database,
			MaxIdleConnections:    c.Mysql.MaxIdleConnections,
			MaxOpenConnections:    c.Mysql.MaxOpenConnections,
			MaxConnectionLifeTime: c.Mysql.MaxConnectionLifeTime,
			//LogLevel:              0,
		}
		factory, err = mysql.GetMysqlFactoryOr(&mysqlOptions)
		if err != nil {
			return nil, err
		}
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
	AndroidJob scan2.ScanJob
	IOSJob     scan2.ScanJob

	// 程序终止时的回调函数
	closeFunc func() error
}

func (c *completedJobConfig) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ready to run scan job 1")
	wg := sync.WaitGroup{}
	finishChan := make(chan struct{})
	var multiErr error
	// 运行两个 2 任务
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := c.IOSJob.RunJob(ctx)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}()

	go func() {
		defer wg.Done()
		err := c.AndroidJob.RunJob(ctx)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}()

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

	return nil
}

func NewRunConfig(config *config.Config) *ParamConfig {
	return &ParamConfig{
		App:       config.App,
		Log:       config.Log,
		Feature:   config.Feature,
		Datastore: config.DataStore,
		Mysql:     config.Mysql,
		Postgres:  config.Postgres,
	}
}
