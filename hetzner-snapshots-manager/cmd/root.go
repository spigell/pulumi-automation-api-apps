package cmd

import (
	"log"
	"fmt"
	"github.com/spigell/pulumi-automation-api-apps/common/version"


	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "hetzner-snapshot-manager",
		Short: "manage snapshots based on pulumi preview events",
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Get())
		},
	}
)




func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().String("hcloud-token", "", "Hetzner Cloud token")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().BoolP("diff", "d", false, "Enable the diff option for pulumi command")
	rootCmd.PersistentFlags().Bool("only-api-server", false, "Run only api server and do not stop it. For testing purposes.")
	rootCmd.PersistentFlags().Int("api-server-port", 0, "default is random")
	rootCmd.PersistentFlags().Int("cleaner-max-keep", 1, "default is keepling only the last snaphot")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("diff", rootCmd.PersistentFlags().Lookup("diff"))
	viper.BindPFlag("hcloud-token", rootCmd.PersistentFlags().Lookup("hcloud-token"))
	viper.BindEnv("hcloud_token")
	viper.BindPFlag("api-server-port", rootCmd.PersistentFlags().Lookup("api-server-port"))
	viper.BindPFlag("max-keep", rootCmd.PersistentFlags().Lookup("cleaner-max-keep"))

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
