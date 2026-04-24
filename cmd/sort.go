package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"vaultpipe/internal/vault"
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort secrets from a Vault path and print them",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault-addr"), viper.GetString("token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		path := viper.GetString("sort.path")
		if path == "" {
			return fmt.Errorf("--path is required")
		}

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		order := vault.SortAsc
		if viper.GetString("sort.order") == "desc" {
			order = vault.SortDesc
		}

		opts := vault.SortOptions{
			Order:      order,
			ByValue:    viper.GetBool("sort.by-value"),
			IgnoreCase: viper.GetBool("sort.ignore-case"),
			Prefix:     viper.GetString("sort.prefix"),
		}

		pairs := vault.SortSecrets(secrets, opts)
		w := os.Stdout
		for _, kv := range pairs {
			fmt.Fprintf(w, "%s=%s\n", kv.Key, kv.Value)
		}
		return nil
	},
}

func init() {
	sortCmd.Flags().String("path", "", "Vault KV path to read secrets from")
	sortCmd.Flags().String("order", "asc", "Sort order: asc or desc")
	sortCmd.Flags().Bool("by-value", false, "Sort by secret value instead of key")
	sortCmd.Flags().Bool("ignore-case", false, "Case-insensitive sort")
	sortCmd.Flags().String("prefix", "", "Only include keys with this prefix")

	_ = viper.BindPFlag("sort.path", sortCmd.Flags().Lookup("path"))
	_ = viper.BindPFlag("sort.order", sortCmd.Flags().Lookup("order"))
	_ = viper.BindPFlag("sort.by-value", sortCmd.Flags().Lookup("by-value"))
	_ = viper.BindPFlag("sort.ignore-case", sortCmd.Flags().Lookup("ignore-case"))
	_ = viper.BindPFlag("sort.prefix", sortCmd.Flags().Lookup("prefix"))

	rootCmd.AddCommand(sortCmd)
}
