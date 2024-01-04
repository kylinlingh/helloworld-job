package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helloworld/internal/job"
	log "helloworld/pkg/logger"
)

func NewRootCommand(programName string) *cobra.Command {
	return &cobra.Command{
		Use:   programName,
		Short: "",
		Long:  "",
		//Example:       job.ServeExample(programName),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

const configDescription = `config file path
order of precedence:
1. --config/-c
2. env var GITLEAKS_CONFIG
3. (--source/-s)/.config.yml
If none of the three options are used, then helloworld will use the default config`

// RegisterRootFlags register general flags for root command
func RegisterRootFlags(cmd *cobra.Command, config *job.ParamConfig) {

	cmd.PersistentFlags().StringP("config", "c", "", configDescription)
	cmd.PersistentFlags().StringP("source", "s", ".", "path to source")
	cmd.PersistentFlags().StringP("downloadfrom-path", "r", "", "downloadfrom file")
	cmd.PersistentFlags().StringP("downloadfrom-format", "f", "json", "output format (json, csv, junit, sarif)")
	cmd.PersistentFlags().StringP("baseline-path", "b", "", "path to baseline with issues that can be ignored")
	cmd.PersistentFlags().StringP("log-level", "l", "info", "log level (trace, debug, info, warn, error, fatal)")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "show verbose output from scan")
	cmd.PersistentFlags().Int("max-target-megabytes", 0, "files larger than this will be skipped")
	cmd.PersistentFlags().BoolP("ignore-gitleaks-allow", "", false, "ignore gitleaks:allow comments")
	//cmd.PersistentFlags().Uint("redact", 0, "redact secrets from logs and stdout. To redact only parts of the secret just apply a percent value from 0..100. For example --redact=20 (default 100%)")
	//cmd.Flag("redact").NoOptDefVal = "100"
	//cmd.PersistentFlags().Bool("no-banner", false, "suppress banner")
	cmd.PersistentFlags().String("log-opts", "", "git log options")
	//cmd.PersistentFlags().StringSlice("enable-rule", []string{}, "only enable specific rules by id, ex: `gitleaks detect --enable-rule=atlassian-api-token --enable-rule=slack-access-token`")
	cmd.PersistentFlags().StringP("gitleaks-ignore-path", "i", ".", "path to .gitleaksignore file or folder containing one")
	cmd.PersistentFlags().Bool("follow-symlinks", false, "scan files that are symlinks to other files")

	cmd.PersistentFlags().String("termination-log-path", "./termination.log", "define the path to the termination log file, which contains a JSON payload to surface as reason for termination")
	err := viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	if err != nil {
		log.Fatal().Msgf("err binding config %s", err.Error())
	}

}
