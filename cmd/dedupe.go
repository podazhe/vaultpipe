package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/vault"
)

var dedupeCmd = &cobra.Command{
	Use:   "dedupe [path...]",
	Short: "Deduplicate secrets from one or more Vault paths",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
		preferLast, _ := cmd.Flags().GetBool("prefer-last")
		report, _ := cmd.Flags().GetBool("report")

		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		var sources []map[string]string
		for _, path := range args {
			secrets, err := vault.ReadSecrets(cmd.Context(), client, path)
			if err != nil {
				return fmt.Errorf("read %s: %w", path, err)
			}
			sources = append(sources, secrets)
		}

		res, err := vault.DedupeSecrets(sources, vault.DedupeOptions{
			CaseSensitive:    caseSensitive,
			PreferLast:       preferLast,
			ReportDuplicates: report,
		})
		if err != nil {
			return err
		}

		if report && len(res.Duplicates) > 0 {
			fmt.Fprintf(cmd.ErrOrStderr(), "duplicates removed: %s\n",
				strings.Join(res.Duplicates, ", "))
		}

		for k, v := range res.Secrets {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	dedupeCmd.Flags().Bool("case-sensitive", false, "treat keys as case-sensitive when deduplicating")
	dedupeCmd.Flags().Bool("prefer-last", false, "keep the last occurrence of a duplicate key")
	dedupeCmd.Flags().Bool("report", false, "print removed duplicate keys to stderr")
	rootCmd.AddCommand(dedupeCmd)
}
