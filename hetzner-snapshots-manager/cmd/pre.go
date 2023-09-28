package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spigell/pulumi-automation-api-apps/hetzner-snapshots-manager/manager"

	"github.com/spf13/cobra"
)

var preCmd = &cobra.Command{
	Use:   "pre",
	Short: "The 'pre' only make preview",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		ctx := context.Background()
		m, err := manager.New(ctx, "preview")
		if err != nil {
			log.Fatalf("create manager: %s", err)
		}

		if err := m.Run(cmd); err != nil {
			m.Logger.Fatal(fmt.Sprintf("run: %s", err))
		}
	},
}

func init() {
	preCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "/path/to/file for config. Required")
	preCmd.MarkPersistentFlagRequired("config")
	rootCmd.AddCommand(preCmd)
}
