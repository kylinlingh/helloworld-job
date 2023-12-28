package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"helloworld/config"
	"helloworld/internal/job"
	"helloworld/pkg/cmd"
	"helloworld/pkg/termination"
)

func NewScanCommand(programName string, config *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "scan for compliance",
		Long:  "static code analysis for helloworld compliance",
		RunE: termination.PublishError(func(command *cobra.Command, args []string) error {
			runConfig := job.NewRunConfig(config)
			runnableJob, err := runConfig.Complete(command.Context())
			if err != nil {
				return err
			}
			signalCtx := cmd.SignalContextWithGracePeriod(context.Background(), runConfig.Feature.ShutdownGracePeriod)
			return runnableJob.Run(signalCtx)
		}),
	}
}

func RegisterScanFlags(command *cobra.Command) {
	command.Flags().Bool("no-git", false, "treat git repo as a regular directory and scan those files, --log-opts has no effect on the scan when --no-git is set")
}
