package main

import (
	"errors"
	"github.com/spf13/cobra"
	"helloworld/internal/cmd"
	"helloworld/internal/job"
	"helloworld/pkg/commonerrors"
	log "helloworld/pkg/logger"
	"os"
)

var errParsing = errors.New("parsing error")

func main() {
	//加载配置文件，生成通用的配置文件
	baseCfg := job.NewRunConfig()

	// 重新初始化 log
	log.New(baseCfg.Log.Level, baseCfg.App.RunMode, baseCfg.LogDir, baseCfg.TaskID)
	// 如果程序 panic 了，将信息重定向到日志中
	log.RewriteStderrFile(baseCfg.LogDir)

	// 注册 root command 和回调函数，在解析命令参数失败时回调
	rootCmd := cmd.NewRootCommand(baseCfg.App.Name)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return errParsing
	})

	// 注册某些通用 flag 到 root command 上
	cmd.RegisterRootFlags(rootCmd, nil)

	// 用于新建项目后尝试运行，检测是否一切ok
	tryConf := job.NewTryJobConfig(baseCfg)
	tryCmd := cmd.NewTryRunCommand(rootCmd.Use, tryConf)
	rootCmd.AddCommand(tryCmd)

	// 增加 scan 命令
	scanConf := job.NewScanJobConfig(baseCfg)
	scanCmd := cmd.NewScanCommand(rootCmd.Use, scanConf)
	cmd.RegisterScanFlags(scanCmd, scanConf)
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
