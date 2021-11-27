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
	Short: "Retreive a configuration of a kratos deployment.",
	Long:  `Download the saved configuration of a kratos deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("gName")
		namespace := viper.GetString("gNamespace")
		destination := viper.GetString("gdestination")

		kratos, err := kratos.New("", viper.GetString("kubeconfig"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := kratos.SaveConfigToDisk(name, namespace, destination); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().String("name", "", "name of the deployment")
	getCmd.MarkFlagRequired("name")

	getCmd.Flags().StringP("namespace", "n", "", "namespace of the deployment")
	getCmd.MarkFlagRequired("namespace")

	getCmd.Flags().StringP("destination", "d", ".", "destination path")

	viper.BindPFlag("gName", getCmd.Flags().Lookup("name"))
	viper.BindPFlag("gNamespace", getCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("gdestination", getCmd.Flags().Lookup("destination"))
}