package main

import (
	"os"

	"github.com/TykTechnologies/governance-action/pkg/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Configure production logger with console encoding for clean CI output
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.DisableStacktrace = true // Disable stack traces in logs

	logger, _ := config.Build()
	defer logger.Sync()

	rootCmd := &cobra.Command{
		Use:   "governance-action",
		Short: "Governance CI Action for analyzing OpenAPI specifications",
		Long: `A CI action that analyzes OpenAPI specifications against governance rules.
This action can be used in GitHub Actions and GitLab CI to ensure API compliance.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.RunAction(logger)
		},
		// Disable help text on error for cleaner CI output
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Action failed", zap.Error(err))
		os.Exit(1)
	}
}
