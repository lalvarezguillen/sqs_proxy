package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type ProxySettings struct {
	Src      string
	Dest     []string
	Interval time.Duration
}

type AppConfig struct {
	ProxyOps []ProxySettings
}

// Pretty produces a pretty JSON representation of the AppConfig object.
func (c *AppConfig) Pretty() string {
	b, _ := json.MarshalIndent(c, "", "    ")
	return string(b)
}

func loadConfig(fpath string) (*AppConfig, error) {
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return &AppConfig{}, err
	}
	var c AppConfig
	json.Unmarshal(file, &c)
	return &c, nil
}
