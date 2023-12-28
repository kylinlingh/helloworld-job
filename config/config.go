package config

//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . Config App Log Feature DataStore Mysql

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
		Postgres  `yaml:"postgres"`
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
		Host                  string        `yaml:"host"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"visible"`
		Username              string        `yaml:"username"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"visible"`
		Password              string        `yaml:"password"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"sensitive"`
		Database              string        `yaml:"database"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"visible"`
		MaxIdleConnections    int           `yaml:"max_idle_connections"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"visible"`
		MaxOpenConnections    int           `yaml:"max_open_connections"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"visible"`
		MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"   env:"DATASTORE_ENGINE" env-default:"none" debugmap:"visible"`
	}

	Postgres struct {
		Uri string `yaml:"uri"   env:"POSTGRES_URI" env-default:"" debugmap:"visible"`
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
