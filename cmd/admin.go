package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/cmd/devtools"
)

func init() {
	adminCmd.AddCommand(devtools.SetMainTagsCmd)
}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "admin utils",
	Long:  "utils for administrator",
}
