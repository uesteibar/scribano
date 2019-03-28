package config

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PathLoader_Load(t *testing.T) {
	loader := PathLoader{Source: "../../fixtures/test/yaml_config.yml"}
	content, err := loader.Load()

	assert.Nil(t, err)
	expected, err := ioutil.ReadFile(loader.Source)
	assert.Equal(t, expected, content)
}

func Test_URLLoader_Load(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/", req.URL.String())
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	loader := URLLoader{Source: server.URL}
	content, err := loader.Load()

	assert.Nil(t, err)
	assert.Equal(t, []byte(`OK`), content)
}
