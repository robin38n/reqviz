package proxy_test

import (
	"bytes"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/robin38n/reqviz/backend/internal/proxy"
)

func TestSanitizeHeaders(t *testing.T) {
	h := http.Header{
		"Content-Type":  []string{"application/json"},
		"Set-Cookie":    []string{"bad=cookie"},
		"Date":          []string{"today"},
		"Authorization": []string{"secret"},
	}

	sanitized := proxy.SanitizeHeaders(h)

	if _, ok := sanitized["Content-Type"]; !ok {
		t.Error("expected Content-Type to be preserved")
	}
	if _, ok := sanitized["Date"]; !ok {
		t.Error("expected Date to be preserved")
	}
	if _, ok := sanitized["Set-Cookie"]; ok {
		t.Error("expected Set-Cookie to be stripped")
	}
	if _, ok := sanitized["Authorization"]; ok {
		t.Error("expected Authorization to be stripped")
	}
}

func TestIsNumericOrEncodedHost(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"example.com", false},
		{"127.0.0.1", false}, // IPs are not numeric-re-matched (handled by net.ParseIP)
		{"2130706433", true}, // Decimal
		{"0x7f000001", true}, // Hex
		{"017700000001", true}, // Octal
		{"[::1]", false}, // IPv6
	}

	for _, tt := range tests {
		if got := proxy.IsNumericOrEncodedHost(tt.host); got != tt.want {
			t.Errorf("IsNumericOrEncodedHost(%q) = %v, want %v", tt.host, got, tt.want)
		}
	}
}

func TestHostInList(t *testing.T) {
	list := []string{"example.com", "API.foo.com"}

	if !proxy.HostInList("example.com", list) {
		t.Error("expected example.com to be in list")
	}
	if !proxy.HostInList("EXAMPLE.COM", list) {
		t.Error("expected host matching to be case-insensitive")
	}
	if !proxy.HostInList("api.foo.com", list) {
		t.Error("expected host matching to be case-insensitive for list items too")
	}
	if proxy.HostInList("google.com", list) {
		t.Error("expected google.com to NOT be in list")
	}
}

func TestContentTypeAllowed(t *testing.T) {
	tests := []struct {
		name string
		ct   string
		want bool
	}{
		{"json", "application/json", true},
		{"json with charset", "application/json; charset=utf-8", true},
		{"text plain with charset", "text/plain; charset=utf-8", true},
		{"csv", "text/csv", true},
		{"xml application", "application/xml", true},
		{"form urlencoded", "application/x-www-form-urlencoded", true},
		{"empty allowed", "", true},
		{"uppercase normalized", "APPLICATION/JSON", true},
		{"html blocked", "text/html", false},
		{"png blocked", "image/png", false},
		{"octet-stream blocked", "application/octet-stream", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := proxy.ContentTypeAllowed(tt.ct); got != tt.want {
				t.Errorf("ContentTypeAllowed(%q) = %v, want %v", tt.ct, got, tt.want)
			}
		})
	}
}

func TestExtractServerHosts(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]any
		want []string
	}{
		{
			name: "dedupe and case-fold across servers",
			raw: map[string]any{"servers": []any{
				map[string]any{"url": "https://API.example.com/v1"},
				map[string]any{"url": "https://api.example.com/v2"},
				map[string]any{"url": "https://other.com"},
			}},
			want: []string{"api.example.com", "other.com"},
		},
		{
			name: "missing servers key",
			raw:  map[string]any{},
			want: nil,
		},
		{
			name: "malformed and empty urls skipped",
			raw: map[string]any{"servers": []any{
				map[string]any{"url": "://nope"},
				map[string]any{"url": ""},
				map[string]any{"nope": "x"},
				map[string]any{"url": "https://good.com"},
			}},
			want: []string{"good.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := proxy.ExtractServerHosts(tt.raw)
			sort.Strings(got)
			want := append([]string(nil), tt.want...)
			sort.Strings(want)
			if strings.Join(got, ",") != strings.Join(want, ",") {
				t.Errorf("ExtractServerHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBodyWithCap(t *testing.T) {
	t.Run("under cap returns full body", func(t *testing.T) {
		data := []byte("hello world")
		buf, truncated, err := proxy.ReadBodyWithCap(bytes.NewReader(data), 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if truncated {
			t.Error("expected truncated=false")
		}
		if !bytes.Equal(buf, data) {
			t.Errorf("body = %q, want %q", buf, data)
		}
	})

	t.Run("over cap truncates and flags", func(t *testing.T) {
		data := bytes.Repeat([]byte("x"), 200)
		buf, truncated, err := proxy.ReadBodyWithCap(bytes.NewReader(data), 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !truncated {
			t.Error("expected truncated=true")
		}
		if len(buf) != 100 {
			t.Errorf("len(body) = %d, want 100", len(buf))
		}
	})
}
