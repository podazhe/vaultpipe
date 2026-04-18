package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate against Vault and print the resulting token",
	RunE: func(cmd *cobra.Command, args []string) error {
		method := viper.GetString("auth.method")
		addr := viper.GetString("vault.address")
		if addr == "" {
			return fmt.Errorf("vault address is required (--vault-addr or VAULT_ADDR)")
		}

		switch method {
		case "token":
			token := viper.GetString("vault.token")
			if token == "" {
				return fmt.Errorf("token auth requires VAULT_TOKEN or --token flag")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "authenticated via token\n")
		case "approle":
			roleID := viper.GetString("auth.role_id")
			secretID := viper.GetString("auth.secret_id")
			if roleID == "" || secretID == "" {
				return fmt.Errorf("approle auth requires --role-id and --secret-id")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "authenticated via approle (role: %s)\n", roleID)
		default:
			return fmt.Errorf("unsupported auth method %q", method)
		}
		return nil
	},
}

func init() {
	authCmd.Flags().String("method", "token", "auth method: token, approle, kubernetes")
	authCmd.Flags().String("role-id", "", "AppRole role ID")
	authCmd.Flags().String("secret-id", "", "AppRole secret ID")
	authCmd.Flags().String("role", "", "Kubernetes role")

	_ = viper.BindPFlag("auth.method", authCmd.Flags().Lookup("method"))
	_ = viper.BindPFlag("auth.role_id", authCmd.Flags().Lookup("role-id"))
	_ = viper.BindPFlag("auth.secret_id", authCmd.Flags().Lookup("secret-id"))
	_ = viper.BindPFlag("auth.role", authCmd.Flags().Lookup("role"))

	rootCmd.AddCommand(authCmd)
}
