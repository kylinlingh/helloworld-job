package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"helloworld/internal/job"
	"helloworld/pkg/cmd"
	"helloworld/pkg/termination"
)

func NewTryRunCommand(programName string, runConfig *job.TryJobParamConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "try",
		Short: "try to run the template project",
		Long:  "try to run the template project",
		RunE: termination.PublishError(func(command *cobra.Command, args []string) error {
			//runConfig := job.NewRunConfig(config)
			runnableJob, err := runConfig.CompleteTryJob(command.Context())
			if err != nil {
				return err
			}
			signalCtx := cmd.SignalContextWithGracePeriod(context.Background(), runConfig.Feature.ShutdownGracePeriod)
			return runnableJob.Run(signalCtx)
		}),
	}
}
