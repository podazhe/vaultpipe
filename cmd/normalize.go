package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpipe/internal/vault"
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize secret keys and values from a Vault path",
	Long: `Read secrets from Vault and apply normalization rules such as
uppercasing keys, trimming whitespace from values, replacing hyphens with
underscores, and stripping non-printable characters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			return fmt.Errorf("--path is required")
		}

		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		upperKeys, _ := cmd.Flags().GetBool("uppercase-keys")
		trimValues, _ := cmd.Flags().GetBool("trim-values")
		replaceHyphens, _ := cmd.Flags().GetBool("replace-hyphens")
		stripNonPrint, _ := cmd.Flags().GetBool("strip-non-print")
		showChanges, _ := cmd.Flags().GetBool("show-changes")

		opts := vault.NormalizeOptions{
			UppercaseKeys:  upperKeys,
			TrimValues:     trimValues,
			ReplaceHyphens: replaceHyphens,
			StripNonPrint:  stripNonPrint,
		}

		res := vault.NormalizeSecrets(secrets, opts)

		if showChanges {
			if len(res.Changes) == 0 {
				fmt.Fprintln(os.Stderr, "no changes")
			} else {
				for _, c := range res.Changes {
					if c.OldKey != "" {
						fmt.Fprintf(os.Stderr, "key renamed: %s -> %s\n", c.OldKey, c.Key)
					}
					if c.OldValue != "" {
						fmt.Fprintf(os.Stderr, "value changed for %s\n", c.Key)
					}
				}
			}
		}

		for k, v := range res.Normalized {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	normalizeCmd.Flags().String("path", "", "Vault secret path (required)")
	normalizeCmd.Flags().Bool("uppercase-keys", false, "Convert all keys to uppercase")
	normalizeCmd.Flags().Bool("trim-values", false, "Trim leading/trailing whitespace from values")
	normalizeCmd.Flags().Bool("replace-hyphens", false, "Replace hyphens with underscores in keys")
	normalizeCmd.Flags().Bool("strip-non-print", false, "Strip non-printable characters from values")
	normalizeCmd.Flags().Bool("show-changes", false, "Print a summary of changes to stderr")
	rootCmd.AddCommand(normalizeCmd)
}
