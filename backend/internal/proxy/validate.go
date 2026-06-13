package proxy

import (
	"net"
	"regexp"
	"strings"
)

// hostnameRe matches a syntactically valid DNS hostname (one or more labels).
var hostnameRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)

// IsValidPublicHost reports whether host is a syntactically valid, non-internal DNS
// hostname suitable for the proxy allowlist. It rejects IP literals, numeric/encoded
// hosts, and localhost. The SSRF dialer remains the authoritative runtime backstop.
func IsValidPublicHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" || len(host) > 253 || host == "localhost" {
		return false
	}
	if net.ParseIP(host) != nil {
		return false
	}
	if IsNumericOrEncodedHost(host) {
		return false
	}
	return hostnameRe.MatchString(host)
}

// ValidHeaderName reports whether name is a valid RFC 7230 header field-name token.
func ValidHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		if !isTokenChar(name[i]) {
			return false
		}
	}
	return true
}

// ValidHeaderValue reports whether value is free of control characters (CR, LF, NUL,
// etc.) that would enable header injection. Tab and printable bytes are allowed.
func ValidHeaderValue(value string) bool {
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c == '\t' {
			continue
		}
		if c < 0x20 || c == 0x7f {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		return true
	}
	switch c {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false
}
