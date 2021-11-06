package cmd

import (
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a deployment in a namespace",
	Long: `Delete the deployment, service and ingress of the deployed application. 
Generated cert-manager secret will not be deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		kratos, err := kratos.New()
		if err != nil {
			panic(err)
		}

		name := viper.GetString("dName")
		namespace := viper.GetString("dNamespace")

		if err := kratos.Delete(name, namespace); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().String("name", "", "name of the deployment")
	deleteCmd.MarkFlagRequired("name")

	deleteCmd.PersistentFlags().String("namespace", "", "namespace of the deployment")
	deleteCmd.MarkFlagRequired("namespace")

	viper.BindPFlag("dName", deleteCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("dNamespace", deleteCmd.PersistentFlags().Lookup("namespace"))
}
