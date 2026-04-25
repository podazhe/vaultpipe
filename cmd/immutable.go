package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	immutableLockKeys   []string
	immutableSecretPath string
)

var immutableCmd = &cobra.Command{
	Use:   "immutable",
	Short: "Lock specific secret keys to prevent modification",
	Long: `Read secrets from Vault, lock the specified keys, and print the
resulting snapshot. Locked keys cannot be overwritten or deleted within
the session, providing a lightweight guard against accidental mutation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newVaultClient()
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(immutableSecretPath)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		im := newImmutableFromVault(secrets)
		im.Lock(immutableLockKeys...)

		fmt.Fprintf(cmd.OutOrStdout(), "Locked keys: %s\n",
			strings.Join(im.LockedKeys(), ", "))

		snap := im.Snapshot()
		for _, k := range sortedMapKeys(snap) {
			marker := " "
			if im.IsLocked(k) {
				marker = "*"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s\n", marker, k)
		}
		return nil
	},
}

// newImmutableFromVault constructs an ImmutableSecrets from a raw map.
// Extracted to allow unit-testing without a live Vault.
func newImmutableFromVault(secrets map[string]string) *vaultImmutable {
	return newImmutableWrapper(secrets)
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort — maps are small in practice
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

func init() {
	immutableCmd.Flags().StringVar(&immutableSecretPath, "path", "", "Vault secret path to read (required)")
	immutableCmd.Flags().StringSliceVar(&immutableLockKeys, "lock", nil, "Comma-separated list of keys to lock")
	_ = immutableCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(immutableCmd)
}
