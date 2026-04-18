package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourorg/vaultpipe/internal/env"
	"github.com/yourorg/vaultpipe/internal/vault"
)

var (
	outputPath   string
	outputFormat string
)

var writeCmd = &cobra.Command{
	Use:   "write [secret-path]",
	Short: "Fetch secrets from Vault and write them to an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretPath := args[0]

		client, err := vault.NewClient(vault.Config{
			Address:   viper.GetString("vault_addr"),
			Token:     viper.GetString("vault_token"),
			RoleID:    viper.GetString("role_id"),
			SecretID:  viper.GetString("secret_id"),
		})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(secretPath)
		if err != nil {
			return fmt.Errorf("reading secrets: %w", err)
		}

		w, err := env.NewWriter(outputPath, env.Format(outputFormat))
		if err != nil {
			return fmt.Errorf("env writer: %w", err)
		}
		defer w.Close()

		if err := w.Write(secrets); err != nil {
			return fmt.Errorf("writing env: %w", err)
		}

		if outputPath != "" {
			log.Printf("secrets written to %s", outputPath)
		}
		return nil
	},
}

func init() {
	writeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output env file (default: stdout)")
	writeCmd.Flags().StringVarP(&outputFormat, "format", "f", "dotenv", "Output format: dotenv or export")
	rootCmd.AddCommand(writeCmd)
}
