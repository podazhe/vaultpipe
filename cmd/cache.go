package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the vaultpipe secret cache",
}

var cacheFlushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush all cached secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalCache == nil {
			fmt.Println("Cache is not initialised.")
			return nil
		}
		globalCache.Flush()
		fmt.Println("Secret cache flushed.")
		return nil
	},
}

var cacheInvalidateCmd = &cobra.Command{
	Use:   "invalidate <path>",
	Short: "Invalidate a single cached secret path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalCache == nil {
			fmt.Println("Cache is not initialised.")
			return nil
		}
		globalCache.Invalidate(args[0])
		fmt.Printf("Invalidated cache entry for path: %s\n", args[0])
		return nil
	},
}

func init() {
	cacheCmd.AddCommand(cacheFlushCmd)
	cacheCmd.AddCommand(cacheInvalidateCmd)
	RootCmd.AddCommand(cacheCmd)
}
