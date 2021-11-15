package kratos

import (
	"fmt"

	"github.com/pterm/pterm"
)

// Delete all objects of the deployment
func (k *Kratos) Delete(name, namespace string) error {
	runWithError := false

	// ingress
	spinner, _ := pterm.DefaultSpinner.Start("deleting ingress ")
	if err := k.DeleteIngress(name, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	// service
	spinner, _ = pterm.DefaultSpinner.Start("deleting service ")
	if err := k.DeleteService(name, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	// deployment
	spinner, _ = pterm.DefaultSpinner.Start("deleting deployment ")
	if err := k.DeleteDeployment(name, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	// configuration
	spinner, _ = pterm.DefaultSpinner.Start("deleting configuration ")
	if err := k.DeleteSecret(name+kratosSuffixConfig, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	if runWithError {
		return fmt.Errorf("some operations failed")
	}

	return nil
}