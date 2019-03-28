package config

import (
	"io/ioutil"
	"net/http"
)

// Loader is an interface that all config loaders should implement
type Loader interface {
	Load() ([]byte, error)
}

// PathLoader loads the content of a config file form a local path
type PathLoader struct {
	Source string
}

// Load configuration file content from local path
func (l *PathLoader) Load() ([]byte, error) {
	data, err := ioutil.ReadFile(l.Source)

	if err != nil {
		return data, err
	}

	return data, nil
}

// URLLoader loads the content of a config file form a local path
type URLLoader struct {
	Source string
}

// Load configuration file content from local path
func (l *URLLoader) Load() ([]byte, error) {
	resp, err := http.Get(l.Source)
	if err != nil {
		var emptyResp []byte
		return emptyResp, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}
