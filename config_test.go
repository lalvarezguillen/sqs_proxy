package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	conf := AppConfig{
		ProxyOps: []ProxySettings{
			ProxySettings{
				Src:      "dummy-source-1",
				Dest:     []string{"dummy-destination-1", "dummy-destination-2"},
				Interval: 20,
			},
			ProxySettings{
				Src:      "dummy-source-2",
				Dest:     []string{"dummy-destination-3", "dummy-destination-4"},
				Interval: 40,
			},
		},
	}
	b, _ := json.Marshal(&conf)
	fname := "/tmp/dummy-config.json"
	ioutil.WriteFile(fname, b, 0644)
	defer os.Remove(fname)

	c, err := loadConfig(fname)
	assert.NoError(t, err)
	assert.Equal(t, conf.ProxyOps, c.ProxyOps)
}

func TestPretty(t *testing.T) {
	conf := AppConfig{}
	pretty := conf.Pretty()
	assert.IsType(t, "", pretty)
}
