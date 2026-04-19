package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	rollbackVersion int
	rollbackMaxSize int
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [version]",
	Short: "Rollback secrets to a previously recorded version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}

		snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
		if snapshotDir == "" {
			snapshotDir = "."
		}

		fmt.Fprintf(os.Stdout, "Rolling back to version %d from snapshot dir %s\n", version, snapshotDir)
		// In a full implementation: load RollbackManager state persisted via
		// SnapshotManager, call Rollback(version), then write env via Writer.
		_ = version
		return nil
	},
}

func init() {
	rollbackCmd.Flags().String("snapshot-dir", ".", "Directory containing rollback snapshots")
	rollbackCmd.Flags().IntVar(&rollbackMaxSize, "max-history", 10, "Maximum number of versions to retain")
	rootCmd.AddCommand(rollbackCmd)
}
