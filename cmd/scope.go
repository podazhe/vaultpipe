package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var (
	scopePrefix  string
	scopeTags    []string
	scopeSecrets []string
)

var scopeCmd = &cobra.Command{
	Use:   "scope <name>",
	Short: "Partition secrets into named scopes by key prefix",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tags := make(map[string]string)
		for _, t := range scopeTags {
			parts := strings.SplitN(t, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid tag %q: expected key=value", t)
			}
			tags[parts[0]] = parts[1]
		}

		sm := vault.NewScopeManager()
		if err := sm.Register(name, scopePrefix, tags); err != nil {
			return fmt.Errorf("register scope: %w", err)
		}

		secrets := make(map[string]string)
		for _, s := range scopeSecrets {
			parts := strings.SplitN(s, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid secret %q: expected KEY=value", s)
			}
			secrets[parts[0]] = parts[1]
		}

		buckets := sm.Partition(secrets)

		for bucket, kvs := range buckets {
			fmt.Fprintf(cmd.OutOrStdout(), "[%s]\n", bucket)
			for k, v := range kvs {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, v)
			}
		}
		return nil
	},
}

func init() {
	scopeCmd.Flags().StringVar(&scopePrefix, "prefix", "", "Key prefix that defines this scope")
	scopeCmd.Flags().StringArrayVar(&scopeTags, "tag", nil, "Metadata tags in key=value format")
	scopeCmd.Flags().StringArrayVar(&scopeSecrets, "secret", nil, "Secrets to partition in KEY=value format")
	rootCmd.AddCommand(scopeCmd)
}
