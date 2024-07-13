package cmd

import (
	"codingstandardsfixer/pkg/csfixing"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var confPath string
var conf csfixing.ApplicationConfig

var rootCmd = &cobra.Command{
	Use:   "coding-standards-rollout-tool",
	Short: "A tool for rolling out coding standards without causing undue merge conflicts",
	Long:  `Fixes coding standards across your whole project, then reverts changes which would have caused a merge conflict with another branch.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", "", "config file (required)")
}

func initConfig() {
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		os.Stderr.WriteString("Config file not found\n")
		os.Exit(1)
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(confPath)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	conf = csfixing.BuildConfig(viper.AllSettings())
}
