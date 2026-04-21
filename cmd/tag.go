package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/example/vaultpipe/internal/vault"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag secrets and filter by tag",
	Long:  `Assign tags to secrets read from Vault and optionally filter output by tag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := viper.GetString("tag.path")
		if path == "" {
			return fmt.Errorf("--path is required")
		}

		client, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		ts := vault.NewTaggedSecrets(secrets)

		rawTags, _ := cmd.Flags().GetStringArray("tag")
		for _, spec := range rawTags {
			parts := strings.SplitN(spec, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid tag spec %q: expected key=tag", spec)
			}
			if err := ts.Tag(parts[0], parts[1]); err != nil {
				return err
			}
		}

		filterTags, _ := cmd.Flags().GetStringArray("filter")
		result := ts.Secrets
		if len(filterTags) > 0 {
			result = ts.FilterByTag(filterTags...)
		}

		for k, v := range result {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	tagCmd.Flags().String("path", "", "Vault secret path to read")
	tagCmd.Flags().StringArray("tag", []string{}, "Tag assignment in key=tag format (repeatable)")
	tagCmd.Flags().StringArray("filter", []string{}, "Return only secrets matching all given tags (repeatable)")
	_ = viper.BindPFlag("tag.path", tagCmd.Flags().Lookup("path"))
	rootCmd.AddCommand(tagCmd)
}
