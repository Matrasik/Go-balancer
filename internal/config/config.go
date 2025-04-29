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

type ClientsCfg struct {
	Clients []BucketConfigs `json:"clients"`
}

type BucketConfigs struct {
	Addr         string `json:"addr"`
	BucketConfig `json:"params"`
}

type BucketConfig struct {
	Capacity int64 `json:"capacity"`
	Rate     int64 `json:"rate"`
}

type BucketDBConfig struct {
	IP       string `gorm:"type:varchar(32);uniqueIndex;not null;primaryKey" json:"ip"`
	Capacity int64  `gorm:"not null" json:"capacity"`
	Rate     int64  `gorm:"not null" json:"rate"`
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

func LoadBucketConfig(filename string) (*ClientsCfg, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	bucketCfg := &ClientsCfg{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&bucketCfg)
	if err != nil {
		return nil, err
	}
	return bucketCfg, nil
}
