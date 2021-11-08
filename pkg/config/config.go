package config

import (
	"fmt"
	"os"

	validator "github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Config of kratos
type Config struct {
	Name      string `yaml:"name" validate:"required,alphanum"`
	Namespace string `yaml:"namespace" validate:"required,alphanum"`
	*Deployment
	*Service
	*Ingress
}

// Deployment object
type Deployment struct {
	Replicas int32 `yaml:"replicas,omitempty" validate:"gte=1,lte=100" `
}

// Service object
type Service struct {
	Port int32 `yaml:"port" validate:"required,gte=1,lte=65535"`
}

// Ingress object
type Ingress struct {
	IngressClass  string      `yaml:"ingressClass" validate:"required,alphanum"`
	ClusterIssuer string      `yaml:"clusterIssuer" validate:"required,alphanum"`
	Hostnames     []Hostnames `yaml:"hostnames" validate:"required,dive,hostname"`
	Port          int32       `yaml:"port" validate:"required,gte=1,lte=65535"`
}

// Hostnames use in ingress object
type Hostnames string

func validateConfig(config *Config) error {
	validate := &validator.Validate{}
	validate = validator.New()

	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validation of config failed: %s", err)
	}

	return nil
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
