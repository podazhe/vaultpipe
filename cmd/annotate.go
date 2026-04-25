package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Attach metadata annotations to secrets and filter or inspect them",
	RunE: func(cmd *cobra.Command, args []string) error {
		secretPath, _ := cmd.Flags().GetString("path")
		annotationFlag, _ := cmd.Flags().GetStringArray("annotation")
		filterKey, _ := cmd.Flags().GetString("filter-key")
		filterVal, _ := cmd.Flags().GetString("filter-value")
		requirePresent, _ := cmd.Flags().GetBool("require-present")

		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, secretPath)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		// Parse --annotation KEY:ANNKEY=ANNVAL flags.
		annotations := make(map[string][]vault.Annotation)
		for _, raw := range annotationFlag {
			// format: SECRET_KEY:ann_key=ann_value
			parts := strings.SplitN(raw, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid annotation format %q, expected KEY:ann=val", raw)
			}
			secretKey := parts[0]
			kv := strings.SplitN(parts[1], "=", 2)
			if len(kv) != 2 {
				return fmt.Errorf("invalid annotation value %q, expected ann=val", parts[1])
			}
			annotations[secretKey] = append(annotations[secretKey], vault.Annotation{
				Key:   kv[0],
				Value: kv[1],
			})
		}

		annotated, err := vault.AnnotateSecrets(secrets, annotations, requirePresent)
		if err != nil {
			return err
		}

		if filterKey != "" {
			annotated = vault.FilterByAnnotation(annotated, filterKey, filterVal)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(annotated)
	},
}

func init() {
	annotateCmd.Flags().String("path", "", "Vault secret path (required)")
	_ = annotateCmd.MarkFlagRequired("path")
	annotateCmd.Flags().StringArray("annotation", nil, "Annotation in KEY:ann_key=ann_value format (repeatable)")
	annotateCmd.Flags().String("filter-key", "", "Filter output to secrets with this annotation key")
	annotateCmd.Flags().String("filter-value", "", "Match value substring when filtering by annotation key")
	annotateCmd.Flags().Bool("require-present", false, "Error if an annotated key is absent from secrets")
	rootCmd.AddCommand(annotateCmd)
}
