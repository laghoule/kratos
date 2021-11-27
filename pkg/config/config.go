package config

import (
	"fmt"
	"os"

	validator "github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	// DeployLabel is a managed-by k8s label for krator
	DeployLabel = "app.kubernetes.io/managed-by=kratos"
)

// Config of kratos
type Config struct {
	Common Common `yaml:"common,omitempty"`
	Deployment
	Ingress
}

// Common object
type Common struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Deployment object
type Deployment struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Replicas    int32             `yaml:"replicas,omitempty" validate:"required,gte=0,lte=100" `
	Port        int32             `yaml:"port" validate:"required,gte=1,lte=65535"`
	Containers  []Container       `yaml:"containers" validate:"required,dive"`
}

// Container object
type Container struct {
	Name      string    `yaml:"name" validate:"required,alphanum,lowercase"`
	Image     string    `yaml:"image" validate:"required,ascii"`
	Tag       string    `yaml:"tag" validate:"required,ascii"`
	Resources Resources `yaml:"resources,omitempty"`
}

// Resources objext
type Resources struct {
	Limits  ResourceType `yaml:"limits,omitempty"`
	Request ResourceType `yaml:"requests,omitempty"`
}

// ResourceType object
type ResourceType struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// Ingress object
type Ingress struct {
	Labels        map[string]string `yaml:"labels,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
	IngressClass  string            `yaml:"ingressClass" validate:"required,alphanum"`
	ClusterIssuer string            `yaml:"clusterIssuer" validate:"required,alphanum"`
	Hostnames     []Hostnames       `yaml:"hostnames" validate:"required,dive,hostname"`
}

// Hostnames use in ingress object
type Hostnames string

// String implement the stringer interface
func (h *Hostnames) String() string {
	return string(*h)
}

func labelsValidation(labels map[string]string) error {
	for name := range labels {
		errors := validation.IsValidLabelValue(name)
		if len(errors) > 0 {
			return fmt.Errorf("validation of labels %s failed: %s", name, errors[len(errors)-1])
		}
	}
	return nil
}

func validateConfig(config *Config) error {
	validate := &validator.Validate{}
	validate = validator.New()

	if err := labelsValidation(config.Common.Labels); err != nil {
		return err
	}

	// validate deployment labels
	if err := labelsValidation(config.Deployment.Labels); err != nil {
		return err
	}

	// validate ingress labels
	if err := labelsValidation(config.Ingress.Labels); err != nil {
		return err
	}

	// common labels must be uniq
	for name := range config.Common.Labels {
		if _, found := config.Deployment.Labels[name]; found {
			return fmt.Errorf("common labels %q cannot be duplicated in deployment labels", name)
		}
		if _, found := config.Ingress.Labels[name]; found {
			return fmt.Errorf("common labels %q cannot be duplicated in ingress labels", name)
		}
	}

	// common annotations must be uniq
	for name := range config.Common.Annotations {
		if _, found := config.Deployment.Annotations[name]; found {
			return fmt.Errorf("common annotations %q cannot be duplicated in deployment annotations", name)
		}
		if _, found := config.Ingress.Annotations[name]; found {
			return fmt.Errorf("common annotations %q cannot be duplicated in ingress annotations", name)
		}
	}

	// validate resource limits/requests
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
				return fmt.Errorf("validation of configuration resources failed: %s\ncontainer: %s -> %s: %s", err, container.Name, rsName, rsValue)
			}
		}
	}

	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validation of configuration failed: %s", err)
	}

	return nil
}

// CreateInit return an sample config
func CreateInit() *Config {
	return &Config{
		Common: Common{
			Labels: map[string]string{
				"commonlabel": "value",
			},
			Annotations: map[string]string{
				"commonannotation": "value",
			},
		},
		Deployment: Deployment{
			Labels: map[string]string{
				"label": "value",
			},
			Annotations: map[string]string{
				"annotation": "value",
			},
			Replicas: 1,
			Port:     8080,
			Containers: []Container{
				{
					Name:  "example",
					Image: "nginx",
					Tag:   "latest",
				},
			},
		},
		Ingress: Ingress{
			Labels: map[string]string{
				"label": "value",
			},
			Annotations: map[string]string{
				"annotation": "value",
			},
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
