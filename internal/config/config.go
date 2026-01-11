package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Token  string `mapstructure:"token"`
	Format string `mapstructure:"format"`
}

func (c *Config) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

// CheckLegacyConfig checks if the old homey-cli config exists and prints migration instructions
func CheckLegacyConfig() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	oldDir := filepath.Join(configDir, "homey-cli")
	newDir := filepath.Join(configDir, "homeyctl")

	// Check if old config exists and new one doesn't
	if _, err := os.Stat(filepath.Join(oldDir, "config.toml")); err == nil {
		if _, err := os.Stat(filepath.Join(newDir, "config.toml")); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "⚠️  Found config from previous version (homey-cli)")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "The binary has been renamed from 'homey' to 'homeyctl' to avoid")
			fmt.Fprintln(os.Stderr, "conflicts with Athom's official Homey CLI for app development.")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "To migrate your config, run:")
			fmt.Fprintf(os.Stderr, "  mv %s %s\n", oldDir, newDir)
			fmt.Fprintln(os.Stderr, "")
		}
	}
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// Config locations
	configDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(configDir, "homeyctl"))
	}
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("HOMEY")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 4859)
	viper.SetDefault("format", "json")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}

	dir := filepath.Join(configDir, "homeyctl")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	viper.Set("host", cfg.Host)
	viper.Set("port", cfg.Port)
	viper.Set("token", cfg.Token)
	viper.Set("format", cfg.Format)

	configPath := filepath.Join(dir, "config.toml")
	return viper.WriteConfigAs(configPath)
}
