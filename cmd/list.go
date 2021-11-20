package cmd

import (
	"fmt"
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List application of managed kratos deployment.",
	Long:  `List application by namespace, or cluster wide of managed kratos deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		kratos, err := kratos.New("", viper.GetString("kubeconfig"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := kratos.List(viper.GetString("listNamespace")); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("namespace", "n", "", "specify a namespace")
	viper.BindPFlag("listNamespace", listCmd.Flags().Lookup("namespace"))
}
