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

func labelsValidation(labels map[string]string) error {
	for name := range labels {
		errors := validation.IsValidLabelValue(name)
		if len(errors) > 0 {
			return fmt.Errorf("validation of labels %s failed: %s", name, errors[len(errors)-1])
		}
	}
	return nil
}

func (c *Config) ensureNoNil() {
	// common
	if c.Common == nil {
		c.Common = &Common{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		}
	}

	// cronjobs
	if c.Cronjob == nil {
		c.Cronjob = &Cronjob{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
			Container: &Container{
				Resources: &Resources{
					Requests: &ResourceType{},
					Limits:   &ResourceType{},
				},
			},
		}
	} else {
		if c.Cronjob.Container.Resources == nil {
			c.Cronjob.Container.Resources = &Resources{
				Requests: &ResourceType{},
				Limits:   &ResourceType{},
			}
		} else {
			if c.Cronjob.Container.Resources.Requests == nil {
				c.Cronjob.Container.Resources.Requests = &ResourceType{}
			}
			if c.Cronjob.Container.Resources.Limits == nil {
				c.Cronjob.Container.Resources.Limits = &ResourceType{}
			}
		}
	}

	// deployment
	for i, container := range c.Deployment.Containers {
		if container.Resources == nil {
			c.Deployment.Containers[i].Resources = &Resources{}
			continue
		}
		if container.Resources.Requests == nil {
			c.Deployment.Containers[i].Resources.Requests = &ResourceType{}
			continue
		}
		if container.Resources.Limits == nil {
			c.Deployment.Containers[i].Resources.Limits = &ResourceType{}
			continue
		}
	}
}

func (c *Config) validateConfig() error {
	validate := &validator.Validate{}
	validate = validator.New()

	// validate config via struct yaml tag
	// must be checked before `ensureNoNil`
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validation of configuration failed: %s", err)
	}

	// TODO integration of all validations with the validator
	// REF: https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go

	// replace nil value in config
	c.ensureNoNil()

	labelsList := []map[string]string{
		c.Common.Labels,
		c.Deployment.Labels,
		c.Cronjob.Labels,
		c.Ingress.Labels,
	}

	// validate labels
	for _, labels := range labelsList {
		if labels != nil {
			if err := labelsValidation(labels); err != nil {
				return err
			}
		}
	}

	// TODO find a way to simplify these statements

	// common labels must be uniq
	for name := range c.Common.Labels {
		if _, found := c.Deployment.Labels[name]; found {
			return fmt.Errorf("common labels %q cannot be duplicated in deployment labels", name)
		}
		if _, found := c.Cronjob.Labels[name]; found {
			return fmt.Errorf("common labels %q cannot be duplicated in cronjobs labels", name)
		}
		if _, found := c.Ingress.Labels[name]; found {
			return fmt.Errorf("common labels %q cannot be duplicated in ingress labels", name)
		}
	}

	// common annotations must be uniq
	for name := range c.Common.Annotations {
		if _, found := c.Deployment.Annotations[name]; found {
			return fmt.Errorf("common annotations %q cannot be duplicated in deployment annotations", name)
		}
		if _, found := c.Cronjob.Annotations[name]; found {
			return fmt.Errorf("common annotations %q cannot be duplicated in cronjobs annotations", name)
		}
		if _, found := c.Ingress.Annotations[name]; found {
			return fmt.Errorf("common annotations %q cannot be duplicated in ingress annotations", name)
		}
	}

	// TODO validate cronjobs schedule
	// via regex: https://stackoverflow.com/questions/14203122/create-a-regular-expression-for-cron-statement

	// validate resource limits/requests
	for _, container := range c.Deployment.Containers {
		resources := map[string]string{
			"requests cpu":    container.Resources.Requests.CPU,
			"requests memory": container.Resources.Requests.Memory,
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

	return nil
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
				Image: "cronjobsimage",
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
