package cmd

import (
	"github.com/laghoule/kratos/pkg/kratos"
	"github.com/pterm/pterm"

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

		if err := k.CheckDependency(); err != nil {
			errorExit(err.Error())
		}

		if exist, err := k.IsReleaseExist(name, namespace); exist && err == nil {
			pterm.Warning.Printfln("existing release found for %s/%s, will do an update", namespace, name)
		} else {
			if err != nil {
				errorExit(err.Error())
			}
		}

		if err := k.Create(name, namespace); err != nil {
			errorExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("config", "c", "", "configuration file")
	if err := createCmd.MarkFlagRequired("config"); err != nil {
		errorExit(err.Error())
	}

	createCmd.Flags().String("name", "", "deployment name")
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		errorExit(err.Error())
	}

	createCmd.Flags().StringP("namespace", "n", "default", "deployment namespace")
	if err := createCmd.MarkFlagRequired("namespace"); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("cConfig", createCmd.Flags().Lookup("config")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("cName", createCmd.Flags().Lookup("name")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("cNamespace", createCmd.Flags().Lookup("namespace")); err != nil {
		errorExit(err.Error())
	}
}
