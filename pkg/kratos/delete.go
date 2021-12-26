package kratos

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/pterm/pterm"
)

// Delete all objects of the deployment
func (k *Kratos) Delete(name, namespace string) error {
	runWithError := false

	secret, err := k.GetSecret(name+config.ConfigSuffix, namespace)
	if err != nil {
		return fmt.Errorf("getting config from secret failed: %s", err)
	}

	conf := ""
	if _, ok := secret.Data[config.ConfigKey]; ok {
		conf = string(secret.Data[config.ConfigKey])
	} else {
		if _, ok := secret.StringData[config.ConfigKey]; ok {
			conf = secret.StringData[config.ConfigKey]
		} else {
			return fmt.Errorf("getting config data from secret failed")
		}
	}

	if err := k.Config.LoadFromString(conf); err != nil {
		return err
	}

	// deployment
	if k.Config.Deployment != nil {
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
	}

	// cronjob
	if k.Config.Cronjob != nil {
		spinner, _ := pterm.DefaultSpinner.Start("deleting cronjob ")
		if err := k.DeleteCronjob(name, namespace); err != nil {
			pterm.Error.Println(err)
			runWithError = true
		}
		spinner.Success()
	}

	// secrets
	if k.Config.Secrets != nil {
		spinner, _ := pterm.DefaultSpinner.Start("deleting secrets ")
		if err := k.DeleteSecrets(name, namespace, k.Config); err != nil {
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
