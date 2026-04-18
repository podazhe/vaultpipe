package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Apply key transformations (prefix, rename, filter) to secrets from Vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix, _ := cmd.Flags().GetString("prefix")
		filter, _ := cmd.Flags().GetStringSlice("filter")
		renameRaw, _ := cmd.Flags().GetStringSlice("rename")

		renames := make(map[string]string)
		for _, r := range renameRaw {
			parts := strings.SplitN(r, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid rename %q, expected old=new", r)
			}
			renames[parts[0]] = parts[1]
		}

		// Read secrets from stdin as JSON for composability.
		var secrets map[string]string
		if err := json.NewDecoder(os.Stdin).Decode(&secrets); err != nil {
			return fmt.Errorf("reading secrets from stdin: %w", err)
		}

		tr := vault.NewTransformer(vault.TransformRule{
			Prefix:  prefix,
			Renames: renames,
			Filter:  filter,
		})

		out, err := tr.Apply(secrets)
		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(out)
	},
}

func init() {
	transformCmd.Flags().String("prefix", "", "Prefix to prepend to all keys (uppercased)")
	transformCmd.Flags().StringSlice("filter", nil, "Only include these keys")
	transformCmd.Flags().StringSlice("rename", nil, "Rename keys in old=new format")
	rootCmd.AddCommand(transformCmd)
}
