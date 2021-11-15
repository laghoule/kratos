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
	kratosSuffixConfig = "-kratos-config"
	secretConfigKey    = "config"
	fileMode           = 0666
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

func (k *Kratos) saveConfigFileToSecret(name, namespace string) error {
	b, err := yaml.Marshal(k.Config)
	if err != nil {
		return fmt.Errorf("saving configuration to kubernetes secret failed: %s", err)
	}

	secret := createSecretString(name, namespace, string(b))

	if err := k.Client.CreateUpdateSecret(secret, namespace); err != nil {
		return err
	}

	return nil
}

func createSecretString(name, namespace, data string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: map[string]string{
			secretConfigKey: data,
		},
	}
}

// SaveConfigFileToDisk get config from secret and write it to disk
func (k *Kratos) SaveConfigFileToDisk(name, namespace, destination string) error {
	secret, err := k.Client.GetSecret(name+kratosSuffixConfig, namespace)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(destination, name)+".yaml", []byte(secret.Data[secretConfigKey]), fileMode); err != nil {
		return fmt.Errorf("writing yaml init file failed: %s", err)
	}

	return nil
}
