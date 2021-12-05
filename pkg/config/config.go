package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	// DeployLabel is a managed-by k8s label for kratos
	DeployLabel = "app.kubernetes.io/managed-by=kratos"
)

// Config of kratos
type Config struct {
	*Common     `yaml:"common,omitempty"`
	*Cronjob    `yaml:"cronjob,omitempty"`
	*Deployment `yaml:"deployment,omitempty"`
	*Ingress    `yaml:"ingress,omitempty" validate:"required_with=deployment"`
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
	Name      string     `yaml:"name" validate:"required,alphanum,lowercase"`
	Image     string     `yaml:"image" validate:"required,ascii"`
	Tag       string     `yaml:"tag" validate:"required,ascii"`
	Resources *Resources `yaml:"resources,omitempty"`
}

// Resources objext
type Resources struct {
	Requests *ResourceType `yaml:"requests,omitempty"`
	Limits   *ResourceType `yaml:"limits,omitempty"`
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
	Hostnames     []string          `yaml:"hostnames" validate:"required,dive,hostname"`
}

// Cronjob object
type Cronjob struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Schedule    string            `yaml:"schedule" validate:"required"`
	Retry       int32             `yaml:"retry,omitempty" validate:"gte=0,lte=100"`
	Container   *Container        `yaml:"container" validate:"required"`
}

// CreateInit return an sample config
func CreateInit() *Config {
	return &Config{
		Common: &Common{
			Labels: map[string]string{
				"commonlabel": "value",
			},
			Annotations: map[string]string{
				"commonannotation": "value",
			},
		},
		Deployment: &Deployment{
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
					Resources: &Resources{
						Requests: &ResourceType{
							CPU:    "25m",
							Memory: "32Mi",
						},
						Limits: &ResourceType{
							CPU:    "50m",
							Memory: "64Mi",
						},
					},
				},
			},
		},
		Cronjob: &Cronjob{
			Labels: map[string]string{
				"label": "value",
			},
			Annotations: map[string]string{
				"annotation": "value",
			},
			Schedule: "0 0 * * *",
			Retry:    3,
			Container: &Container{
				Name:  "example",
				Image: "cronjobimage",
				Tag:   "latest",
				Resources: &Resources{
					Requests: &ResourceType{
						CPU:    "25m",
						Memory: "32Mi",
					},
					Limits: &ResourceType{
						CPU:    "50m",
						Memory: "64Mi",
					},
				},
			},
		},
		Ingress: &Ingress{
			Labels: map[string]string{
				"label": "value",
			},
			Annotations: map[string]string{
				"annotation": "value",
			},
			IngressClass:  "nginx",
			ClusterIssuer: "letsencrypt",
			Hostnames: []string{
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

	if err := c.validateConfig(); err != nil {
		return err
	}

	return nil
}

// LoadFromString load configuration from a specified string
func (c *Config) LoadFromString(conf string) error {
	if err := yaml.Unmarshal([]byte(conf), c); err != nil {
		return fmt.Errorf("unmarshaling yaml failed: %s", err)
	}

	return nil
}
