package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	vaultclient "vaultpipe/internal/vault"
	"vaultpipe/internal/template"
)

var (
	tmplFile   string
	tmplPaths  []string
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render a template file with secrets from Vault",
	Long: `Reads secrets from one or more Vault KV paths and renders a
Go text/template file, writing the result to stdout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if tmplFile == "" {
			return fmt.Errorf("--template flag is required")
		}
		if len(tmplPaths) == 0 {
			return fmt.Errorf("at least one --path is required")
		}

		client, err := vaultclient.NewClient()
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets := make(map[string]string)
		for _, p := range tmplPaths {
			data, err := client.ReadSecrets(p)
			if err != nil {
				return fmt.Errorf("read %s: %w", p, err)
			}
			for k, v := range data {
				secrets[k] = v
			}
		}

		src, err := os.ReadFile(tmplFile)
		if err != nil {
			return fmt.Errorf("read template file: %w", err)
		}

		r := template.NewRenderer(os.Stdout)
		return r.Render(string(src), secrets)
	},
}

func init() {
	templateCmd.Flags().StringVarP(&tmplFile, "template", "t", "", "path to template file (required)")
	templateCmd.Flags().StringArrayVarP(&tmplPaths, "path", "p", nil, "Vault secret path (repeatable)")
	rootCmd.AddCommand(templateCmd)
}
