package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Log configuration
const (
	LogPath = "tmp/service.log"

	// Max log file size in megabyte
	MaxSize = 1

	MaxBackups = 3

	// Max log retention in day
	MaxAge = 28
)

type Config struct {
	App struct {
		Service         string        `mapstructure:"service"`
		Environment     string        `mapstructure:"environment"`
		Timezone        string        `mapstructure:"timezone"`
		ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
		TargetHost      string        `mapstructure:"targetHost"`
		DestinationHost string        `mapstructure:"destinationHost"`
		ProxyTimeout    time.Duration `mapstructure:"proxyTimeout"`
		Mirrors         []struct {
			Name      string   `mapstructure:"name"`
			Methods   []string `mapstructure:"methods"`
			Endpoints []string `mapstructure:"endpoint"`
		} `mapstructure:"mirrors"`
	} `mapstructure:"app"`
	Log struct {
		Level      string `mapstructure:"level"`
		Path       string `mapstructure:"path"`
		MaxSize    int    `mapstructure:"maxSize"`
		MaxBackups int    `mapstructure:"maxBackups"`
		MaxAge     int    `mapstructure:"maxAge"`
	} `mapstructure:"log"`
	Server struct {
		Rest struct {
			Host    string `mapstructure:"host"`
			Port    int    `mapstructure:"port"`
			Prefork bool   `mapstructure:"prefork"`
		} `mapstructure:"rest"`
	} `mapstructure:"server"`
}

type RetryConfig struct {
	Max   int
	Delay time.Duration
}

func Load(path string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	var conf Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("%s", fmt.Sprintf("[Config] Error loading config: %v", err))
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("%s", fmt.Sprintf("[Config] Error unmarshaling config: %v", err))
	}

	return &conf
}
