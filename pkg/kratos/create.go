package kratos

import (
	"fmt"

	"github.com/pterm/pterm"
)

// Create the deployment of all objects
func (k *Kratos) Create(name, namespace string) error {
	runWithError := false

	// deployment
	spinner, _ := pterm.DefaultSpinner.Start("creating deployment")
	if err := k.CreateUpdateDeployment(name, namespace, k.Config); err != nil {
		spinner.Fail(err)
		runWithError = true
	} else {
		spinner.Success()
	}

	// service
	spinner, _ = pterm.DefaultSpinner.Start("creating service")
	if err := k.CreateUpdateService(name, namespace, k.Config); err != nil {
		spinner.Fail(err)
		runWithError = true
	} else {
		spinner.Success()
	}

	// ingress
	spinner, _ = pterm.DefaultSpinner.Start("creating ingress")
	if err := k.CreateUpdateIngress(name, namespace, k.Config); err != nil {
		spinner.Fail(err)
		runWithError = true
	} else {
		spinner.Success()
	}

	// configuration
	spinner, _ = pterm.DefaultSpinner.Start("saving configuration")
	if err := k.saveConfigToSecret(name+configSuffix, namespace); err != nil {
		spinner.Fail(err)
		runWithError = true
	} else {
		spinner.Success()
	}

	if runWithError {
		return fmt.Errorf("some operations failed")
	}

	return nil
}
