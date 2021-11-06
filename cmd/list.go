package cmd

import (
	"os"

	"github.com/laghoule/kratos/pkg/kratos"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List application of managed kratos deployment",
	Long: `List application by namespace, or cluster wide of managed
kratos deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		kratos, err := kratos.New()
		if err != nil {
			panic(err)
		}

		if err := kratos.List(viper.GetString("listNamespace")); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().String("namespace", "", "specify a namespace")
	viper.BindPFlag("listNamespace", listCmd.PersistentFlags().Lookup("namespace"))
}
