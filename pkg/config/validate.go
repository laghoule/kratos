package config

import (
	"fmt"
	"regexp"

	"github.com/laghoule/kratos/pkg/common"

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
	} else {
		c.Common = &Common{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
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

	if c.ConfigMaps != nil {
		if err := c.ConfigMaps.validateConfig(c); err != nil {
			return err
		}
	}

	if c.Secrets != nil {
		if err := c.Secrets.validateConfig(c); err != nil {
			return err
		}
	}

	// validate that mount path are uniq by container for configmaps & secrets
	if err := c.validateMountPath(); err != nil {
		return err
	}

	return nil
}

// validateConfig validate common config
func (c *Common) validateConfig() error {
	if c.Labels != nil {
		if err := labelsValidation(c.Labels); err != nil {
			return err
		}
	} else {
		c.Labels = map[string]string{}
	}

	if c.Annotations == nil {
		c.Annotations = map[string]string{}
	}

	return nil
}

// validateLabelsAnnotations validate labels and annotations
func validateLabelsAnnotations(common *Common, labels, annotations map[string]string) error {
	if err := labelsValidation(labels); err != nil {
		return err
	}

	// common labels & annotations must be uniq
	if common != nil {
		if common.Labels != nil && labels != nil {
			if err := mapKeyUniq(common.Labels, labels); err != nil {
				return err
			}
		}
		if common.Annotations != nil && annotations != nil {
			if err := mapKeyUniq(common.Annotations, annotations); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateConfig validate deployment config
func (d *Deployment) validateConfig(common *Common) error {
	if err := validateLabelsAnnotations(common, d.Labels, d.Annotations); err != nil {
		return err
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
	if err := validateLabelsAnnotations(common, c.Labels, c.Annotations); err != nil {
		return err
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

// validateConfig validate ingress labels & annotations
func (i *Ingress) validateConfig(common *Common) error {
	if err := validateLabelsAnnotations(common, i.Labels, i.Annotations); err != nil {
		return err
	}

	return nil
}

// validateConfig validate resources requests and limits
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

// validateConfig validate CPU & Memory resources
func (r *ResourceType) validateConfig(container, rType string) error {
	if _, err := resource.ParseQuantity(r.CPU); err != nil {
		return fmt.Errorf("validation of configuration resources failed: %s\ncontainer: %s -> %s cpu: %s", err, container, rType, r.CPU)
	}
	if _, err := resource.ParseQuantity(r.Memory); err != nil {
		return fmt.Errorf("validation of configuration resources failed: %s\ncontainer: %s -> %s memory: %s", err, container, rType, r.Memory)
	}

	return nil
}

// TODO: use generic for secrets & config validation

// validateConfig validate secrets labels & annotations
func (s *Secrets) validateConfig(conf *Config) error {
	if err := validateLabelsAnnotations(conf.Common, s.Labels, s.Annotations); err != nil {
		return err
	}

	if err := validateExposedTo(s.Files, conf); err != nil {
		return err
	}

	return nil
}

// validateConfig validate configmaps labels & annotations
func (c *ConfigMaps) validateConfig(conf *Config) error {
	if err := validateLabelsAnnotations(conf.Common, c.Labels, c.Annotations); err != nil {
		return err
	}

	if err := validateExposedTo(c.Files, conf); err != nil {
		return err
	}

	return nil
}

// validateExposedTo check if exposedTo container exist in deployment or in cronjob
func validateExposedTo(files []File, conf *Config) error {
	exposedToFound := false
	for _, file := range files {
		for _, exposedTo := range file.Mount.ExposedTo {

			if conf.Deployment != nil {
				if common.ListContain(conf.Deployment.listContainerNames(), exposedTo) {
					exposedToFound = true
					continue
				}
			}

			if conf.Cronjob != nil {
				if conf.Cronjob.Container.Name == exposedTo {
					exposedToFound = true
					continue
				}
			}

		}
	}

	if !exposedToFound {
		return fmt.Errorf("exposedTo container don't exist")
	}

	return nil
}

// validateMountPath ensure mount path are uniq in container
func (c *Config) validateMountPath() error {
	mountPathFound := map[string]map[string]bool{}
	files := []File{}

	if c.ConfigMaps != nil {
		files = append(files, c.ConfigMaps.Files...)
	}

	if c.Secrets != nil {
		files = append(files, c.Secrets.Files...)
	}

	for _, file := range files {
		for _, exposedTo := range file.Mount.ExposedTo {
			if _, found := mountPathFound[exposedTo][file.Mount.Path]; found {
				return fmt.Errorf("mount.path %s must be uniq", file.Mount.Path)
			}
			mountPathFound[exposedTo] = map[string]bool{file.Mount.Path: true}
		}
	}

	return nil
}
