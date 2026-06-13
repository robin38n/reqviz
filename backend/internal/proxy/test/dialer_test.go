package proxy_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/robin38n/reqviz/backend/internal/proxy"
)

// TestResolveAndValidate covers the SSRF allow/deny decision for literal IPs and
// special hostnames. It uses only IP literals, "localhost", and numeric hosts, so
// it never performs DNS resolution and runs fully offline. Private-range coverage
// here transitively exercises the unexported isPrivateIP helper.
func TestResolveAndValidate(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		wantError bool
	}{
		{"public v4", "8.8.8.8", false},
		{"public v6", "2001:4860:4860::8888", false},
		{"loopback v4", "127.0.0.1", true},
		{"private 10/8", "10.0.0.1", true},
		{"private 172.16/12", "172.16.0.1", true},
		{"private 192.168/16", "192.168.1.1", true},
		{"cgnat 100.64/10", "100.64.0.1", true},
		{"link-local 169.254/16", "169.254.0.1", true},
		{"unspecified", "0.0.0.0", true},
		{"multicast", "224.0.0.1", true},
		{"loopback v6", "::1", true},
		{"link-local v6", "fe80::1", true},
		{"ula v6", "fc00::1", true},
		{"localhost", "localhost", true},
		{"decimal encoded", "2130706433", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := proxy.ResolveAndValidate(context.Background(), tt.host)
			if (err != nil) != tt.wantError {
				t.Errorf("ResolveAndValidate(%q) error = %v, wantError = %v", tt.host, err, tt.wantError)
			}
		})
	}
}

// TestValidateIPs covers the shared private-IP check applied to DNS-resolved hosts.
// Hermetic tests can't trigger real DNS, so this directly exercises the multi-IP
// rejection path (the DNS-rebinding defense) that ResolveAndValidate relies on.
func TestValidateIPs(t *testing.T) {
	tests := []struct {
		name      string
		ips       []net.IP
		wantError bool
	}{
		{"empty", nil, false},
		{"all public", []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("1.1.1.1")}, false},
		{"single private", []net.IP{net.ParseIP("10.0.0.1")}, true},
		{"public then private", []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("192.168.1.1")}, true},
		{"private v6", []net.IP{net.ParseIP("::1")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := proxy.ValidateIPs(tt.ips); (err != nil) != tt.wantError {
				t.Errorf("ValidateIPs(%v) error = %v, wantError = %v", tt.ips, err, tt.wantError)
			}
		})
	}
}

// TestValidateConnAddr directly tests the rebinding-proof SSRF gate that the dialer's
// Control hook runs on every Happy-Eyeballs candidate before connecting.
func TestValidateConnAddr(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		wantError bool
	}{
		{"public v4", "8.8.8.8:443", false},
		{"public v6", "[2001:4860:4860::8888]:443", false},
		{"private v4", "10.0.0.1:80", true},
		{"loopback v4", "127.0.0.1:80", true},
		{"loopback v6", "[::1]:443", true},
		{"non-ip host", "example.com:80", true},
		{"missing port", "8.8.8.8", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := proxy.ValidateConnAddr(tt.address)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateConnAddr(%q) error = %v, wantError = %v", tt.address, err, tt.wantError)
			}
		})
	}
}

// TestSafeCheckRedirect covers the redirect-following SSRF guard: hop cap, scheme
// restriction, and re-validation of the redirect target host.
func TestSafeCheckRedirect(t *testing.T) {
	mkReq := func(rawurl string) *http.Request {
		req, err := http.NewRequest("GET", rawurl, nil)
		if err != nil {
			t.Fatalf("building request for %q: %v", rawurl, err)
		}
		return req
	}

	tests := []struct {
		name      string
		url       string
		viaCount  int
		wantError bool
	}{
		{"public target ok", "https://8.8.8.8/next", 1, false},
		{"too many hops", "https://8.8.8.8/next", 10, true},
		{"non-http scheme", "ftp://8.8.8.8/file", 1, true},
		{"private target", "http://169.254.169.254/latest", 1, true},
		{"loopback target", "http://127.0.0.1/admin", 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			via := make([]*http.Request, tt.viaCount)
			err := proxy.SafeCheckRedirect(mkReq(tt.url), via)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeCheckRedirect(%q, via=%d) error = %v, wantError = %v",
					tt.url, tt.viaCount, err, tt.wantError)
			}
		})
	}
}

// TestSafeClient_BlocksLoopback is an end-to-end check that the safe client's dialer
// refuses to connect to a loopback httptest server (offline; no external network).
func TestSafeClient_BlocksLoopback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := proxy.NewSafeClient()
	resp, err := client.Get(ts.URL)
	if err == nil {
		resp.Body.Close()
		t.Fatalf("expected safe client to block loopback %s, got success", ts.URL)
	}
}

// TestHappyEyeballs_LiveDualStack confirms the safe dialer connects to a real
// dual-stack host via Happy Eyeballs fast-fallback. It is gated behind REQVIZ_NET_TEST
// so the default suite stays hermetic and offline-safe. Run with `task test:net`.
func TestHappyEyeballs_LiveDualStack(t *testing.T) {
	if os.Getenv("REQVIZ_NET_TEST") == "" {
		t.Skip("set REQVIZ_NET_TEST=1 to run live network tests (e.g. `task test:net`)")
	}

	url := os.Getenv("REQVIZ_NET_TEST_URL")
	if url == "" {
		url = "https://pokeapi.co/api/v2/pokemon/ditto"
	}

	client := proxy.NewSafeClient()
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("live dual-stack GET %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET %s status = %d, want 200", url, resp.StatusCode)
	}
	// Happy Eyeballs falls back within ~300ms even when IPv6 is unreachable; a 5s
	// budget proves we are not stalling on a dead IPv6 connect (which took ~10s before).
	if elapsed > 5*time.Second {
		t.Errorf("GET %s took %v, want < 5s (Happy Eyeballs fast-fallback)", url, elapsed)
	}
	t.Logf("live dual-stack GET %s -> %d in %v", url, resp.StatusCode, elapsed)
}
