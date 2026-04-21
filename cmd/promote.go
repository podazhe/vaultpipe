package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Promote secrets from one Vault path to another",
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := viper.GetString("promote.src")
		dstPath := viper.GetString("promote.dst")
		overwrite := viper.GetBool("promote.overwrite")
		dryRun := viper.GetBool("promote.dry-run")
		keys := viper.GetStringSlice("promote.keys")

		if srcPath == "" || dstPath == "" {
			return fmt.Errorf("--src and --dst are required")
		}

		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		src, err := client.ReadSecrets(srcPath)
		if err != nil {
			return fmt.Errorf("read src %s: %w", srcPath, err)
		}

		dst, err := client.ReadSecrets(dstPath)
		if err != nil {
			return fmt.Errorf("read dst %s: %w", dstPath, err)
		}

		opts := vault.PromoteOptions{
			Overwrite: overwrite,
			DryRun:    dryRun,
			Keys:      keys,
		}

		_, result := vault.PromoteSecrets(src, dst, opts)

		if dryRun {
			fmt.Printf("[dry-run] would promote: %s\n", strings.Join(result.Promoted, ", "))
			fmt.Printf("[dry-run] would skip:    %s\n", strings.Join(result.Skipped, ", "))
		} else {
			fmt.Println(result.Summary())
		}
		return nil
	},
}

func init() {
	promoteCmd.Flags().String("src", "", "Source Vault path")
	promoteCmd.Flags().String("dst", "", "Destination Vault path")
	promoteCmd.Flags().Bool("overwrite", false, "Overwrite existing keys in destination")
	promoteCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	promoteCmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys to promote")

	_ = viper.BindPFlag("promote.src", promoteCmd.Flags().Lookup("src"))
	_ = viper.BindPFlag("promote.dst", promoteCmd.Flags().Lookup("dst"))
	_ = viper.BindPFlag("promote.overwrite", promoteCmd.Flags().Lookup("overwrite"))
	_ = viper.BindPFlag("promote.dry-run", promoteCmd.Flags().Lookup("dry-run"))
	_ = viper.BindPFlag("promote.keys", promoteCmd.Flags().Lookup("keys"))

	rootCmd.AddCommand(promoteCmd)
}
