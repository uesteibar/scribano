package yamlconfig

import (
	"errors"

	"github.com/uesteibar/scribano/watcher"
	"github.com/uesteibar/scribano/watcher/config"
	"gopkg.in/yaml.v2"
)

// YamlConfig loads watcher configuration from yaml file
type YamlConfig struct {
	Loader config.Loader
}

type parsedConfig struct {
	Host         string `yaml:"host"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchange_type"`
	RoutingKey   string `yaml:"routing_key"`
}

const defaultExchangeType = "topic"

func validConfig(c parsedConfig) bool {
	return c.Host != "" && c.Exchange != "" && c.RoutingKey != ""
}

func validConfigs(configs []parsedConfig) bool {
	for _, c := range configs {
		if !validConfig(c) {
			return false
		}
	}

	return true
}

func invalidConfigErr() error {
	return errors.New("Invalid config")
}

// New returns a new YamlConfig instance
func New(loader config.Loader) YamlConfig {
	return YamlConfig{Loader: loader}
}

// Parse yaml file into watcher configuration
func (c YamlConfig) Parse() ([]watcher.Config, error) {
	configs := []watcher.Config{}
	parsedConfigs := []parsedConfig{}

	data, err := c.Loader.Load()
	if err != nil {
		return configs, err
	}
	err = yaml.Unmarshal(data, &parsedConfigs)

	if err != nil {
		return configs, err
	}

	if !validConfigs(parsedConfigs) {
		return configs, invalidConfigErr()
	}

	for _, parsedConfig := range parsedConfigs {
		exchangeType := parsedConfig.ExchangeType
		if exchangeType == "" {
			exchangeType = defaultExchangeType
		}

		configs = append(
			configs,
			watcher.Config{
				Host:         parsedConfig.Host,
				Exchange:     parsedConfig.Exchange,
				ExchangeType: exchangeType,
				RoutingKey:   parsedConfig.RoutingKey,
			})
	}

	return configs, nil
}
