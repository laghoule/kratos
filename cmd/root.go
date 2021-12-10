package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

var rootCmd = &cobra.Command{
	Use:               "kratos",
	Short:             "Easy deployment tool for Kubernetes.",
	Long:              `Alternative to helm for deploying simple container, without the pain of managing Kubernetes yaml templates.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	TraverseChildren:  true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	kubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	rootCmd.PersistentFlags().StringP("kubeconfig", "k", kubeconfig, "kubernetes configuration file")
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
}

func initConfig() {
	viper.AutomaticEnv()
}

func errorExit(err string) {
	fmt.Println(err)
	os.Exit(1)
}
