package cmd

import (
	"codingstandardsfixer/pkg/csfixing"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "A tool for rolling out coding standards without causing undue merge conflicts",
	Long:  `Fixes coding standards across your whole project, then reverts changes which would have caused a merge conflict with another branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		csfixing.Fix(conf)
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
