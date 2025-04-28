package config

import (
	"Golang_balancer/internal/balancer"
	"encoding/json"
	"net/url"
	"os"
)

type Config struct {
	Port         string                        `json:"port"`
	BackendsInfo []*balancer.BackendServerInfo `json:"backends"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	for i := range cfg.BackendsInfo {
		cfg.BackendsInfo[i].Address, err = url.Parse(cfg.BackendsInfo[i].UrlString)
		cfg.BackendsInfo[i].SetAlive(true)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
