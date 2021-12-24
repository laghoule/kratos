package cmd

import (
	"fmt"
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an application.",
	Long:  `Delete an application in a namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("dName")
		namespace := viper.GetString("dNamespace")

		k, err := kratos.New("", viper.GetString("kubeconfig"))
		if err != nil {
			errorExit(err.Error())
		}

		if found, err := k.IsReleaseExist(name, namespace); !found && err == nil {
			fmt.Printf("%s don't exist\n", name)
			os.Exit(1)
		} else {
			if err != nil {
				errorExit(err.Error())
			}
		}

		if err := k.Delete(name, namespace); err != nil {
			errorExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().String("name", "", "name of the deployment")
	if err := deleteCmd.MarkFlagRequired("name"); err != nil {
		errorExit(err.Error())
	}

	deleteCmd.Flags().StringP("namespace", "n", "", "namespace of the deployment")
	if err := deleteCmd.MarkFlagRequired("namespace"); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("dName", deleteCmd.Flags().Lookup("name")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("dNamespace", deleteCmd.Flags().Lookup("namespace")); err != nil {
		errorExit(err.Error())
	}
}
