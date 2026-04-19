package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	vaultpkg "vaultpipe/internal/vault"
)

func init() {
	var paths []string
	var required []string
	var noEmpty bool
	var warnLong int

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate secrets fetched from Vault against key/value rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := vaultpkg.NewClient(viper.GetString("vault_addr"), viper.GetString("vault_token"))
			if err != nil {
				return fmt.Errorf("vault client: %w", err)
			}
			merged := map[string]string{}
			for _, p := range paths {
				secrets, err := vaultpkg.ReadSecrets(client, p)
				if err != nil {
					return fmt.Errorf("read %s: %w", p, err)
				}
				for k, v := range secrets {
					merged[k] = v
				}
			}
			opts := vaultpkg.ValidateOptions{
				NoEmpty:      noEmpty,
				WarnLong:     warnLong,
				RequiredKeys: required,
			}
			res := vaultpkg.ValidateSecrets(merged, opts)
			if s := res.Summary(); s != "" {
				fmt.Fprint(os.Stderr, s)
			}
			if !res.OK() {
				return fmt.Errorf("validation failed with %d error(s)", len(res.Errors))
			}
			fmt.Println("validation passed " + strings.Repeat(".", 3))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&paths, "path", nil, "Vault secret paths to validate")
	cmd.Flags().StringSliceVar(&required, "require", nil, "Keys that must be present")
	cmd.Flags().BoolVar(&noEmpty, "no-empty", false, "Fail if any value is empty")
	cmd.Flags().IntVar(&warnLong, "warn-long", 0, "Warn if value exceeds N characters (0 = off)")
	_ = cmd.MarkFlagRequired("path")

	rootCmd.AddCommand(cmd)
}
