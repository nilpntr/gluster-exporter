package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var appVersion string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gluster Exporter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Gluster Exporter version: %s", appVersion))
	},
}
