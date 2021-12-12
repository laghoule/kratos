package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	// DeployLabel is a managed-by k8s label for kratos
	DeployLabel = "app.kubernetes.io/managed-by=kratos"
	// LiveProbe represent the live config keyword
	LiveProbe = "live"
	// ReadyProbe represent the ready config keyword
	ReadyProbe = "ready"
)

// Config of kratos
type Config struct {
	*Common     `yaml:"common,omitempty"`
	*Cronjob    `yaml:"cronjob,omitempty"`
	*Deployment `yaml:"deployment,omitempty"`
}

// Common represent the common fields
type Common struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Deployment represent the Kubernetes deployment
type Deployment struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Replicas    int32             `yaml:"replicas,omitempty" validate:"required,gte=0,lte=100" `
	Port        int32             `yaml:"port" validate:"required,gte=1,lte=65535"`
	Containers  []Container       `yaml:"containers" validate:"required,dive"`
	Ingress     *Ingress          `yaml:"ingress" validate:"required,dive"`
}

// Container represent the Kubernetes container
type Container struct {
	Name      string     `yaml:"name" validate:"required,alphanum,lowercase"`
	Image     string     `yaml:"image" validate:"required,ascii"`
	Tag       string     `yaml:"tag" validate:"required,ascii"`
	Resources *Resources `yaml:"resources,omitempty"`
	Health    *Health    `yaml:"health,omitempty"`
}

// Health represent the healthcheck for the container
type Health struct {
	Live  *Check `yaml:"live,omitempty"`
	Ready *Check `yaml:"ready,omitempty"`
}

// Check represent the information about the healthcheck
type Check struct {
	Probe               string `yaml:"probe" validate:"required,uri"`
	Port                int32  `yaml:"port" validate:"required,gte=1,lte=65535"`
	InitialDelaySeconds int32  `yaml:"initialDelaySeconds,omitempty" validate:"omitempty,gte=1,lte=600"`
	PeriodSeconds       int32  `yaml:"periodSeconds,omitempty" validate:"omitempty,gte=1,lte=600"`
}

// Resources represent requests and limits allocations
type Resources struct {
	Requests *ResourceType `yaml:"requests,omitempty"`
	Limits   *ResourceType `yaml:"limits,omitempty"`
}

// ResourceType represent CPU & Memory allocations
type ResourceType struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// Ingress represent the Kubernetes ingress
type Ingress struct {
	Labels        map[string]string `yaml:"labels,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
	IngressClass  string            `yaml:"ingressClass" validate:"required,alphanum"`
	ClusterIssuer string            `yaml:"clusterIssuer" validate:"required,alphanum"`
	Hostnames     []string          `yaml:"hostnames" validate:"required,dive,hostname"`
}

// Cronjob represent the Kubernetes cronjobs
type Cronjob struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Schedule    string            `yaml:"schedule" validate:"required"`
	Retry       int32             `yaml:"retry,omitempty" validate:"omitempty,gte=0,lte=100"`
	Container   *Container        `yaml:"container" validate:"required"`
}

// Configmaps represent the Kubernetes configmaps
type Configmaps struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	*File
}

// Secrets represent the Kubernetes secrets
type Secrets struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	*File
}

// File contains secrets and configmaps informations
type File struct {
	Name       string   `yaml:"name" validate:"required"`
	MountPath  string   `yaml:"mountPath" validate:"required, dir"`
	Data       string   `yaml:"data" validate:"requires"`
	Containers []string `yaml:"containers" validate:"required"`
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
					Health: &Health{
						Live: &Check{
							Probe:               "/isLive",
							Port:                8080,
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
						},
						Ready: &Check{
							Probe:               "/isReady",
							Port:                8080,
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
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
