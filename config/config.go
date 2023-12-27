package config

//go:generate go run github.com/ecordell/optgen -output zz_generated.options.go . Config App Log Feature PG

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	// Config 用于描述配置文件
	Config struct {
		App     `yaml:"app"`
		Log     `yaml:"logger"`
		PG      `yaml:"postgres"`
		Feature `yaml:"feature"`
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

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX" debugmap:"visible"`
		URE     string `yaml:"url" env:"PG_URL" debugmap:"visible"`
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
