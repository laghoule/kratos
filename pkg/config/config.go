package config

import (
	"fmt"
	"os"

	validator "github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	// DeployLabel is a managed-by k8s label for krator
	DeployLabel = "app.kubernetes.io/managed-by=kratos"
)

// Config of kratos
type Config struct {
	*Deployment
	*Ingress
}

// Deployment object
type Deployment struct {
	Replicas   int32       `yaml:"replicas,omitempty" validate:"gte=0,lte=100" `
	Containers []Container `yaml:"containers" validate:"required,dive"`
}

// Container object
type Container struct {
	Name      string    `yaml:"name" validate:"required,alphanum,lowercase"`
	Image     string    `yaml:"image" validate:"required,ascii"`
	Tag       string    `yaml:"tag" validate:"required,ascii"`
	Port      int32     `yaml:"port" validate:"required,gte=1,lte=65535"`
	Resources Resources `yaml:"resources,omitempty"`
}

// Resources objext
type Resources struct {
	Limits  ResourceType `yaml:"limits,omitempty"`
	Request ResourceType `yaml:"requests,omitempty"`
}

// ResourceType object
type ResourceType struct {
	// todo add validate
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// Ingress object
type Ingress struct {
	IngressClass  string      `yaml:"ingressClass" validate:"required,alphanum"`
	ClusterIssuer string      `yaml:"clusterIssuer" validate:"required,alphanum"`
	Hostnames     []Hostnames `yaml:"hostnames" validate:"required,dive,hostname"`
}

// Hostnames use in ingress object
type Hostnames string

// String implement the stringer interface
func (h *Hostnames) String() string {
	return string(*h)
}

func validateConfig(config *Config) error {
	validate := &validator.Validate{}
	validate = validator.New()

	for _, container := range config.Containers {
		resources := map[string]string{
			"requests cpu":    container.Resources.Request.CPU,
			"requests memory": container.Resources.Request.Memory,
			"limits cpu":      container.Resources.Limits.CPU,
			"limits memory":   container.Resources.Limits.Memory,
		}
		for rsName, rsValue := range resources {
			if rsValue == "" {
				continue
			}
			_, err := resource.ParseQuantity(rsValue)
			if err != nil {
				return fmt.Errorf("validating configuration resources failed: %s\ncontainer: %s -> %s: %s", err, container.Name, rsName, rsValue)
			}
		}
	}

	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validation of config failed: %s", err)
	}

	return nil
}

// CreateInit return an sample config
func CreateInit() *Config {
	return &Config{
		Deployment: &Deployment{
			Replicas: 1,
			Containers: []Container{
				{
					Name:  "example",
					Image: "nginx",
					Tag:   "latest",
					Port:  8080,
				},
			},
		},
		Ingress: &Ingress{
			IngressClass:  "nginx",
			ClusterIssuer: "letsencrypt",
			Hostnames: []Hostnames{
				"www.example.com",
			},
		},
	}
}

// Load configuration of the specified file
func (c *Config) Load(file string) error {
	configFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read configuration failed: %s", err)
	}

	if err := yaml.Unmarshal(configFile, c); err != nil {
		return fmt.Errorf("unmarshaling yaml failed: %s", err)
	}

	if err := validateConfig(c); err != nil {
		return err
	}

	return nil
}
