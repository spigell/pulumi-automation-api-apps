package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spigell/pulumi-automation-api-apps/hetzner-snapshots-manager/manager"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "The 'up' command runs pulumi",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		m, err := manager.New(ctx, "up")
		if err != nil {
			log.Fatalf("create manager: %s", err)
		}

		if err := m.Run(cmd); err != nil {
			m.Logger.Fatal(fmt.Sprintf("run: %s", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
