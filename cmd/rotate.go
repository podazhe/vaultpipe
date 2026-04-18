package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Watch Vault paths and re-fetch secrets on a schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault_addr"), viper.GetString("vault_token"))
		if err != nil {
			return err
		}

		interval, _ := cmd.Flags().GetDuration("interval")
		paths, _ := cmd.Flags().GetStringSlice("path")

		cache := vault.NewSecretCache()
		cfg := vault.RotateConfig{
			Paths:    paths,
			Interval: interval,
			OnRotate: func(path string, data map[string]interface{}) {
				cmd.Printf("[rotate] refreshed %s (%d keys)\n", path, len(data))
			},
		}

		rotator := vault.NewRotator(client, cache, cfg)

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		cmd.Println("Starting rotation loop — press Ctrl+C to stop")
		rotator.Start(ctx)
		return nil
	},
}

func init() {
	rotateCmd.Flags().StringSlice("path", []string{}, "Vault secret paths to watch (repeatable)")
	rotateCmd.Flags().Duration("interval", 60*time.Second, "Poll interval for secret rotation")
	rotateCmd.MarkFlagRequired("path") //nolint:errcheck
	rootCmd.AddCommand(rotateCmd)
}
