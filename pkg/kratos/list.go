package kratos

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
)

// List deployment
func (k *Kratos) List(namespace string) error {
	depList, err := k.ListDeployments(namespace)
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
	return nil
}
