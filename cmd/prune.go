package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpipe/internal/vault"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove unwanted secrets by prefix, key, or empty value",
	RunE:  runPrune,
}

var (
	pruneRemoveEmpty   bool
	pruneRemovePrefixes []string
	pruneRemoveKeys    []string
	pruneDryRun        bool
	prunePath          string
)

func init() {
	pruneCmd.Flags().StringVar(&prunePath, "path", "", "Vault KV path to read secrets from (required)")
	pruneCmd.Flags().BoolVar(&pruneRemoveEmpty, "remove-empty", false, "Remove secrets with empty values")
	pruneCmd.Flags().StringSliceVar(&pruneRemovePrefixes, "prefix", nil, "Remove secrets matching these key prefixes")
	pruneCmd.Flags().StringSliceVar(&pruneRemoveKeys, "key", nil, "Remove secrets with these exact keys")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "Report removals without modifying secrets")
	_ = pruneCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, _ []string) error {
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("prune: vault client: %w", err)
	}

	secrets, err := vault.ReadSecrets(client, prunePath)
	if err != nil {
		return fmt.Errorf("prune: read secrets: %w", err)
	}

	opts := vault.PruneOptions{
		RemoveEmpty:    pruneRemoveEmpty,
		RemovePrefixes: pruneRemovePrefixes,
		RemoveKeys:     pruneRemoveKeys,
		DryRun:         pruneDryRun,
	}

	_, result, err := vault.PruneSecrets(secrets, opts)
	if err != nil {
		return fmt.Errorf("prune: %w", err)
	}

	if pruneDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "[dry-run] "+result.Summary())
		if len(result.Removed) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  would remove: "+strings.Join(result.Removed, ", "))
		}
		return nil
	}

	fmt.Fprintln(os.Stdout, result.Summary())
	return nil
}
