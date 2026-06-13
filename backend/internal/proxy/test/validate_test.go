package proxy_test

import (
	"strings"
	"testing"

	"github.com/robin38n/reqviz/backend/internal/proxy"
)

func TestIsValidPublicHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{"simple domain", "example.com", true},
		{"subdomains", "api.foo.co.uk", true},
		{"trailing whitespace trimmed", "  example.com  ", true},
		{"uppercase normalized", "Example.COM", true},
		{"empty", "", false},
		{"localhost", "localhost", false},
		{"too long", strings.Repeat("a", 254), false},
		{"ipv4 literal", "127.0.0.1", false},
		{"ipv6 literal", "[::1]", false},
		{"decimal encoded", "2130706433", false},
		{"hex encoded", "0x7f000001", false},
		{"leading dash", "-bad.com", false},
		{"underscore", "bad_host", false},
		{"empty label", "a..b", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := proxy.IsValidPublicHost(tt.host); got != tt.want {
				t.Errorf("IsValidPublicHost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestValidHeaderName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"simple token", "X-Test", true},
		{"token chars", "X-Custom_Header.1", true},
		{"empty", "", false},
		{"space", "X Test", false},
		{"colon", "X:Test", false},
		{"newline", "X\nTest", false},
		{"control char", "X\x00Test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := proxy.ValidHeaderName(tt.input); got != tt.want {
				t.Errorf("ValidHeaderName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidHeaderValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"printable ascii", "application/json", true},
		{"tab allowed", "a\tb", true},
		{"empty allowed", "", true},
		{"carriage return", "a\rb", false},
		{"line feed", "a\nb", false},
		{"crlf injection", "value\r\nX-Injected: 1", false},
		{"null byte", "a\x00b", false},
		{"del char", "a\x7fb", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := proxy.ValidHeaderValue(tt.input); got != tt.want {
				t.Errorf("ValidHeaderValue(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
