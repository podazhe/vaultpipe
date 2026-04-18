package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var runPaths []string

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Inject Vault secrets as env vars and execute a command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(runPaths)
		if err != nil {
			return fmt.Errorf("reading secrets: %w", err)
		}

		// Build env: inherit current env then overlay secrets.
		env := os.Environ()
		for k, v := range secrets {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}

		binary, err := exec.LookPath(args[0])
		if err != nil {
			return fmt.Errorf("resolving binary %q: %w", args[0], err)
		}

		return syscall.Exec(binary, args, env)
	},
}

func init() {
	runCmd.Flags().StringArrayVarP(&runPaths, "path", "p", nil,
		"Vault secret path(s) to inject (repeatable)")
	_ = runCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(runCmd)
}
