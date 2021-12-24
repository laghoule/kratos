package cmd

import (
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// build informations
var (
	version   = "devel"
	buildDate = time.Now().String()
	gitCommit = ""
	gitRef    = ""
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of kratos.",
	Long:  `Show version, build time, git commit, git reference informations.`,
	Run: func(cmd *cobra.Command, args []string) {
		pdata := pterm.TableData{
			{"Version", "Build date", "Git commit", "Git reference"},
			{version, buildDate, gitCommit, gitRef},
		}
		if err := pterm.DefaultTable.WithHasHeader().WithData(pdata).Render(); err != nil {
			errorExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
