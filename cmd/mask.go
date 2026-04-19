package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "Print secrets with values partially masked",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		showChars, _ := cmd.Flags().GetInt("show-chars")
		sensitive, _ := cmd.Flags().GetStringSlice("sensitive")
		maskChar, _ := cmd.Flags().GetString("mask-char")

		client, err := vault.NewClient(cfgFile)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(cmd.Context(), client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		opts := &vault.MaskOptions{
			ShowChars: showChars,
			MaskChar:  maskChar,
		}
		masked := vault.MaskSecrets(secrets, sensitive, opts)

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(masked); err != nil {
			return fmt.Errorf("encode output: %w", err)
		}
		return nil
	},
}

func init() {
	maskCmd.Flags().String("path", "", "Vault secret path (required)")
	_ = maskCmd.MarkFlagRequired("path")
	maskCmd.Flags().Int("show-chars", 4, "Number of trailing characters to reveal")
	maskCmd.Flags().StringSlice("sensitive", nil, "Keys to mask fully regardless of show-chars")
	maskCmd.Flags().String("mask-char", "*", "Character used for masking")
	rootCmd.AddCommand(maskCmd)
}
