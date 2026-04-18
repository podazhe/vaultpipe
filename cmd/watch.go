package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var watchCmd = &cobra.Command{
	Ulient, err := vault.NewClient(viper.GetString("vault.address"), viper.GetString("vault.token"))
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		interval := viper.GetDuration("watch.interval")
		if interval == 0 {
			interval = 30 * time.Second
		}

		watcher := vault.NewWatcher(client, args, interval)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		watcher.Start(ctx)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case ev := <-watcher.Events():
				if ev.Err != nil {
					fmt.Fprintf(os.Stderr, "[watch] error %s: %v\n", ev.Path, ev.Err)
					continue
				}
				fmt.Printf("[watch] change detected: %s\n", ev.Path)
				for k, v := range ev.Data {
					fmt.Printf("  %s=%s\n", k, v)
				}
			case <-sig:
				watcher.Stop()
				return nil
			}
		}
	},
}

func init() {
	watchCmd.Flags().Duration("interval", 30*time.Second, "polling interval")
	_ = viper.BindPFlag("watch.interval", watchCmd.Flags().Lookup("interval"))
	rootCmd.AddCommand(watchCmd)
}
