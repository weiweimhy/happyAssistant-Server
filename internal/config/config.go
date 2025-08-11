package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ServerConfig Server配置
type ServerConfig struct {
	Port  int    `yaml:"port"`
	Route string `yaml:"route"`
}

// MongoConfig MongoDB配置
type MongoConfig struct {
	URI         string        `yaml:"uri"`
	Database    string        `yaml:"database"`
	Username    string        `yaml:"username"`
	Password    string        `yaml:"password"`
	MaxPoolSize int           `yaml:"maxPoolSize"`
	Timeout     time.Duration `yaml:"timeout"`
	OpTimeout   time.Duration `yaml:"opTimeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

// Config 总配置
type Config struct {
	Server  ServerConfig `yaml:"server"`
	MongoDB MongoConfig  `yaml:"mongodb"`
	Log     LogConfig    `yaml:"logger"`
}

var Cfg Config

func LoadConfig(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open config file error: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("close config file error: %v\n", err)
		}
	}(f)
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&Cfg); err != nil {
		log.Fatalf("decode config yaml error: %v", err)
	}
	// 设置默认操作超时
	if Cfg.MongoDB.OpTimeout == 0 {
		Cfg.MongoDB.OpTimeout = 5 * time.Second
	}
}
