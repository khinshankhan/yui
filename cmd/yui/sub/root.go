package sub

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yui",
	Short: "A bespoke kit of tools — each built to solve one small problem very well.",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(caseCmd)
}
