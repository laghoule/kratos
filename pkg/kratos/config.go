package kratos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/laghoule/kratos/pkg/config"

	"gopkg.in/yaml.v3"
)

const (
	fileMode = 0666
	// YamlExt is the default yaml extension
	YamlExt = ".yaml"
)

// CreateInit create sample configuration file
func (k *Kratos) CreateInit(file string) error {
	b, err := yaml.Marshal(config.CreateInit())
	if err != nil {
		return fmt.Errorf("marshaling yaml failed: %s", err)
	}

	if err := os.WriteFile(file, b, fileMode); err != nil {
		return fmt.Errorf("writing yaml init file failed: %s", err)
	}

	return nil
}

// saveConfigToSecret save configuration in secret DataString
func (k *Kratos) saveConfigToSecret(name, namespace string) error {
	b, err := yaml.Marshal(k.Config)
	if err != nil {
		return fmt.Errorf("saving configuration to kubernetes secret failed: %s", err)
	}

	if err := k.Client.SaveConfig(name, namespace, config.ConfigKey, string(b)); err != nil {
		return err
	}

	return nil
}

// SaveConfigToDisk get config from secret and write it to disk
func (k *Kratos) SaveConfigToDisk(name, namespace, destination string) error {
	secret, err := k.Client.Get(name+config.ConfigSuffix, namespace)
	if err != nil {
		return err
	}

	if _, ok := secret.Data[config.ConfigKey]; ok {
		if err := os.WriteFile(filepath.Join(destination, name)+YamlExt, []byte(secret.Data[config.ConfigKey]), fileMode); err != nil {
			return fmt.Errorf("writing yaml init file failed: %s", err)
		}
	} else {
		if _, ok := secret.StringData[config.ConfigKey]; ok {
			if err := os.WriteFile(filepath.Join(destination, name)+YamlExt, []byte(secret.StringData[config.ConfigKey]), fileMode); err != nil {
				return fmt.Errorf("writing yaml init file failed: %s", err)
			}
		} else {
			return fmt.Errorf("unexpected missing data in secret %s", secret.Name)
		}
	}

	return nil
}
