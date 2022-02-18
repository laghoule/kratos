package kratos

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/pterm/pterm"
)

// Delete all objects of the deployment
func (k *Kratos) Delete(name, namespace string) error {
	runWithError := false

	if err := k.loadConfigFromSecret(name, namespace); err != nil {
		return err
	}

	// deployment
	if k.Config.Deployment != nil {
		// ingress
		spinner, _ := pterm.DefaultSpinner.Start("deleting ingress ")
		if err := k.Client.Ingress.Delete(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()

		// service
		spinner, _ = pterm.DefaultSpinner.Start("deleting service ")
		if err := k.Client.Service.Delete(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()

		// deployment
		spinner, _ = pterm.DefaultSpinner.Start("deleting deployment ")
		if err := k.Client.Deployment.Delete(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()
	}

	// cronjob
	if k.Config.Cronjob != nil {
		spinner, _ := pterm.DefaultSpinner.Start("deleting cronjob ")
		if err := k.Client.Cronjob.Delete(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()
	}

	// configmaps
	if k.Config.ConfigMaps != nil {
		spinner, _ := pterm.DefaultSpinner.Start("deleting configmaps ")
		if err := k.Client.ConfigMaps.Delete(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()
	}

	// secrets
	if k.Config.Secrets != nil {
		spinner, _ := pterm.DefaultSpinner.Start("deleting secrets ")
		if err := k.Client.Secrets.Delete(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()
	}

	// configuration
	spinner, _ := pterm.DefaultSpinner.Start("deleting configuration ")
	if err := k.DeleteConfig(name+config.ConfigSuffix, namespace); err != nil {
		pterm.Error.Println(err)
		runWithError = true
	}
	spinner.Success()

	if runWithError {
		return fmt.Errorf("some operations failed")
	}

	return nil
}
