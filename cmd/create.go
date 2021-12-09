package cmd

import (
	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploy an application.",
	Long:  `Create and update an application in a namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := viper.GetString("cConfig")
		name := viper.GetString("cName")
		namespace := viper.GetString("cNamespace")

		k, err := kratos.New(config, viper.GetString("kubeconfig"))
		if err != nil {
			errorExit(err.Error())
		}

		if err := k.IsDependencyMeet(); err != nil {
			errorExit(err.Error())
		}

		if err := k.Create(name, namespace); err != nil {
			errorExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("config", "c", "", "configuration file")
	createCmd.MarkFlagRequired("config")

	createCmd.Flags().String("name", "", "deployment name")
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().StringP("namespace", "n", "default", "deployment namespace")
	createCmd.MarkFlagRequired("namespace")

	viper.BindPFlag("cConfig", createCmd.Flags().Lookup("config"))
	viper.BindPFlag("cName", createCmd.Flags().Lookup("name"))
	viper.BindPFlag("cNamespace", createCmd.Flags().Lookup("namespace"))
}
