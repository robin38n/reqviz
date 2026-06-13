package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

var privateRanges []*net.IPNet

func init() {
	cidrs := []string{
		"127.0.0.0/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
		"100.64.0.0/10", "0.0.0.0/8", "169.254.0.0/16", "224.0.0.0/4", "255.255.255.255/32",
		"::1/128", "fe80::/10", "fc00::/7",
	}
	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("invalid CIDR: " + cidr)
		}
		privateRanges = append(privateRanges, network)
	}
}

func isPrivateIP(ip net.IP) bool {
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	for _, network := range privateRanges {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// ResolveAndValidate resolves a hostname and ensures no returned IPs are private.
func ResolveAndValidate(ctx context.Context, host string) ([]net.IP, error) {
	if host == "localhost" {
		return nil, fmt.Errorf("requests to localhost are not allowed")
	}
	if IsNumericOrEncodedHost(host) {
		return nil, fmt.Errorf("numeric/encoded hostnames are not allowed")
	}

	if ip := net.ParseIP(host); ip != nil {
		if isPrivateIP(ip) {
			return nil, fmt.Errorf("requests to private/internal addresses are not allowed")
		}
		return []net.IP{ip}, nil
	}

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, fmt.Errorf("DNS resolution failed: %w", err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("DNS resolution returned no addresses")
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return nil, fmt.Errorf("requests to private/internal addresses are not allowed")
		}
	}
	return ips, nil
}

func safeDialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address %q: %w", addr, err)
	}
	// Early, explicit rejection of localhost, numeric/encoded hosts, and literal
	// private IPs (fail-closed if any resolved IP is private).
	if _, err := ResolveAndValidate(ctx, host); err != nil {
		return nil, err
	}
	// Dial the hostname so the stdlib runs Happy Eyeballs (RFC 6555/8305): it races
	// IPv4/IPv6 and falls back within FallbackDelay (~300ms) if a family is unreachable.
	// Control validates the actual IP the kernel is about to connect to — on every
	// candidate — so a private/rebound address is never connected to.
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
		Control: func(_, address string, _ syscall.RawConn) error {
			h, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}
			ip := net.ParseIP(h)
			if ip == nil || isPrivateIP(ip) {
				return fmt.Errorf("connection to disallowed address blocked: %s", address)
			}
			return nil
		},
	}
	return dialer.DialContext(ctx, network, addr)
}

func safeCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("too many redirects (max 10)")
	}
	scheme := req.URL.Scheme
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("redirect to disallowed scheme: %s", scheme)
	}
	host := req.URL.Hostname()
	if _, err := ResolveAndValidate(req.Context(), host); err != nil {
		return fmt.Errorf("redirect blocked: %w", err)
	}
	return nil
}
