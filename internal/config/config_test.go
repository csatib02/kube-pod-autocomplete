package config

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		wantConfig *Config
		wantError  error
	}{
		{
			name:    "Default values",
			envVars: map[string]string{},
			wantConfig: &Config{
				ListenAddress:    defaultListenAddress,
				TrustedProxies:   nil,
				Mode:             defaultMode,
				LogLevel:         "info",
				JSONLog:          false,
				LogServerAddress: "",
			},
			wantError: nil,
		},
		{
			name: "Custom values",
			envVars: map[string]string{
				"KPA_LISTEN_ADDRESS":  "127.0.0.1:9090",
				"KPA_TRUSTED_PROXIES": "192.168.1.1,192.168.1.2",
				"KPA_MODE":            gin.ReleaseMode,
				"KPA_LOG_LEVEL":       "debug",
				"KPA_JSON_LOG":        "true",
				"KPA_LOG_SERVER":      "logserver.local",
			},
			wantConfig: &Config{
				ListenAddress:    "127.0.0.1:9090",
				TrustedProxies:   []string{"192.168.1.1", "192.168.1.2"},
				Mode:             gin.ReleaseMode,
				LogLevel:         "debug",
				JSONLog:          true,
				LogServerAddress: "logserver.local",
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			for envKey, envVal := range ttp.envVars {
				os.Setenv(envKey, envVal)
			}
			t.Cleanup(func() {
				os.Clearenv()
			})

			config, err := LoadConfig()
			if err != nil {
				assert.EqualError(t, ttp.wantError, err.Error(), "Unexpected error message")
			}

			if ttp.wantConfig != nil {
				assert.Equal(t, ttp.wantConfig, config, "Unexpected config")
			}
		})
	}
}
