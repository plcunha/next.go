package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents Next.go configuration (similar to next.config.js)
type Config struct {
	// Server configuration
	Server ServerConfig `yaml:"server,omitempty"`

	// Build configuration
	Build BuildConfig `yaml:"build,omitempty"`

	// Experimental features
	Experimental ExperimentalConfig `yaml:"experimental,omitempty"`

	// Image optimization
	Images ImageConfig `yaml:"images,omitempty"`

	// Internationalization
	I18n I18nConfig `yaml:"i18n,omitempty"`

	// Redirects
	Redirects []Redirect `yaml:"redirects,omitempty"`

	// Rewrites
	Rewrites []Rewrite `yaml:"rewrites,omitempty"`

	// Headers
	Headers []Header `yaml:"headers,omitempty"`
}

type ServerConfig struct {
	Port        int    `yaml:"port,omitempty"`
	Host        string `yaml:"host,omitempty"`
	Compression bool   `yaml:"compression,omitempty"`
}

type BuildConfig struct {
	Output                string `yaml:"output,omitempty"` // standalone, export
	Target                string `yaml:"target,omitempty"`
	DistDir               string `yaml:"distDir,omitempty"`
	GenerateBuildManifest bool   `yaml:"generateBuildManifest,omitempty"`
}

type ExperimentalConfig struct {
	AppDir        bool `yaml:"appDir,omitempty"`
	Turbopack     bool `yaml:"turbopack,omitempty"`
	ServerActions bool `yaml:"serverActions,omitempty"`
}

type ImageConfig struct {
	Domains    []string `yaml:"domains,omitempty"`
	Formats    []string `yaml:"formats,omitempty"`
	DeviceSizes []int   `yaml:"deviceSizes,omitempty"`
}

type I18nConfig struct {
	Locales       []string `yaml:"locales,omitempty"`
	DefaultLocale string   `yaml:"defaultLocale,omitempty"`
}

type Redirect struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	Permanent   bool   `yaml:"permanent,omitempty"`
}

type Rewrite struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type Header struct {
	Source  string     `yaml:"source"`
	Headers []HeaderKV `yaml:"headers"`
}

type HeaderKV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Load loads configuration from nextgo.yaml or next.go.yaml
func Load(dir string) (*Config, error) {
	config := defaultConfig()

	// Try to load config file
	configFiles := []string{
		"nextgo.yaml",
		"nextgo.yml",
		"next.go.yaml",
		"next.go.yml",
	}

	for _, file := range configFiles {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return loadFromFile(path, config)
		}
	}

	return config, nil
}

// loadFromFile loads config from a YAML file
func loadFromFile(path string, config *Config) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return config, nil
}

// defaultConfig returns default configuration
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        3000,
			Host:        "localhost",
			Compression: true,
		},
		Build: BuildConfig{
			Output:               "standalone",
			DistDir:              ".next",
			GenerateBuildManifest: true,
		},
		Experimental: ExperimentalConfig{
			AppDir: true,
		},
		Images: ImageConfig{
			Domains:    []string{},
			Formats:    []string{"image/webp"},
			DeviceSizes: []int{640, 750, 828, 1080, 1200, 1920, 2048, 3840},
		},
	}
}

// Save saves the configuration to a file
func (c *Config) Save(dir string) error {
	path := filepath.Join(dir, "nextgo.yaml")

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
