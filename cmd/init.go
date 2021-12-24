package cmd

import (
	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the list command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty configuration file.",
	Long:  `Initialise a configuration file with all available options.`,
	Run: func(cmd *cobra.Command, args []string) {
		k, err := kratos.New("", viper.GetString("kubeconfig"))
		if err != nil {
			errorExit(err.Error())
		}

		if err := k.CreateInit(viper.GetString("initName")); err != nil {
			errorExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	
	initCmd.Flags().String("name", "", "name of the configuration file")
	if err := initCmd.MarkFlagRequired("name"); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("initName", initCmd.Flags().Lookup("name")); err != nil {
		errorExit(err.Error())
	}
}
