package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	indexPath   string
	indexSearch string
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Build and query an in-memory index of secret keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, indexPath)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		idx := vault.NewSecretIndex()
		for k := range secrets {
			idx.Add(indexPath, k, nil, nil)
		}

		var entries []*vault.IndexEntry
		if indexSearch != "" {
			entries = idx.Search(indexSearch)
		} else {
			entries = idx.ListByPath(indexPath)
		}

		if len(entries) == 0 {
			fmt.Println("no entries found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PATH\tKEY")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\n", e.Path, e.Key)
		}
		return w.Flush()
	},
}

func init() {
	indexCmd.Flags().StringVar(&indexPath, "path", "", "Vault secret path to index (required)")
	indexCmd.Flags().StringVar(&indexSearch, "search", "", "Substring to search within indexed keys")
	_ = indexCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(indexCmd)
}
