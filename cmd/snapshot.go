package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vaultpipe/vaultpipe/internal/vault"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Save or load a snapshot of secrets from Vault",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save [path]",
	Short: "Save current secrets at path to a local snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(viper.GetString("vault_addr"), viper.GetString("vault_token"))
		if err != nil {
			return err
		}
		data, err := vault.ReadSecrets(client, args[0])
		if err != nil {
			return err
		}
		dir := viper.GetString("snapshot_dir")
		if dir == "" {
			dir = ".vaultpipe/snapshots"
		}
		m := vault.NewSnapshotManager(dir)
		if err := m.Save(args[0], data); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Snapshot saved for %s\n", args[0])
		return nil
	},
}

var snapshotShowCmd = &cobra.Command{
	Use:   "show [path]",
	Short: "Show a previously saved snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := viper.GetString("snapshot_dir")
		if dir == "" {
			dir = ".vaultpipe/snapshots"
		}
		m := vault.NewSnapshotManager(dir)
		snap, err := m.Load(args[0])
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Path:\t%s\n", snap.Path)
		fmt.Fprintf(w, "Captured:\t%s\n", snap.CapturedAt.Format("2006-01-02 15:04:05 UTC"))
		fmt.Fprintln(w, "---")
		for k, v := range snap.Data {
			fmt.Fprintf(w, "%s\t%s\n", k, v)
		}
		return w.Flush()
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotShowCmd)
	snapshotCmd.PersistentFlags().String("snapshot-dir", ".vaultpipe/snapshots", "Directory to store snapshots")
	_ = viper.BindPFlag("snapshot_dir", snapshotCmd.PersistentFlags().Lookup("snapshot-dir"))
	rootCmd.AddCommand(snapshotCmd)
}
