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
			fmt.Println(err)
			os.Exit(1)
		}

		if found, err := k.IsReleaseExist(name, namespace); !found && err == nil {
			fmt.Printf("%s don't exist\n", name)
			os.Exit(1)
		} else {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if err := k.Delete(name, namespace); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().String("name", "", "name of the deployment")
	deleteCmd.MarkFlagRequired("name")

	deleteCmd.Flags().StringP("namespace", "n", "", "namespace of the deployment")
	deleteCmd.MarkFlagRequired("namespace")

	viper.BindPFlag("dName", deleteCmd.Flags().Lookup("name"))
	viper.BindPFlag("dNamespace", deleteCmd.Flags().Lookup("namespace"))
}
