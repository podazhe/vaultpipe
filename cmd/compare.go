package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var (
	compareLeftPath  string
	compareRightPath string
	compareShowKeys  bool
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare secrets between two Vault paths",
	Long:  `Reads secrets from two Vault KV paths and reports keys that differ, are missing, or are identical.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		left, err := vault.ReadSecrets(cmd.Context(), client, compareLeftPath)
		if err != nil {
			return fmt.Errorf("reading left path %q: %w", compareLeftPath, err)
		}

		right, err := vault.ReadSecrets(cmd.Context(), client, compareRightPath)
		if err != nil {
			return fmt.Errorf("reading right path %q: %w", compareRightPath, err)
		}

		result := vault.CompareSecrets(left, right)
		fmt.Print(result.Summary())

		if compareShowKeys {
			printKeys("Only in left", result.OnlyInLeft)
			printKeys("Only in right", result.OnlyInRight)
			printKeys("Different", result.Different)
		}

		if result.HasDifferences() {
			os.Exit(1)
		}
		return nil
	},
}

func printKeys(label string, keys []string) {
	if len(keys) == 0 {
		return
	}
	fmt.Printf("%s:\n", label)
	for _, k := range keys {
		fmt.Printf("  - %s\n", k)
	}
}

func init() {
	compareCmd.Flags().StringVar(&compareLeftPath, "left", "", "Left Vault secret path (required)")
	compareCmd.Flags().StringVar(&compareRightPath, "right", "", "Right Vault secret path (required)")
	compareCmd.Flags().BoolVar(&compareShowKeys, "show-keys", false, "Print the names of differing keys")
	_ = compareCmd.MarkFlagRequired("left")
	_ = compareCmd.MarkFlagRequired("right")
	rootCmd.AddCommand(compareCmd)
}
