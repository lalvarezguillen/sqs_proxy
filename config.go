package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type AppConfig struct {
	ProxySettingsFile string
	AWSRegion         string
}

func loadRequiredEnvVariable(n string) (string, error) {
	v := os.Getenv(n)
	if v == "" {
		return "", fmt.Errorf("Missing required %s environment variable", n)
	}
	return v, nil
}

func loadConfig() (*AppConfig, error) {
	f, err := loadRequiredEnvVariable("PROXY_SETTINGS_FILE")
	if err != nil {
		return nil, err
	}
	r, err := loadRequiredEnvVariable("AWS_REGION")
	if err != nil {
		return nil, err
	}
	return &AppConfig{ProxySettingsFile: f, AWSRegion: r}, nil
}

type ProxySettings struct {
	Src      string
	Dest     []string
	Interval time.Duration
}

type Settings struct {
	Proxies []ProxySettings
}

func loadProxySettings(fpath string) ([]ProxySettings, error) {
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	var config Settings
	json.Unmarshal(file, &config)
	return config.Proxies, nil
}
