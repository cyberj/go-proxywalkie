package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "X.X.X"

func init() {
	rootCmd.AddCommand(versionCmd)
}

func SetVersion(version string) {
	Version = version
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of proxywalkie",
	// Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
		os.Exit(0)
	},
}
