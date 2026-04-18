package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpipe/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <path>",
	Short: "Show changes in a Vault secret path compared to a local env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, _ := cmd.Flags().GetString("env-file")

		client, err := vault.NewClient(cfgFile)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		next, err := vault.ReadSecrets(client, args[0])
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		old := map[string]string{}
		if envFile != "" {
			old, err = vault.LoadEnvFile(envFile)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("load env file: %w", err)
			}
		}

		diff := vault.DiffSecrets(old, next)
		if !diff.HasChanges() {
			fmt.Println("No changes detected.")
			return nil
		}

		fmt.Println(diff.Summary())
		for k := range diff.Added {
			fmt.Printf("+ %s\n", k)
		}
		for k := range diff.Removed {
			fmt.Printf("- %s\n", k)
		}
		for k := range diff.Changed {
			fmt.Printf("~ %s\n", k)
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().String("env-file", "", "Path to existing .env file to compare against")
	rootCmd.AddCommand(diffCmd)
}
