package config

import (
	"fmt"

	"regexp"

	validator "github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation"
)

func labelsValidation(labels map[string]string) error {
	for name := range labels {
		errors := validation.IsValidLabelValue(name)
		if len(errors) > 0 {
			return fmt.Errorf("validation of labels %s failed: %s", name, errors[len(errors)-1])
		}
	}
	return nil
}

func mapKeyUniq(m1, m2 map[string]string) error {
	for name := range m1 {
		if _, found := m2[name]; found {
			return fmt.Errorf("common labels/annotations %q cannot be duplicated in others sections", name)
		}
	}
	return nil
}

func (c *Config) validateConfig() error {
	validate := validator.New()

	// validate config via struct yaml tag
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validation of configuration failed: %s", err)
	}

	// TODO integration of all validations with the validator
	// REF: https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go

	if c.Common != nil {
		if err := c.Common.validateConfig(); err != nil {
			return err
		}
	}

	if c.Deployment != nil {
		if err := c.Deployment.validateConfig(c.Common); err != nil {
			return err
		}
	}

	if c.Cronjob != nil {
		if err := c.Cronjob.validateConfig(c.Common); err != nil {
			return err
		}
	}

	return nil
}

// validateConfig validate common config
func (c *Common) validateConfig() error {
	if c.Labels != nil {
		if err := labelsValidation(c.Labels); err != nil {
			return err
		}
	}

	return nil
}

// validateConfig validate deployment config
func (d *Deployment) validateConfig(common *Common) error {
	if d.Labels != nil {
		if err := labelsValidation(d.Labels); err != nil {
			return err
		}
	}

	// common labels & annotations must be uniq
	if common != nil {
		if common.Labels != nil && d.Labels != nil {
			if err := mapKeyUniq(common.Labels, d.Labels); err != nil {
				return err
			}
		}
		if common.Annotations != nil && d.Annotations != nil {
			if err := mapKeyUniq(common.Annotations, d.Annotations); err != nil {
				return err
			}
		}
	}

	// containers
	for _, container := range d.Containers {
		if container.Resources != nil {
			if err := container.Resources.validateConfig(container.Name); err != nil {
				return err
			}
		}
	}

	// ingress
	if d.Ingress != nil {
		if err := d.Ingress.validateConfig(common); err != nil {
			return err
		}
	}

	return nil
}

// validateConfig validate deployment config
func (c *Cronjob) validateConfig(common *Common) error {
	if c.Labels != nil {
		if err := labelsValidation(c.Labels); err != nil {
			return err
		}
	}

	// common labels & annotations must be uniq
	if common != nil {
		if common.Labels != nil && c.Labels != nil {
			if err := mapKeyUniq(common.Labels, c.Labels); err != nil {
				return err
			}
		}
		if common.Annotations != nil && c.Annotations != nil {
			if err := mapKeyUniq(common.Annotations, c.Annotations); err != nil {
				return err
			}
		}
	}

	// cronjob schedule validation
	// TODO better validation, probably check how k8s handle this
	re := regexp.MustCompile(`(((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7}`)
	if !re.MatchString(c.Schedule) {
		return fmt.Errorf("cronjob schedule isn't valid")
	}

	if c.Container.Resources != nil {
		if err := c.Container.Resources.validateConfig(c.Container.Name); err != nil {
			return err
		}
	}

	return nil
}

// validateConfig validate ingress config
func (i *Ingress) validateConfig(common *Common) error {
	if i.Labels != nil {
		if err := labelsValidation(i.Labels); err != nil {
			return err
		}
	}

	// common labels & annotations must be uniq
	if common != nil {
		if common.Labels != nil && i.Labels != nil {
			if err := mapKeyUniq(common.Labels, i.Labels); err != nil {
				return err
			}
		}
		if common.Annotations != nil && i.Annotations != nil {
			if err := mapKeyUniq(common.Annotations, i.Annotations); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Resources) validateConfig(container string) error {
	if r.Requests != nil {
		if err := r.Requests.validateConfig(container, "requests"); err != nil {
			return err
		}
	}
	if r.Limits != nil {
		if err := r.Limits.validateConfig(container, "limits"); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResourceType) validateConfig(container, rType string) error {
	if _, err := resource.ParseQuantity(r.CPU); err != nil {
		return fmt.Errorf("validation of configuration resources failed: %s\ncontainer: %s -> %s cpu: %s", err, container, rType, r.CPU)
	}
	if _, err := resource.ParseQuantity(r.Memory); err != nil {
		return fmt.Errorf("validation of configuration resources failed: %s\ncontainer: %s -> %s memory: %s", err, container, rType, r.Memory)
	}

	return nil
}
