package cmd

import (
	"log"
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploy an application in an namespace",
	Long: `Create deployment, service and ingress of the deployed application. 
	Cert-manager will create the necessary TLS certificate.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		kratos, err := kratos.New()
		if err != nil {
			log.Fatal(err)
		}

		name := viper.GetString("cName")
		namespace := viper.GetString("cNamespace")
		image := viper.GetString("cImage")
		tag := viper.GetString("cTag")
		replicas := viper.GetInt32("cReplicas")
		port := viper.GetInt32("cPort")
		ingresClass := viper.GetString("cIngressClass")
		clusterIssuer := viper.GetString("cClusterIssuer")
		hostnames := viper.GetStringSlice("cHostnames")

		if err := kratos.Create(name, namespace, image, tag, ingresClass, clusterIssuer, hostnames, replicas, port); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("name", "", "createment name")
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().String("namespace", "default", "specify a namespace")

	createCmd.Flags().String("image", "", "image to create")
	createCmd.MarkFlagRequired("image")

	createCmd.Flags().String("tag", "latest", "tag to use")

	createCmd.Flags().Int32("replicas", 1, "number of replicas")

	createCmd.Flags().Int32("port", 80, "container port")

	createCmd.Flags().String("ingressClass", "nginx", "ingress class name")
	createCmd.Flags().String("clusterIssuer", "letsencrypt", "letsencrypt cluster issuer name")
	createCmd.Flags().StringArray("hostnames", []string{}, "list of hostname to assign")
	createCmd.MarkFlagRequired("hostnames")

	viper.BindPFlag("cName", createCmd.Flags().Lookup("name"))
	viper.BindPFlag("cNamespace", createCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("cImage", createCmd.Flags().Lookup("image"))
	viper.BindPFlag("cTag", createCmd.Flags().Lookup("tag"))
	viper.BindPFlag("cReplicas", createCmd.Flags().Lookup("replicas"))
	viper.BindPFlag("cPort", createCmd.Flags().Lookup("port"))
	viper.BindPFlag("cIngressClass", createCmd.Flags().Lookup("ingressClass"))
	viper.BindPFlag("cClusterIssuer", createCmd.Flags().Lookup("clusterIssuer"))
	viper.BindPFlag("cHostnames", createCmd.Flags().Lookup("hostnames"))
}
