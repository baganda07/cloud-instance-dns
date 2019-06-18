package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input     map[interface{}]interface{}
		domain    string
		port      string
		awsConfig *AwsConfig
		gcpConfig *GcpConfig
		err       bool
	}{
		"empty":       {input: nil, err: true},
		"emptyDomain": {input: make(map[interface{}]interface{}), err: true},
		"success": {input: map[interface{}]interface{}{
			"domain": "localhost"}, err: false, domain: "localhost."},
	}

	for _, t := range tests {
		do, _, ac, gc, err := ParseConfig(t.input)
		assert.Equal(t.domain, do)
		assert.Equal(t.awsConfig, ac)
		assert.Equal(t.gcpConfig, gc)
		if t.err {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}

	yamlPath := os.Getenv("TEST_YAML_PATH")
	if yamlPath != "" {
		config := make(map[interface{}]interface{})
		bys, err := ioutil.ReadFile(yamlPath)
		assert.NoError(err)
		err = yaml.Unmarshal(bys, &config)
		assert.NoError(err)

		// both enable
		do, po, ac, gc, err := ParseConfig(config)
		assert.NoError(err)
		assert.NotEqual("", do)
		assert.NotEqual("", po)
		assert.True(strings.HasSuffix(do, "."))
		assert.NotEmpty(ac)
		assert.NotEmpty(gc)
		assert.NotEqual(0, len(ac.clients))
		assert.NotEmpty(gc.projectId)
		assert.NotEmpty(gc.zones)
		assert.NotEmpty(gc.client)

		// gcp enable off
		config["gcp"].(map[interface{}]interface{})["enable"] = "false"
		do, _, ac, gc, err = ParseConfig(config)
		assert.NoError(err)
		assert.NotEqual("", do)
		assert.NotEmpty(ac)
		assert.Empty(gc)

		// aws enable off
		config["gcp"].(map[interface{}]interface{})["enable"] = "true"
		config["aws"].(map[interface{}]interface{})["enable"] = "false"
		do, _, ac, gc, err = ParseConfig(config)
		assert.NoError(err)
		assert.NotEqual("", do)
		assert.Empty(ac)
		assert.NotEmpty(gc)
	}
}