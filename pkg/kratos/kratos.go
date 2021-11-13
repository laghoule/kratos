package kratos

import (
	"fmt"
	"os"
	"strconv"

	"github.com/laghoule/kratos/pkg/certmanager"
	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kratosSuffixConfig = "-kratos-config"
	kratosConfigKey    = "config"
)

// Kratosphere is the kratos interface
type Kratosphere interface {
	List(namespace string) error
	Create(name, namespace, image, tag, ingresClass, clusterIssuer string, hostnames []string, replicas, port int32) error
	Delete(name, namespace string)
}

// Kratos contains info for deployment
type Kratos struct {
	*k8s.Client
	*config.Config
}

// New return a kratos struct
func New(confFile string) (*Kratos, error) {
	kclient, err := k8s.New()
	if err != nil {
		return nil, err
	}

	// check if we meet k8s version requirement
	kclient.CheckVersionDepency()

	// loading configuration
	confYAML := &config.Config{}
	if confFile != "" {
		if err := confYAML.Load(confFile); err != nil {
			return nil, err
		}
	}

	return &Kratos{
		Client: kclient,
		Config: confYAML,
	}, nil
}

// Create the deployment of all objects
func (k *Kratos) Create(name, namespace string) error {
	runWithError := false

	// validate clusterIssuer
	cm, err := certmanager.New(*k.Client)
	if err != nil {
		return err
	}

	if err := cm.CheckClusterIssuer(k.Client, k.Config.ClusterIssuer); err != nil {
		return err
	}

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

	// save config in secret
	spinner, _ = pterm.DefaultSpinner.Start("saving configuration")
	if err := k.saveConfigFile(name+kratosSuffixConfig, namespace); err != nil {
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

// CreateInit create sample configuration file
func (k *Kratos) CreateInit(file string) error {
	b, err := yaml.Marshal(config.CreateInit())
	if err != nil {
		return fmt.Errorf("marshaling yaml failed: %s", err)
	}

	if err := os.WriteFile(file, b, 0666); err != nil {
		return fmt.Errorf("writing yaml init file failed: %s", err)
	}

	return nil
}

func (k *Kratos) saveConfigFile(name, namespace string) error {
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
			kratosConfigKey: data,
		},
	}
}
