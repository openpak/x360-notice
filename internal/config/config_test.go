package config

import "testing"

func env(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestDefaults(t *testing.T) {
	c, err := Load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if c.Addr != ":8080" || c.StatusText != DefaultStatusText {
		t.Fatalf("defaults: %+v", c)
	}
}

func TestOverrides(t *testing.T) {
	c, _ := Load(env(map[string]string{"X360_NOTICE_ADDR": "127.0.0.1:18080", "X360_NOTICE_STATUS_TEXT": " Maintenance "}))
	if c.Addr != "127.0.0.1:18080" || c.StatusText != "Maintenance" {
		t.Fatalf("overrides: %+v", c)
	}
}
