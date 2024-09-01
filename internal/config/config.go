package config

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	defaultEnvPrefix     = "KPA"
	defaultListenAddress = "0.0.0.0:8080"
	defaultMode          = gin.DebugMode
)

type Config struct {
	ListenAddress    string   `mapstructure:"listen_address"`
	TrustedProxies   []string `mapstructure:"trusted_proxies"`
	Mode             string   `mapstructure:"mode"`
	LogLevel         string   `mapstructure:"log_level"`
	JSONLog          bool     `mapstructure:"json_log"`
	LogServerAddress string   `mapstructure:"log_server"`
}

func LoadConfig() (*Config, error) {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("."),
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
	)

	v.SetEnvPrefix("KPA")
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("listen_address")
	v.SetDefault("listen_address", defaultListenAddress)

	_ = v.BindEnv("trusted_proxies")
	v.SetDefault("trusted_proxies", nil)

	_ = v.BindEnv("mode")
	v.SetDefault("mode", gin.DebugMode)

	_ = v.BindEnv("log_level")
	v.SetDefault("log_level", "info")

	_ = v.BindEnv("json_log")

	_ = v.BindEnv("log_server")

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return config, nil
}
