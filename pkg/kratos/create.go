package kratos

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/pterm/pterm"
)

// Create the deployment of all objects
func (k *Kratos) Create(name, namespace string) error {
	runWithError := false

	// configuration is saved first
	spinner, _ := pterm.DefaultSpinner.Start("saving configuration")
	if err := k.saveConfigToSecret(name+config.ConfigSuffix, namespace); err != nil {
		spinner.Fail(err)
		runWithError = true
	} else {
		spinner.Success()
	}

	// deployment
	if k.Config.Deployment != nil {
		spinner, _ := pterm.DefaultSpinner.Start("creating deployment")
		if err := k.Client.Deployment.CreateUpdate(name, namespace); err != nil {
			spinner.Fail(err)
			runWithError = true
		} else {
			spinner.Success()
		}

		// service
		spinner, _ = pterm.DefaultSpinner.Start("creating service")
		if err := k.Client.Service.CreateUpdate(name, namespace); err != nil {
			spinner.Fail(err)
			runWithError = true
		} else {
			spinner.Success()
		}

		// ingress
		spinner, _ = pterm.DefaultSpinner.Start("creating ingress")
		if err := k.Client.Ingress.CreateUpdate(name, namespace); err != nil {
			spinner.Fail(err)
			runWithError = true
		} else {
			spinner.Success()
		}
	}

	// cronjob
	if k.Config.Cronjob != nil {
		spinner, _ := pterm.DefaultSpinner.Start("creating cronjob")
		if err := k.Client.Cronjob.CreateUpdate(name, namespace); err != nil {
			spinner.Fail(err)
			runWithError = true
		} else {
			spinner.Success()
		}
	}

	// configmaps
	if k.Config.ConfigMaps != nil {
		spinner, _ := pterm.DefaultSpinner.Start("creating configmaps")
		if err := k.Client.ConfigMaps.CreateUpdate(name, namespace); err != nil {
			spinner.Fail(err)
			runWithError = true
		} else {
			spinner.Success()
		}
	}

	// secrets
	if k.Config.Secrets != nil {
		spinner, _ := pterm.DefaultSpinner.Start("creating secrets")
		if err := k.Client.Secret.CreateUpdate(name, namespace); err != nil {
			spinner.Fail(err)
			runWithError = true
		} else {
			spinner.Success()
		}
	}

	if runWithError {
		return fmt.Errorf("some operations failed")
	}

	return nil
}
