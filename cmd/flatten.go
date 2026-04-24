package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/vaultpipe/internal/vault"
)

var flattenCmd = &cobra.Command{
	Use:   "flatten",
	Short: "Flatten a nested JSON secret map into a single-level key=value set",
	Long: `Reads a nested JSON object from stdin and emits a flat map where
nested keys are joined by a separator (default "_").

Example:
  echo '{"db":{"host":"localhost","port":"5432"}}' | vaultpipe flatten --upper`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sep, _ := cmd.Flags().GetString("separator")
		prefix, _ := cmd.Flags().GetString("prefix")
		upper, _ := cmd.Flags().GetBool("upper")

		var nested map[string]any
		dec := json.NewDecoder(os.Stdin)
		if err := dec.Decode(&nested); err != nil {
			return fmt.Errorf("flatten: failed to decode JSON input: %w", err)
		}

		flat, err := vault.FlattenSecrets(nested, vault.FlattenOptions{
			Separator: sep,
			Prefix:    prefix,
			UpperCase: upper,
		})
		if err != nil {
			return err
		}

		outFmt, _ := cmd.Flags().GetString("output")
		switch outFmt {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(flat)
		default:
			for k, v := range flat {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			}
		}
		return nil
	},
}

func init() {
	flattenCmd.Flags().String("separator", "_", "Separator used between nested key segments")
	flattenCmd.Flags().String("prefix", "", "Prefix prepended to every key")
	flattenCmd.Flags().Bool("upper", false, "Convert all keys to upper case")
	flattenCmd.Flags().StringP("output", "o", "dotenv", "Output format: dotenv or json")
	rootCmd.AddCommand(flattenCmd)
}
