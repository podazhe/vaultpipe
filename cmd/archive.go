package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"vaultpipe/internal/vault"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a snapshot of secrets from Vault to disk",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := viper.GetString("archive.path")
		dir := viper.GetString("archive.dir")
		note := viper.GetString("archive.note")

		if path == "" {
			return fmt.Errorf("--path is required")
		}
		if dir == "" {
			dir = ".vaultpipe-archives"
		}

		client, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		am, err := vault.NewArchiveManager(dir)
		if err != nil {
			return fmt.Errorf("archive manager: %w", err)
		}

		entry, err := am.Save(path, secrets, note)
		if err != nil {
			return fmt.Errorf("save archive: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Archived %d secrets from %q at %s\n",
			len(entry.Secrets), entry.Path, entry.ArchivedAt.Format("2006-01-02T15:04:05Z"))
		return nil
	},
}

func init() {
	archiveCmd.Flags().String("path", "", "Vault secret path to archive")
	archiveCmd.Flags().String("dir", ".vaultpipe-archives", "Directory to store archive files")
	archiveCmd.Flags().String("note", "", "Optional note to attach to the archive entry")

	_ = viper.BindPFlag("archive.path", archiveCmd.Flags().Lookup("path"))
	_ = viper.BindPFlag("archive.dir", archiveCmd.Flags().Lookup("dir"))
	_ = viper.BindPFlag("archive.note", archiveCmd.Flags().Lookup("note"))

	rootCmd.AddCommand(archiveCmd)
}
