package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Enforce key/value quota limits on secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		maxKeys, _ := cmd.Flags().GetInt("max-keys")
		maxValueLen, _ := cmd.Flags().GetInt("max-value-len")
		maxTotalSize, _ := cmd.Flags().GetInt("max-total-size")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		client, err := vault.NewClient(cfgFile)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(cmd.Context(), client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		opts := vault.QuotaOptions{
			MaxKeys:      maxKeys,
			MaxValueLen:  maxValueLen,
			MaxTotalSize: maxTotalSize,
			DryRun:       dryRun,
		}

		result, err := vault.EnforceQuota(secrets, opts)

		fmt.Fprintf(os.Stdout, "keys=%d total_size=%d violations=%d\n",
			result.TotalKeys, result.TotalSize, len(result.Violations))

		for _, v := range result.Violations {
			fmt.Fprintf(os.Stderr, "  [%s] %s\n", v.Rule, v.Detail)
		}

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	quotaCmd.Flags().String("path", "", "Vault secret path (required)")
	_ = quotaCmd.MarkFlagRequired("path")
	quotaCmd.Flags().Int("max-keys", 0, "Maximum number of secret keys (0 = unlimited)")
	quotaCmd.Flags().Int("max-value-len", 0, "Maximum length of any single value (0 = unlimited)")
	quotaCmd.Flags().Int("max-total-size", 0, "Maximum total byte size of all values (0 = unlimited)")
	quotaCmd.Flags().Bool("dry-run", false, "Report all violations without failing")
	rootCmd.AddCommand(quotaCmd)
}
