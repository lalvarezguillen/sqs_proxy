package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func loadRequiredEnvVariable(n string) (string, error) {
	v := os.Getenv(n)
	if v == "" {
		return "", fmt.Errorf("Missing required %s environment variable", n)
	}
	return v, nil
}

type ProxySettings struct {
	Src      string
	Dest     []string
	Interval time.Duration
}

type AppConfig struct {
	ProxyOps []ProxySettings
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
