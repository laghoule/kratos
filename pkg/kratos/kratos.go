package kratos

import (
	"fmt"
	"strconv"

	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/pterm/pterm"
)

// Kratos contains info for deployment
type Kratos struct {
	k8s.Client
}

// New return a kratos struct
func New() (*Kratos, error) {
	kclient, err := k8s.New()
	if err != nil {
		return nil, err
	}

	return &Kratos{
		Client: *kclient,
	}, nil
}

// Create the deployment of all objects
func (k *Kratos) Create(name, namespace, image, tag, ingresClass, clusterIssuer string, hostnames []string, replicas, port int32) error {
	runWithError := false

	// deployment
	spinner, _ := pterm.DefaultSpinner.Start("creating deployment ", name)
	if err := k.CreateUpdateDeployment(name, namespace, image, tag, replicas, port); err != nil {
		spinner.Fail(err)
		runWithError = true
	}
	spinner.Success()

	// service
	spinner, _ = pterm.DefaultSpinner.Start("creating service ", name)
	if err := k.CreateUpdateService(name, namespace, port); err != nil {
		spinner.Fail(err)
		runWithError = true
	}
	spinner.Success()

	// ingress
	spinner, _ = pterm.DefaultSpinner.Start("creating ingress ", name)
	if err := k.CreateUpdateIngress(name, namespace, ingresClass, clusterIssuer, hostnames, port); err != nil {
		spinner.Fail(err)
		runWithError = true
	}
	spinner.Success()

	if runWithError {
		return fmt.Errorf("some operations failed")
	}

	return nil
}

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

// Delete all objects of the deployment
func (k *Kratos) Delete(name, namespace string) error {
	runWithError := false

	// ingress
	spinner, _ := pterm.DefaultSpinner.Start("deleting ingress ", name)
	if err := k.DeleteIngress(name, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	// service
	spinner, _ = pterm.DefaultSpinner.Start("deleting service ", name)
	if err := k.DeleteService(name, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	// deployment
	spinner, _ = pterm.DefaultSpinner.Start("deleting deployment ", name)
	if err := k.DeleteDeployment(name, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	if runWithError {
		return fmt.Errorf("some operations failed")
	}

	return nil
}
