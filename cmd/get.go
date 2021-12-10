package cmd

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a configuration of deployed application.",
	Long:  `Download the saved configuration of a deployed application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("gName")
		namespace := viper.GetString("gNamespace")
		destination := viper.GetString("gdestination")

		k, err := kratos.New("", viper.GetString("kubeconfig"))
		if err != nil {
			errorExit(err.Error())
		}

		if found, err := k.IsReleaseExist(name, namespace); !found && err == nil {
			errorExit(fmt.Sprintf("%s don't exist\n", name))
		} else {
			if err != nil {
				errorExit(err.Error())
			}
		}

		if err := k.SaveConfigToDisk(name, namespace, destination); err != nil {
			errorExit(err.Error())
		}

		fmt.Printf("configuration saved at %s/%s", destination, name+kratos.YamlExt)
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
