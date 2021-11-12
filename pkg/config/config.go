package config

import (
	"fmt"
	"os"

	validator "github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
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
	Name  string `yaml:"name" validate:"required,alphanum,lowercase"`
	Image string `yaml:"image" validate:"required,ascii"`
	Tag   string `yaml:"tag" validate:"required,ascii"`
	Port  int32  `yaml:"port" validate:"required,gte=1,lte=65535"`
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
