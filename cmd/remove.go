package cmd

import (
	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a deployment in a namespace",
	Long: `Remove the deployment, service and ingress of the deployed application. 
Generated cert-manager secret will not be deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := k8s.New()
		if err != nil {
			panic(err)
		}

		name := viper.GetString("delName")
		namespace := viper.GetString("delNamespace")

		// ingress
		spinner, _ := pterm.DefaultSpinner.Start("removing ingress ", name)
		if err := client.DeleteIngress(name, namespace); err != nil {
			pterm.Error.Println(err)
		}
		spinner.Success()

		// service
		spinner, _ = pterm.DefaultSpinner.Start("removing service ", name)
		if err := client.DeleteService(name, namespace); err != nil {
			pterm.Error.Println(err)
		}
		spinner.Success()

		// deployment
		spinner, _ = pterm.DefaultSpinner.Start("removing deployment ", name)
		if err := client.DeleteDeployment(name, namespace); err != nil {
			pterm.Error.Println(err)
		}
		spinner.Success()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.PersistentFlags().String("name", "", "name of the deployment")
	removeCmd.MarkFlagRequired("name")
	removeCmd.PersistentFlags().String("namespace", "", "namespace of the deployment")
	removeCmd.MarkFlagRequired("namespace")
	viper.BindPFlag("delName", removeCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("delNamespace", removeCmd.PersistentFlags().Lookup("namespace"))
}
