package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Jenkins    JenkinsConfig         `yaml:"jenkins"`
	MCP        MCPConfig             `yaml:"mcp"`
	Tools      map[string]ToolConfig `yaml:"tools"`
	Pagination PaginationConfig      `yaml:"pagination"`
}

type PaginationConfig struct {
	MaxCharsPerPage         int64 `yaml:"max_chars_per_page"`
	MaxContextLines         int64 `yaml:"max_context_lines"`
	MaxSearchResultsPerPage int64 `yaml:"max_search_results_per_page"`
}

type JenkinsConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MCPConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ToolConfig struct {
	IsEnabled   bool            `yaml:"is-enabled"`
	Name        string          `yaml:"name"`
	Description ToolDescription `yaml:"description"`
}

type ToolDescription struct {
	RU string `yaml:"ru"`
	EN string `yaml:"en"`
}

// Tool returns config for a named tool, defaulting to enabled if not configured.
func (c *Config) Tool(key string) ToolConfig {
	if t, ok := c.Tools[key]; ok {
		return t
	}
	return ToolConfig{IsEnabled: true, Name: key, Description: ToolDescription{EN: key, RU: key}}
}

func Load() (*Config, error) {
	path := "config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		path = p
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Jenkins.URL == "" {
		return fmt.Errorf("jenkins.url is required")
	}
	if c.Jenkins.Username == "" {
		return fmt.Errorf("jenkins.username is required")
	}
	if c.Jenkins.Password == "" {
		return fmt.Errorf("jenkins.password is required")
	}
	if c.MCP.Name == "" {
		c.MCP.Name = "jenkins"
	}
	if c.MCP.Version == "" {
		c.MCP.Version = "1.0.0"
	}
	if c.Pagination.MaxCharsPerPage <= 0 {
		c.Pagination.MaxCharsPerPage = 50_000
	}
	if c.Pagination.MaxContextLines <= 0 {
		c.Pagination.MaxContextLines = 10
	}
	if c.Pagination.MaxSearchResultsPerPage <= 0 {
		c.Pagination.MaxSearchResultsPerPage = 100
	}
	return nil
}
