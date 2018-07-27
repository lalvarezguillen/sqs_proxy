package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRequiredEnvVariable(t *testing.T) {
	os.Setenv("dummy-var", "dummy-val")

	val, err := loadRequiredEnvVariable("dummy-var")
	assert.Nil(t, err)
	assert.Equal(t, "dummy-val", val)
}

func TestLoadConfigMissingData(t *testing.T) {
	_, err := loadConfig()
	assert.Error(t, err)
}

func TestLoadConfig(t *testing.T) {
	os.Setenv("PROXY_SETTINGS_FILE", "config.json")
	os.Setenv("AWS_REGION", "us-east-1")
	c, err := loadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "config.json", c.ProxySettingsFile)
	assert.Equal(t, "us-east-1", c.AWSRegion)
}

func TestLoadProxySettings(t *testing.T) {
	s := Settings{
		Proxies: []ProxySettings{
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
	b, _ := json.Marshal(&s)
	fname := "/tmp/dummy-config.json"
	ioutil.WriteFile(fname, b, 0644)

	p, err := loadProxySettings(fname)
	assert.NoError(t, err)
	assert.Equal(t, s.Proxies, p)
}
