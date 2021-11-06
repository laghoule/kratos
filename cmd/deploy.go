package cmd

import (
	"log"

	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application in an namespace",
	Long: `Create deployment, service and ingress of the deployed application. 
	Cert-manager will create the necessary TLS certificate.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := k8s.New()
		if err != nil {
			log.Fatal(err)
		}

		name := viper.GetString("depName")
		namespace := viper.GetString("depNamespace")
		image := viper.GetString("depImage")
		tag := viper.GetString("depTag")
		replicas := viper.GetInt32("depReplicas")
		port := viper.GetInt32("depPort")
		ingresClass := viper.GetString("depIngressClass")
		clusterIssuer := viper.GetString("depClusterIssuer")
		hostnames := viper.GetStringSlice("depHostnames")

		// deployment
		spinner, _ := pterm.DefaultSpinner.Start("deploying deployment ", name)
		if err := client.CreateUpdateDeployment(name, namespace, image, tag, replicas, port); err != nil {
			spinner.Fail(err)
			return
		}
		spinner.Success()

		// service
		spinner, _ = pterm.DefaultSpinner.Start("deploying service ", name)
		if err := client.CreateUpdateService(name, namespace, port); err != nil {
			spinner.Fail(err)
			return
		}
		spinner.Success()

		// ingress
		spinner, _ = pterm.DefaultSpinner.Start("deploying ingress ", name)
		if err := client.CreateUpdateIngress(name, namespace, ingresClass, clusterIssuer, hostnames, port); err != nil {
			spinner.Fail(err)
			return
		}
		spinner.Success()
		
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().String("name", "", "deployment name")
	deployCmd.MarkFlagRequired("name")

	deployCmd.Flags().String("namespace", "default", "specify a namespace")

	deployCmd.Flags().String("image", "", "image to deploy")
	deployCmd.MarkFlagRequired("image")

	deployCmd.Flags().String("tag", "latest", "tag to use")

	deployCmd.Flags().Int32("replicas", 1, "number of replicas")

	deployCmd.Flags().Int32("port", 80, "container port")

	deployCmd.Flags().String("ingressClass", "nginx", "ingress class name")
	deployCmd.Flags().String("clusterIssuer", "letsencrypt", "letsencrypt cluster issuer name")
	deployCmd.Flags().StringArray("hostnames", []string{}, "list of hostname to assign")
	deployCmd.MarkFlagRequired("hostnames")

	viper.BindPFlag("depName", deployCmd.Flags().Lookup("name"))
	viper.BindPFlag("depNamespace", deployCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("depImage", deployCmd.Flags().Lookup("image"))
	viper.BindPFlag("depTag", deployCmd.Flags().Lookup("tag"))
	viper.BindPFlag("depReplicas", deployCmd.Flags().Lookup("replicas"))
	viper.BindPFlag("depPort", deployCmd.Flags().Lookup("port"))
	viper.BindPFlag("depIngressClass", deployCmd.Flags().Lookup("ingressClass"))
	viper.BindPFlag("depClusterIssuer", deployCmd.Flags().Lookup("clusterIssuer"))
	viper.BindPFlag("depHostnames", deployCmd.Flags().Lookup("hostnames"))
}
