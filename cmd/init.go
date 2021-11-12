package cmd

import (
	"fmt"
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the list command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty configuration file.",
	Long: `Initialise a configuration file with all available options.`,
	Run: func(cmd *cobra.Command, args []string) {
		kratos, err := kratos.New("")
		if err != nil {
			panic(err)
		}

		if err := kratos.CreateInit(viper.GetString("initName")); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().String("name", "", "name of the configuration file")
	initCmd.MarkFlagRequired("name")
	viper.BindPFlag("initName", initCmd.PersistentFlags().Lookup("name"))
}