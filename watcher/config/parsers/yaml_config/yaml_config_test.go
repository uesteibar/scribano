package yamlconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/watcher"
)

func TestParse(t *testing.T) {
	parser := New("../../../../fixtures/test/yaml_config.yml")
	config, err := parser.Parse()

	assert.Nil(t, err)
	expected := watcher.Config{
		Host:       "amqp://guest:guest@localhost",
		Exchange:   "/",
		RoutingKey: "#",
	}
	assert.Equal(t, expected, config)
}

func TestParse_NoFile(t *testing.T) {
	parser := New("non_existing_file")
	_, err := parser.Parse()

	assert.NotNil(t, err)
}

func TestParse_InvalidFile(t *testing.T) {
	parser := New("../../../../fixtures/test/invalid_yaml_config.yml")
	_, err := parser.Parse()

	assert.NotNil(t, err)
}
