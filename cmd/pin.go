package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	pinVersion int
	globalPinManager = vault.NewPinManager()
)

var pinCmd = &cobra.Command{
	Use:   "pin <path> <version>",
	Short: "Pin a Vault secret path to a specific version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		ver, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("version must be an integer: %w", err)
		}

		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		if err := globalPinManager.Pin(path, ver, secrets); err != nil {
			return fmt.Errorf("pin: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Pinned %s at version %d (%d keys)\n", path, ver, len(secrets))
		return nil
	},
}

var unpinCmd = &cobra.Command{
	Use:   "unpin <path>",
	Short: "Remove a pin for a Vault secret path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		globalPinManager.Unpin(args[0])
		fmt.Fprintf(cmd.OutOrStdout(), "Unpinned %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pinCmd)
	rootCmd.AddCommand(unpinCmd)
}
