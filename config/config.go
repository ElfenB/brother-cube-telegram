package config

import (
	"brother-cube-telegram/logger"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Holds all configuration values for the application
type Config struct {
	Printer PrinterConfig `yaml:"printer"`
	GPIO    GPIOConfig    `yaml:"gpio"`
	Logging LoggingConfig `yaml:"logging"`
}

// Holds printer-specific configuration
type PrinterConfig struct {
	// Number of retry attempts when connecting to printer
	RetryAttempts int `yaml:"retry_attempts"`

	// Auto-shutdown delay in minutes
	AutoShutdownDelayMinutes int `yaml:"auto_shutdown_delay_minutes"`

	// Path to store draft/preview images
	DraftsFolder string `yaml:"drafts_folder"`

	// Font size for label printing
	FontSize int `yaml:"font_size"`

	// Base delay in seconds for retries
	RetryBaseDelaySeconds int `yaml:"retry_base_delay_seconds"`

	// File permissions for created directories (octal)
	FolderPermissions int `yaml:"folder_permissions"`
}

// GPIOConfig holds GPIO-specific configuration
type GPIOConfig struct {
	// GPIO pin number for the relay
	RelayPin int `yaml:"relay_pin"`
}

// LoggingConfig holds logging-specific configuration
type LoggingConfig struct {
	// Log level: DEBUG, INFO, WARN, ERROR
	Level string `yaml:"level"`
}

// Global config instance
var cfg *Config

// Loads the configuration from the specified YAML file
func Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Process the drafts folder path (expand ~ to home directory)
	cfg.Printer.DraftsFolder = expandPath(cfg.Printer.DraftsFolder)

	return nil
}

// Returns the loaded configuration
func Get() *Config {
	if cfg == nil {
		panic("configuration not loaded - call config.Load() first")
	}
	return cfg
}

// Returns the auto-shutdown delay as a time.Duration
func (p *PrinterConfig) GetAutoShutdownDelay() time.Duration {
	return time.Duration(p.AutoShutdownDelayMinutes) * time.Minute
}

// Returns the retry delay for a given attempt number
func (p *PrinterConfig) GetRetryDelay(attemptNumber int) time.Duration {
	return time.Duration(attemptNumber+p.RetryBaseDelaySeconds) * time.Second
}

// Returns the folder permissions as os.FileMode
func (p *PrinterConfig) GetFolderPermissions() os.FileMode {
	return os.FileMode(p.FolderPermissions)
}

// Returns the logger.LogLevel from the string configuration
func (l *LoggingConfig) GetLogLevel() logger.LogLevel {
	switch strings.ToUpper(l.Level) {
	case "DEBUG":
		return logger.DEBUG
	case "INFO":
		return logger.INFO
	case "WARN":
		return logger.WARN
	case "ERROR":
		return logger.ERROR
	default:
		return logger.INFO // Default to INFO if invalid level
	}
}

// Expands ~ to the user's home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path // Return original path if we can't get home dir
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// Loads configuration from the default location (config.yaml in the project root)
func LoadDefault() error {
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	configPath := filepath.Join(filepath.Dir(execPath), "config.yaml")

	// If config.yaml doesn't exist next to executable, try current working directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config.yaml"
	}

	return Load(configPath)
}
