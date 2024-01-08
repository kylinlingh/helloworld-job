package config

//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . Config App Log Feature DataStore Mysql Postgres Upload Download Backends CSVOpt

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	// Config 用于描述配置文件
	Config struct {
		App       `yaml:"app"`
		Log       `yaml:"logger"`
		Feature   `yaml:"feature"`
		DataStore `yaml:"datastore"`
		Mysql     `yaml:"mysql"`
		Postgres  `yaml:"psql"`
		Upload    `yaml:"upload"`
		Download  `yaml:"download"`
		Backends  `yaml:"backends"`
		CSVOpt    `yaml:"csv"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME" debugmap:"visible"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION" debugmap:"visible"`
		RunMode string `yaml:"run_mode" env:"APP_RUN_MODE" debugmap:"visible"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL" env-default:"debug" debugmap:"visible"`
	}

	Feature struct {
		ShutdownGracePeriod time.Duration `yaml:"shutdown_grace_period"   env:"FEATURE_SHUTDOWN_GRACE_PERIOD" env-default:"0s" debugmap:"visible"`
	}

	DataStore struct {
		Engine             string        `yaml:"engine"   env:"DATASTORE_ENGINE" env-default:"mysql" debugmap:"visible"`
		GcWindows          time.Duration `debugmap:"visible"`
		GcMaxOperationTime time.Duration `debugmap:"visible"`
		MigrationPhase     string        `debugmap:"visible"`
	}

	Mysql struct {
		Host                  string        `yaml:"host"   env:"MYSQL_HOST" env-default:"none" debugmap:"visible"`
		Port                  int           `yaml:"port"   env:"MYSQL_PORT" env-default:"3306" debugmap:"visible"`
		Username              string        `yaml:"username"   env:"MYSQL_USER_NAME" env-default:"none" debugmap:"visible"`
		Password              string        `yaml:"password"   env:"MYSQL_PASSWORD" env-default:"none" debugmap:"sensitive"`
		DBName                string        `yaml:"db_name"   env:"MYSQL_DB_NAME" env-default:"none" debugmap:"visible"`
		MaxIdleConnections    int           `yaml:"max_idle_connections"   env:"MYSQL_MAX_IDLE_CONNECTIONS" env-default:"none" debugmap:"visible"`
		MaxOpenConnections    int           `yaml:"max_open_connections"   env:"MYSQL_MAX_OPEN_CONNECTIONS" env-default:"none" debugmap:"visible"`
		MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"   env:"MYSQL_MAX_CONNECTION_LIFE_TIME" env-default:"none" debugmap:"visible"`
	}

	Postgres struct {
		PHost                  string        `yaml:"host"   env:"POSTGRES_HOST" env-default:"127.0.0.1" debugmap:"visible"`
		PPort                  int           `yaml:"port"   env:"POSTGRES_PORT" env-default:"5432" debugmap:"visible"`
		DBNAME                 string        `yaml:"db_name"   env:"POSTGRES_DB_NAME" env-default:"postgres" debugmap:"visible"`
		DBUser                 string        `yaml:"db_user"   env:"POSTGRES_DBU_SER" env-default:"postgres" debugmap:"visible"`
		DBPassword             string        `yaml:"db_password"   env:"POSTGRES_DB_PASSWORD" env-default:"" debugmap:"visible"`
		PMaxIdleConnections    int           `yaml:"max_idle_connections"   env:"MYSQL_MAX_IDLE_CONNECTIONS" env-default:"none" debugmap:"visible"`
		PMaxOpenConnections    int           `yaml:"max_open_connections"   env:"MYSQL_MAX_OPEN_CONNECTIONS" env-default:"none" debugmap:"visible"`
		PMaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"   env:"MYSQL_MAX_CONNECTION_LIFE_TIME" env-default:"none" debugmap:"visible"`
	}

	Upload struct {
		Enable                  bool          `env-required:"true" yaml:"enable"   env:"REPORT_ENABLE" env-default:"false" debugmap:"visible"`
		Storage                 string        `yaml:"storage"   env:"REPORT_STORAGE" env-default:"memory" debugmap:"visible"`
		WorkersNum              int           `yaml:"workers-num"   env:"REPORT_POOL_SIZE" env-default:"50" debugmap:"visible"`
		RecordsBufferSize       uint64        `yaml:"records-buffer-size"   env:"REPORT_RECORDS_BUFFER_SIZE" env-default:"2000" debugmap:"visible"`
		FlushInterval           time.Duration `yaml:"flush-interval"   env:"REPORT_FLUSH_INTERVAL" env-default:"200ms" debugmap:"visible"`
		EnableDetailedRecording bool          `yaml:"enable-detailed-recording"   env:"REPORT_ENABLE_DETAILED_RECORDING" env-default:"true" debugmap:"visible"`
	}

	Download struct {
		PurgeDelay time.Duration `yaml:"purge-delay"   env:"DOWNLOAD_PURGE_DELAY" env-default:"10s" debugmap:"visible"`
		Backends   []string      `yaml:"backends"   env:"DOWNLOAD_BACKENDS" env-default:"" debugmap:"visible"`
	}

	Backends struct {
		CSV CSVOpt `yaml:"csv"   env:"BACKENDS_CSV" env-default:"" debugmap:"visible"`
	}

	CSVOpt struct {
		CSVDIR string `yaml:"csv_dir"   env:"CSV_DIR" env-default:"" debugmap:"visible"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
