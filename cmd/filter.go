package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vaultpipe/vaultpipe/internal/vault"
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter secrets by prefix, key, or exclusion list",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		path, _ := cmd.Flags().GetString("path")
		prefixes, _ := cmd.Flags().GetStringSlice("prefix")
		keys, _ := cmd.Flags().GetStringSlice("key")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		filtered := vault.FilterSecrets(secrets, vault.FilterOptions{
			Prefixes: prefixes,
			Keys:     keys,
			Exclude:  exclude,
		})

		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(filtered)
	},
}

func init() {
	filterCmd.Flags().String("path", "", "Vault secret path (required)")
	_ = filterCmd.MarkFlagRequired("path")
	filterCmd.Flags().StringSlice("prefix", nil, "Include keys matching these prefixes")
	filterCmd.Flags().StringSlice("key", nil, "Include only these exact keys")
	filterCmd.Flags().StringSlice("exclude", nil, "Exclude these keys")
	rootCmd.AddCommand(filterCmd)
}
