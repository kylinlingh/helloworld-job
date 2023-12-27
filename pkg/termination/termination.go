package termination

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"helloworld/pkg/cobrautil"
	"helloworld/pkg/commonerrors"
	log "helloworld/pkg/logger"
	"os"
	"path"
)

const (
	terminationLogFlagName = "termination-log-path"
)

// PublishError returns a new wrapping cobra run function that executes the provided argument runFunc, and
// writes to disk an error returned by the latter if it is of type termination.TerminationError.
func PublishError(runFunc cobrautil.CobraRunFunc) cobrautil.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) error {
		runFuncErr := runFunc(cmd, args)
		var termErr commonerrors.TerminationError
		if runFuncErr != nil && errors.As(runFuncErr, &termErr) {
			ctx := context.Background()
			if cmd.Context() != nil {
				ctx = cmd.Context()
			}
			terminationLogPath, err := cmd.Flags().GetString(terminationLogFlagName)
			if err != nil || terminationLogPath == "" {
				terminationLogPath = "./termination.log"
			}
			//terminationLogPath := cobrautil.MustGetString(cmd, terminationLogFlagName)
			//if terminationLogPath == "" {
			//	return runFuncErr
			//}
			bytes, err := json.Marshal(termErr)
			if err != nil {
				log.Ctx(ctx).Error().Err(fmt.Errorf("unable to marshall termination log: %w", err)).Msg("failed to report termination log")
				return runFuncErr
			}

			if _, err := os.Stat(path.Dir(terminationLogPath)); os.IsNotExist(err) {
				mkdirErr := os.MkdirAll(path.Dir(terminationLogPath), 0o700) // Create your file
				if mkdirErr != nil {
					log.Ctx(ctx).Error().Err(fmt.Errorf("unable to create directory for termination log: %w", err)).Msg("failed to report termination log")
					return runFuncErr
				}
			}
			if err := os.WriteFile(terminationLogPath, bytes, 0o600); err != nil {
				log.Ctx(ctx).Error().Err(fmt.Errorf("unable to write terminationlog file: %w", err)).Msg("failed to report termination log")
			}
		}
		return runFuncErr
	}
}
