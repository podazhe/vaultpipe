package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Apply a named pipeline of transformations to secrets",
	Long: `chain reads secrets from Vault and applies a sequence of built-in
transformation steps (prefix, redact, upper, lower, filter) in the order
specified by --steps.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		stepNames, _ := cmd.Flags().GetStringSlice("steps")
		pfx, _ := cmd.Flags().GetString("prefix")

		client, err := vault.NewClient(cfgFile)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		chain := vault.NewChain()
		for _, name := range stepNames {
			switch strings.ToLower(name) {
			case "prefix":
				p := pfx
				chain.Add("prefix", func(m map[string]string) (map[string]string, error) {
					out := make(map[string]string, len(m))
					for k, v := range m {
						out[p+k] = v
					}
					return out, nil
				})
			case "redact":
				chain.Add("redact", func(m map[string]string) (map[string]string, error) {
					return vault.RedactSecrets(m, vault.RedactOptions{}), nil
				})
			case "upper":
				chain.Add("upper", func(m map[string]string) (map[string]string, error) {
					out := make(map[string]string, len(m))
					for k, v := range m {
						out[k] = strings.ToUpper(v)
					}
					return out, nil
				})
			default:
				return fmt.Errorf("unknown step: %q", name)
			}
		}

		result, err := chain.Run(secrets)
		if err != nil {
			return err
		}

		for k, v := range result {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	chainCmd.Flags().String("path", "", "Vault secret path (required)")
	chainCmd.Flags().StringSlice("steps", []string{}, "Ordered list of steps: prefix,redact,upper")
	chainCmd.Flags().String("prefix", "APP_", "Prefix to apply when using the prefix step")
	_ = chainCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(chainCmd)
}
