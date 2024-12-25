package main

import (
	"log"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Init     string   `yaml:"init"`
	Plugins  []string `yaml:"plugins"`
	AllowIPs []string `yaml:"allow_ips"`
	Tokens   []string `yaml:"tokens"`
}

var (
	cfg = &Config{AllowIPs: []string{"127.0.0.1"}}
	mu  sync.Mutex
)

// 加载配置
// method: auto 自动更新, signal 信号更新
func LoadConfig(filename string, method string) error {
	var config Config

	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 解析 YAML 内容到 Config 结构体
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	// 更新全局配置
	updateConfig(&config)

	// 监控文件变化
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(filename); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}
				switch method {
				case "auto":
					log.Println("配置文件已修改，重新加载配置...")
					if err := reloadConfig(filename); err != nil {
						log.Println("重新加载配置失败:", err)
					}
				case "signal":
					log.Println("配置文件已修改...")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("监控错误:", err)
			}
		}
	}()

	if method == "signal" {
		// 监听信号
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGUSR1)

		go func() {
			for {
				select {
				case sig := <-signalChan:
					log.Printf("接收到信号 %v，重新加载配置...", sig)
					if err := reloadConfig(filename); err != nil {
						log.Println("重新加载配置失败:", err)
					}
				}
			}
		}()
	}

	return nil
}

func updateConfig(newConfig *Config) {
	mu.Lock()
	defer mu.Unlock()

	// 深拷贝或逐字段更新
	*cfg = *newConfig

	// 确保 AllowIPs 包含 "127.0.0.1"
	if len(cfg.AllowIPs) == 0 || !slices.Contains(cfg.AllowIPs, "127.0.0.1") {
		cfg.AllowIPs = append(cfg.AllowIPs, "127.0.0.1")
	}
}

func reloadConfig(path string) error {
	var config Config

	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析 YAML 内容到 Config 结构体
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	// 更新全局配置
	updateConfig(&config)

	log.Println("配置已更新:", cfg)

	return nil
}
