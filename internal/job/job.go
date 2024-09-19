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
	"helloworld/pkg/db"
	log "helloworld/pkg/logger"
	"io"
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

	// From env
	TaskID string `debugmap:"visible"`
	LogDir string `debugmap:"visible"`
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

func NewRunConfig() *ParamConfig {
	fileConfig, err := config.NewConfigFromFile()
	if err != nil {
		log.Fatal().Msgf("failed to initialize because cannot read config from file: %s", err)
	}
	envConfig, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatal().Msgf("failed to initialize because cannot read config from env: %s", err)
	}
	return &ParamConfig{
		App:       &fileConfig.App,
		Log:       &fileConfig.Log,
		Feature:   &fileConfig.Feature,
		Datastore: &fileConfig.DataStore,
		Mysql:     &fileConfig.Mysql,
		Postgres:  &fileConfig.Postgres,
		Upload:    &fileConfig.Upload,
		Download:  &fileConfig.Download,
		Backends:  &fileConfig.Backends,
		TaskID:    envConfig.TaskID,
		LogDir:    envConfig.LogDir,
	}
}
