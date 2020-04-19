package prometheus

import (
	yaml "gopkg.in/yaml.v3"
)

// Config represents a subset of Prometheus config
// that we are interested in.
type Config struct {
	RuleFiles []string `yaml:"rule_files,omitempty"`
}

// LoadConfig parses config from []byte
func LoadConfig(raw []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
