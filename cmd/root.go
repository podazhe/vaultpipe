package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd is the base command for the vaultpipe CLI.
var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Stream secrets from HashiCorp Vault into your environment",
	Long: `vaultpipe fetches secrets from HashiCorp Vault and injects them
as environment variables — either written to a .env file or passed
directly to a subprocess at runtime.`,
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "",
		"config file (default: $HOME/.vaultpipe.yaml or ./vaultpipe.yaml)",
	)

	// Vault connection flags shared across subcommands.
	rootCmd.PersistentFlags().String("vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().String("vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().String("vault-role-id", "", "AppRole role ID (overrides VAULT_ROLE_ID)")
	rootCmd.PersistentFlags().String("vault-secret-id", "", "AppRole secret ID (overrides VAULT_SECRET_ID)")

	_ = viper.BindPFlag("vault.addr", rootCmd.PersistentFlags().Lookup("vault-addr"))
	_ = viper.BindPFlag("vault.token", rootCmd.PersistentFlags().Lookup("vault-token"))
	_ = viper.BindPFlag("vault.role_id", rootCmd.PersistentFlags().Lookup("vault-role-id"))
	_ = viper.BindPFlag("vault.secret_id", rootCmd.PersistentFlags().Lookup("vault-secret-id"))
}

// initConfig reads in the config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "warning: could not determine home directory:", err)
		}

		// Search order: current directory, then home directory.
		viper.AddConfigPath(".")
		if home != "" {
			viper.AddConfigPath(home)
		}
		viper.SetConfigName(".vaultpipe")
		viper.SetConfigType("yaml")
	}

	// Allow all config keys to be overridden by environment variables
	// prefixed with VAULTPIPE_ (e.g. VAULTPIPE_VAULT_ADDR).
	viper.SetEnvPrefix("VAULTPIPE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
