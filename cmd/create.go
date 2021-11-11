package cmd

import (
	"fmt"
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploy an application in an namespace",
	Long: `Create deployment, service and ingress for k8s application. 
Cert-manager will create the necessary TLS certificate.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		config := viper.GetString("cConfig")
		name := viper.GetString("cName")
		namespace := viper.GetString("cNamespace")

		kratos, err := kratos.New(config)
		if err != nil {
			panic(err)
		}

		if err := kratos.Create(name, namespace); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("config", "", "configuration file")
	createCmd.MarkFlagRequired("config")

	createCmd.Flags().String("name", "", "deployment name")
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().String("namespace", "default", "deployment namespace")
	createCmd.MarkFlagRequired("namespace")

	viper.BindPFlag("cConfig", createCmd.Flags().Lookup("config"))
	viper.BindPFlag("cName", createCmd.Flags().Lookup("name"))
	viper.BindPFlag("cNamespace", createCmd.Flags().Lookup("namespace"))
}
