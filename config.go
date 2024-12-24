package main

import (
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Init     string   `yaml:"init"`
	Plugins  []string `yaml:"plugins"`
	AllowIPs []string `yaml:"allow_ips"`
	Tokens   []string `yaml:"tokens"`
}

var cfg = &Config{
	AllowIPs: []string{"127.0.0.1"},
}

func LoadConfig(path string) error {
	var config Config

	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析 YAML 内容到 Config 结构体
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// 更新全局配置
	cfg = &config

	if len(config.AllowIPs) == 0 || slices.Contains(cfg.AllowIPs, "127.0.0.1") {
		config.AllowIPs = append(config.AllowIPs, "127.0.0.1")
	}

	return nil
}
