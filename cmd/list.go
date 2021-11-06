package cmd

import (
	"fmt"
	"strconv"

	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/pterm/pterm"
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
		client, err := k8s.New()
		if err != nil {
			panic(err)
		}

		depList, err := client.ListDeployments(viper.GetString("listNamespace"))
		if err != nil {
			pterm.Error.Println(err)
		}

		pdata := pterm.TableData{
			{"Name", "Namespace", "Replicas", "Creation", "Revision"},
		}

		for _, item := range depList {
			pdata = append(pdata, []string{
				item.Name,
				item.Namespace,
				strconv.Itoa(int(*item.Spec.Replicas)),
				item.CreationTimestamp.UTC().String(),
				fmt.Sprint(item.Generation),
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(pdata).Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().String("namespace", "", "specify a namespace")
	viper.BindPFlag("listNamespace", listCmd.PersistentFlags().Lookup("namespace"))
}
