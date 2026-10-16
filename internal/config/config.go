// Package config loads x360-notice configuration from the environment.
package config

import "strings"

type Config struct {
	Addr       string // X360_NOTICE_ADDR, default :8080
	StatusText string // X360_NOTICE_STATUS_TEXT, default "Xbox LIVE is up and running"
}

// DefaultStatusText is the service-status message the dashboard shows.
const DefaultStatusText = "Xbox LIVE is up and running"

// Load reads the environment through getenv (os.Getenv in production).
func Load(getenv func(string) string) (*Config, error) {
	c := &Config{Addr: getenv("X360_NOTICE_ADDR"), StatusText: strings.TrimSpace(getenv("X360_NOTICE_STATUS_TEXT"))}
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.StatusText == "" {
		c.StatusText = DefaultStatusText
	}
	return c, nil
}
