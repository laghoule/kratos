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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := k8s.New()
		if err != nil {
			panic(err)
		}

		depList, err := client.ListDeployments(viper.GetString("namespace"))
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
	viper.BindPFlag("namespace", listCmd.PersistentFlags().Lookup("namespace"))
}
