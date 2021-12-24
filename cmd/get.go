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

		fmt.Printf("configuration saved at %s/%s\n", destination, name+kratos.YamlExt)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().String("name", "", "name of the deployment")
	if err := getCmd.MarkFlagRequired("name"); err != nil {
		errorExit(err.Error())
	}

	getCmd.Flags().StringP("namespace", "n", "", "namespace of the deployment")
	if err := getCmd.MarkFlagRequired("namespace"); err != nil {
		errorExit(err.Error())
	}

	getCmd.Flags().StringP("destination", "d", ".", "destination path")

	if err := viper.BindPFlag("gName", getCmd.Flags().Lookup("name")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("gNamespace", getCmd.Flags().Lookup("namespace")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("gdestination", getCmd.Flags().Lookup("destination")); err != nil {
		errorExit(err.Error())
	}
}
