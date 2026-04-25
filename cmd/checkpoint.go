package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"vaultpipe/internal/vault"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Manage named secret checkpoints",
}

var checkpointSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current secrets as a named checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := viper.GetString("path")
		dir := viper.GetString("checkpoint-dir")
		client, err := vault.NewClient(nil)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}
		m, err := vault.NewCheckpointManager(dir)
		if err != nil {
			return err
		}
		if err := m.Save(args[0], secrets); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Checkpoint %q saved.\n", args[0])
		return nil
	},
}

var checkpointDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a named checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := viper.GetString("checkpoint-dir")
		m, err := vault.NewCheckpointManager(dir)
		if err != nil {
			return err
		}
		if err := m.Delete(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Checkpoint %q deleted.\n", args[0])
		return nil
	},
}

func init() {
	checkpointCmd.AddCommand(checkpointSaveCmd)
	checkpointCmd.AddCommand(checkpointDeleteCmd)

	checkpointCmd.PersistentFlags().String("checkpoint-dir", ".vaultpipe/checkpoints", "Directory to store checkpoints")
	_ = viper.BindPFlag("checkpoint-dir", checkpointCmd.PersistentFlags().Lookup("checkpoint-dir"))

	checkpointSaveCmd.Flags().String("path", "", "Vault secret path to snapshot")
	_ = viper.BindPFlag("path", checkpointSaveCmd.Flags().Lookup("path"))

	rootCmd.AddCommand(checkpointCmd)
}
