package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var copyCmd = &cobra.Command{
	Use:   "copy <src-path> <dst-path>",
	Short: "Copy secrets from one Vault path to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		srcPath, dstPath := args[0], args[1]

		src, err := client.ReadSecrets(cmd.Context(), srcPath)
		if err != nil {
			return fmt.Errorf("read source %q: %w", srcPath, err)
		}

		dst, err := client.ReadSecrets(cmd.Context(), dstPath)
		if err != nil {
			// destination may not exist yet — start with empty map
			dst = map[string]string{}
		}

		opts := vault.CopyOptions{
			Overwrite: viper.GetBool("copy.overwrite"),
			Keys:      viper.GetStringSlice("copy.keys"),
			DryRun:    viper.GetBool("copy.dry-run"),
		}

		result, err := vault.CopySecrets(dst, src, opts)
		if err != nil {
			return err
		}

		if opts.DryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] %s\n", result.Summary())
			return nil
		}

		if err := client.WriteSecrets(cmd.Context(), dstPath, dst); err != nil {
			return fmt.Errorf("write destination %q: %w", dstPath, err)
		}

		fmt.Fprintln(os.Stdout, result.Summary())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().Bool("overwrite", false, "overwrite existing keys in the destination")
	copyCmd.Flags().StringSlice("keys", nil, "copy only these keys (comma-separated)")
	copyCmd.Flags().Bool("dry-run", false, "preview changes without writing to Vault")
	_ = viper.BindPFlag("copy.overwrite", copyCmd.Flags().Lookup("overwrite"))
	_ = viper.BindPFlag("copy.keys", copyCmd.Flags().Lookup("keys"))
	_ = viper.BindPFlag("copy.dry-run", copyCmd.Flags().Lookup("dry-run"))
}
