package main

import (
	"errors"
	"github.com/spf13/cobra"
	"helloworld/config"
	"helloworld/internal/cmd"
	"helloworld/pkg/commonerrors"
	log "helloworld/pkg/logger"
	"os"
)

var errParsing = errors.New("parsing error")

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Msgf("Initialization failed: %s", err)
	}

	// 重新初始化 log
	log.New(cfg.Log.Level, cfg.App.RunMode)

	// 注册 root command 和回调函数，在解析命令参数失败时回调
	rootCmd := cmd.NewRootCommand(cfg.App.Name)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return errParsing
	})

	// 注册某些通用 flag 到 root command 上
	cmd.RegisterRootFlags(rootCmd, nil)

	// 增加 scan 命令
	scanCmd := cmd.NewScanCommand(rootCmd.Use, cfg)
	cmd.RegisterScanFlags(scanCmd)
	rootCmd.AddCommand(scanCmd)

	// 运行命令
	if err := rootCmd.Execute(); err != nil {
		if !errors.Is(err, errParsing) {
			log.Err(err).Msg("terminated with errors")
		}
		var termErr commonerrors.TerminationError
		if errors.As(err, &termErr) {
			os.Exit(-1)
		}
		os.Exit(1)
	}
}
