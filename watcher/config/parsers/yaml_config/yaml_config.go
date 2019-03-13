package yamlconfig

import (
	"errors"
	"io/ioutil"

	"github.com/uesteibar/asyncapi-watcher/watcher"
	"gopkg.in/yaml.v2"
)

// YamlConfig loads watcher configuration from yaml file
type YamlConfig struct {
	Path string
}

type config struct {
	Host       string `yaml:"host"`
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routing_key"`
}

func validConfig(c config) bool {
	return c.Host != "" && c.Exchange != "" && c.RoutingKey != ""
}

func invalidConfigErr() error {
	return errors.New("Invalid config")
}

// New returns a new YamlConfig instance
func New(path string) YamlConfig {
	return YamlConfig{Path: path}
}

// Parse yaml file into watcher configuration
func (c YamlConfig) Parse() (watcher.Config, error) {
	config := config{}
	data, err := ioutil.ReadFile(c.Path)
	if err == nil {
		err = yaml.Unmarshal(data, &config)

		if err == nil {
			if validConfig(config) {
				return watcher.Config{
					Host:       config.Host,
					Exchange:   config.Exchange,
					RoutingKey: config.RoutingKey,
				}, nil
			}

			err = invalidConfigErr()
		}
	}

	return watcher.Config{}, err
}
