package config

//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . Config App Log Feature DataStore Mysql Postgres Upload Download Backends CSVOpt EnvConfig

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	// FileConfig 用于描述配置文件
	FileConfig struct {
		App       `yaml:"app"`
		Log       `yaml:"logger"`
		Feature   `yaml:"feature"`
		DataStore `yaml:"datastore"`
		Mysql     `yaml:"mysql"`
		Postgres  `yaml:"psql"`
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
)

// NewConfigFromFile returns app config.
func NewConfigFromFile() (*FileConfig, error) {
	cfg := &FileConfig{}

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

type EnvConfig struct {
	TaskID string `env-required:"true" env:"TASK_ID" debugmap:"visible"`
	LogDir string `env-required:"true" env:"LOG_DIR" debugmap:"visible"`
}

func NewConfigFromEnv() (*EnvConfig, error) {
	var cfg EnvConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
