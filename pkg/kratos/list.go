package kratos

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
)

// IsReleaseExist check if a release of an app exist in the list of release
func (k *Kratos) IsReleaseExist(name, namespace string) (bool, error) {
	list, err := k.GetList(namespace)
	if err != nil {
		return false, fmt.Errorf("getting list of release failed: %s", err)
	}

	for _, release := range list {
		if release == name {
			return true, nil
		}
	}

	return false, nil
}

// GetList return a list of kratos deployed release name
func (k *Kratos) GetList(namespace string) ([]string, error) {
	list := []string{}

	// deployments
	depList, err := k.ListDeployments(namespace)
	if err != nil {
		return []string{}, err
	}

	for _, deployment := range depList {
		list = append(list, deployment.Name)
	}

	// cronjobs
	cronList, err := k.ListCronjobs(namespace)
	if err != nil {
		return []string{}, err
	}

	for _, cronjob := range cronList {
		list = append(list, cronjob.Name)
	}

	return list, nil
}

// PrintList list deployment to stdout
func (k *Kratos) PrintList(namespace string) error {
	depList, err := k.ListDeployments(namespace)
	if err != nil {
		pterm.Error.Println(err)
	}

	cronList, err := k.ListCronjobs(namespace)
	if err != nil {
		pterm.Error.Println(err)
	}

	pdata := pterm.TableData{
		{"Name", "Namespace", "Type", "Replicas", "Creation", "Revision"},
	}

	// deployments
	for _, item := range depList {
		pdata = append(pdata, []string{
			item.Name,
			item.Namespace,
			"deployment",
			strconv.Itoa(int(*item.Spec.Replicas)),
			item.CreationTimestamp.UTC().String(),
			fmt.Sprint(item.Generation),
		})
	}

	// cronjobs
	for _, item := range cronList {
		pdata = append(pdata, []string{
			item.Name,
			item.Namespace,
			"cronjob",
			"-",
			item.CreationTimestamp.UTC().String(),
			fmt.Sprint(item.Generation),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(pdata).Render()
	return nil
}
