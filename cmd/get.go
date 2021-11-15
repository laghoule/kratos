package cmd

import (
	"fmt"
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get retreive a configuration of a kratos deployment.",
	Long: `Get download the saved configuration of a kratos deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("gName")
		namespace := viper.GetString("gNamespace")
		destination := viper.GetString("gdestination")

		kratos, err := kratos.New("")
		if err != nil {
			panic(err)
		}

		if err := kratos.SaveConfigFileToDisk(name, namespace, destination); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().String("name", "", "name of the deployment")
	getCmd.MarkFlagRequired("name")

	getCmd.PersistentFlags().String("namespace", "", "namespace of the deployment")
	getCmd.MarkFlagRequired("namespace")

	getCmd.PersistentFlags().String("destination", ".", "destination path")

	viper.BindPFlag("gName", getCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("gNamespace", getCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("gdestination", getCmd.PersistentFlags().Lookup("destination"))
}
