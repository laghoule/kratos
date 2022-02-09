package cmd

import (
	"github.com/laghoule/kratos/pkg/kratos"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the create command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an application.",
	Long:  `Update an existing application in a namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := viper.GetString("uConfig")
		name := viper.GetString("uName")
		namespace := viper.GetString("uNamespace")

		k, err := kratos.New(config, viper.GetString("kubeconfig"))
		if err != nil {
			errorExit(err.Error())
		}

		if err := k.CheckDependency(); err != nil {
			errorExit(err.Error())
		}

		if exist, err := k.IsReleaseExist(name, namespace); exist && err == nil {
			if err := k.Create(name, namespace); err != nil {
				errorExit(err.Error())
			}
		} else {
			if err != nil {
				errorExit(err.Error())
			}
			pterm.Error.Printf("application %s don't exist in namespace %s\n", name, namespace)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("config", "c", "", "configuration file")
	if err := updateCmd.MarkFlagRequired("config"); err != nil {
		errorExit(err.Error())
	}

	updateCmd.Flags().String("name", "", "deployment name")
	if err := updateCmd.MarkFlagRequired("name"); err != nil {
		errorExit(err.Error())
	}

	updateCmd.Flags().StringP("namespace", "n", "default", "deployment namespace")
	if err := updateCmd.MarkFlagRequired("namespace"); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("uConfig", updateCmd.Flags().Lookup("config")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("uName", updateCmd.Flags().Lookup("name")); err != nil {
		errorExit(err.Error())
	}

	if err := viper.BindPFlag("uNamespace", updateCmd.Flags().Lookup("namespace")); err != nil {
		errorExit(err.Error())
	}
}
