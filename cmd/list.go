package cmd

import (
	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications.",
	Long:  `List applications in namespace, or cluster wide.`,
	Run: func(cmd *cobra.Command, args []string) {
		k, err := kratos.New("", viper.GetString("kubeconfig"))
		if err != nil {
			errorExit(err.Error())
		}

		if err := k.PrintList(viper.GetString("listNamespace")); err != nil {
			errorExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("namespace", "n", "", "specify a namespace")
	viper.BindPFlag("listNamespace", listCmd.Flags().Lookup("namespace"))
}
