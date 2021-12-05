package kratos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/laghoule/kratos/pkg/config"

	"gopkg.in/yaml.v3"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configSuffix = "-kratos-config"
	configKey    = "config"
	fileMode     = 0666
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

	secret := k.createConfigSecret(name, namespace, string(b))

	if err := k.Client.CreateUpdateSecret(secret, namespace); err != nil {
		return err
	}

	return nil
}

// createConfigSecret return the configuration as a secret object
func (k *Kratos) createConfigSecret(name, namespace, data string) *corev1.Secret {
	if k.Config.Common == nil {
		k.Config.Common = &config.Common{} // FIXME pass common labels & annotations
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      k.Config.Common.Labels,
			Annotations: k.Config.Common.Annotations,
		},
		StringData: map[string]string{
			configKey: data,
		},
	}
}

// SaveConfigToDisk get config from secret and write it to disk
func (k *Kratos) SaveConfigToDisk(name, namespace, destination string) error {
	secret, err := k.Client.GetSecret(name+configSuffix, namespace)
	if err != nil {
		return err
	}

	if _, ok := secret.Data[configKey]; ok {
		if err := os.WriteFile(filepath.Join(destination, name)+YamlExt, []byte(secret.Data[configKey]), fileMode); err != nil {
			return fmt.Errorf("writing yaml init file failed: %s", err)
		}
	} else {
		if _, ok := secret.StringData[configKey]; ok {
			if err := os.WriteFile(filepath.Join(destination, name)+YamlExt, []byte(secret.StringData[configKey]), fileMode); err != nil {
				return fmt.Errorf("writing yaml init file failed: %s", err)
			}
		} else {
			return fmt.Errorf("unexpected missing data in secret %s", secret.Name)
		}
	}

	return nil
}
