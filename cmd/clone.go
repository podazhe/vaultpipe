package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <src-path> <dst-path>",
	Short: "Clone secrets from one Vault path into another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("clone: vault client: %w", err)
		}

		src, err := client.ReadSecrets(args[0])
		if err != nil {
			return fmt.Errorf("clone: read src %q: %w", args[0], err)
		}

		dst, err := client.ReadSecrets(args[1])
		if err != nil {
			dst = map[string]string{}
		}

		prefix := viper.GetString("clone.prefix")
		overwrite := viper.GetBool("clone.overwrite")
		dryRun := viper.GetBool("clone.dry-run")
		keys := viper.GetStringSlice("clone.keys")

		res, err := vault.CloneSecrets(src, dst, vault.CloneOptions{
			Prefix:    prefix,
			KeyFilter: keys,
			Overwrite: overwrite,
			DryRun:    dryRun,
		})
		if err != nil {
			return fmt.Errorf("clone: %w", err)
		}

		if dryRun {
			fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] %s\n", res.Summary())
			if len(res.Cloned) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  would clone: %s\n", strings.Join(res.Cloned, ", "))
			}
			return nil
		}

		if err := client.WriteSecrets(args[1], dst); err != nil {
			return fmt.Errorf("clone: write dst %q: %w", args[1], err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", res.Summary())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().String("prefix", "", "prefix to prepend to cloned keys")
	cloneCmd.Flags().StringSlice("keys", nil, "restrict clone to these keys")
	cloneCmd.Flags().Bool("overwrite", false, "overwrite existing keys in destination")
	cloneCmd.Flags().Bool("dry-run", false, "show what would be cloned without writing")
	viper.BindPFlag("clone.prefix", cloneCmd.Flags().Lookup("prefix"))     //nolint:errcheck
	viper.BindPFlag("clone.keys", cloneCmd.Flags().Lookup("keys"))         //nolint:errcheck
	viper.BindPFlag("clone.overwrite", cloneCmd.Flags().Lookup("overwrite")) //nolint:errcheck
	viper.BindPFlag("clone.dry-run", cloneCmd.Flags().Lookup("dry-run"))   //nolint:errcheck
}
