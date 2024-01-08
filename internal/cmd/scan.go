package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"helloworld/internal/job"
	"helloworld/pkg/cmd"
	"helloworld/pkg/termination"
)

func NewScanCommand(programName string, runConfig *job.ParamConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "scan for compliance",
		Long:  "static code analysis for helloworld compliance",
		RunE: termination.PublishError(func(command *cobra.Command, args []string) error {
			//runConfig := job.NewRunConfig(config)
			runnableJob, err := runConfig.Complete(command.Context())
			if err != nil {
				return err
			}
			signalCtx := cmd.SignalContextWithGracePeriod(context.Background(), runConfig.Feature.ShutdownGracePeriod)
			return runnableJob.Run(signalCtx)
		}),
	}
}

func RegisterScanFlags(command *cobra.Command, cfg *job.ParamConfig) {
	command.Flags().StringSliceVarP(&cfg.ScanTargets, "target", "t", []string{}, "scan target(sc-android, sc-ios, af-android, af-ios)")
	command.Flags().StringVar(&cfg.ScanMode, "mode", "dev", "running mode: dev, pro; data will be exported to csv under dev mode, or will be exported to db under pro mode")
}
