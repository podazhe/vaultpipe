package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/vaultpipe/internal/vault"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint secrets from Vault against style and safety rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		allowEmpty, _ := cmd.Flags().GetBool("allow-empty")
		maxLen, _ := cmd.Flags().GetInt("max-value-len")
		forbidPrefixes, _ := cmd.Flags().GetStringSlice("forbid-prefix")

		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		opts := vault.LintOptions{
			AllowEmpty:   allowEmpty,
			MaxValueLen:  maxLen,
			ForbidPrefix: forbidPrefixes,
		}

		results := vault.LintSecrets(secrets, opts)
		if len(results) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "✔ no lint violations found")
			return nil
		}

		for _, r := range results {
			fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s — %s\n", r.Rule, r.Key, r.Message)
		}

		fmt.Fprintf(os.Stderr, "lint failed: %d violation(s)\n", len(results))
		os.Exit(1)
		return nil
	},
}

func init() {
	lintCmd.Flags().String("path", "", "Vault secret path to lint (required)")
	_ = lintCmd.MarkFlagRequired("path")
	lintCmd.Flags().Bool("allow-empty", false, "Allow secrets with empty values")
	lintCmd.Flags().Int("max-value-len", 0, "Maximum allowed value length (0 = unlimited)")
	lintCmd.Flags().StringSlice("forbid-prefix", nil, "Key prefixes that are not allowed")
	rootCmd.AddCommand(lintCmd)
}
