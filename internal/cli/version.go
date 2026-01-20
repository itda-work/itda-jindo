package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version and BuildDate are set via ldflags at build time
var (
	Version   = "dev"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of jindo.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("jindo version %s (built %s)\n", Version, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
