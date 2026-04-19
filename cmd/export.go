package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export Vault secrets to stdout in dotenv, JSON, or YAML format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(args[0])
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		format := vault.ExportFormat(viper.GetString("export.format"))
		opts := vault.ExportOptions{
			Format: format,
			Export: viper.GetBool("export.export_keyword"),
			Prefix: viper.GetString("export.prefix"),
		}

		out, err := vault.ExportSecrets(secrets, opts)
		if err != nil {
			return err
		}

		fmt.Fprint(os.Stdout, out)
		return nil
	},
}

func init() {
	exportCmd.Flags().String("format", "dotenv", "Output format: dotenv, json, yaml")
	exportCmd.Flags().Bool("export-keyword", false, "Prefix each line with 'export' (dotenv only)")
	exportCmd.Flags().String("prefix", "", "Optional key prefix to prepend to all keys")

	_ = viper.BindPFlag("export.format", exportCmd.Flags().Lookup("format"))
	_ = viper.BindPFlag("export.export_keyword", exportCmd.Flags().Lookup("export-keyword"))
	_ = viper.BindPFlag("export.prefix", exportCmd.Flags().Lookup("prefix"))

	rootCmd.AddCommand(exportCmd)
}
