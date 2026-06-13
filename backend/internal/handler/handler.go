package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/robin38n/reqviz/backend/internal/parser"
	"github.com/robin38n/reqviz/backend/internal/proxy"
	"github.com/robin38n/reqviz/backend/internal/store"
)

type Server struct {
	store *store.SpecStore
	proxy *proxy.Executor
}

func NewServer(s *store.SpecStore) *Server {
	return &Server{
		store: s,
		proxy: proxy.New(proxy.NewSafeClient(), proxy.NewLimiter(), slog.Default()),
	}
}

func (s *Server) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// UploadSpec validates, stores, and returns a SpecSummary for an uploaded OpenAPI spec.
func (s *Server) UploadSpec(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20)) // 10 MB limit
	if err != nil {
		writeError(w, http.StatusBadRequest, "Could not read the request body.")
		return
	}

	format := detectFormat(body)
	var result *parser.ParseResult
	if format == "JSON" {
		result, err = parser.FromJSON(body)
	} else {
		result, err = parser.FromYAML(body)
	}

	if err != nil {
		writeError(w, http.StatusBadRequest, "The spec couldn't be parsed. Make sure it's valid JSON or YAML.")
		return
	}

	writeJSON(w, http.StatusCreated, s.storeResult(result))
}

func (s *Server) GetSpec(w http.ResponseWriter, _ *http.Request, id openapi_types.UUID) {
	stored, err := s.store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, NotFound{strPtr("Spec not found. It may have expired, or the server was restarted.")})
		return
	}

	writeJSON(w, http.StatusOK, ParsedSpec{
		Id:  openapi_types.UUID(stored.ID),
		Raw: stored.Raw,
		Summary: SpecSummary{
			Id:            openapi_types.UUID(stored.ID),
			Title:         stored.Title,
			Version:       stored.Version,
			EndpointCount: stored.EndpointCount,
			SchemaCount:   stored.SchemaCount,
			Tags:          &stored.Tags,
			CreatedAt:     &stored.CreatedAt,
			Approved:      stored.Approved,
			AllowedHosts:  stored.AllowedHosts,
		},
	})
}

// ApproveSpec marks a spec as approved and optionally updates its host allowlist.
func (s *Server) ApproveSpec(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	stored, err := s.store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, NotFound{strPtr("Spec not found. It may have expired, or the server was restarted.")})
		return
	}

	hosts := stored.AllowedHosts
	if r.Body != nil {
		var req ApproveSpecRequest
		_ = json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&req)
		if req.AllowedHosts != nil && len(*req.AllowedHosts) > 0 {
			for _, h := range *req.AllowedHosts {
				if !proxy.IsValidPublicHost(h) {
					writeError(w, http.StatusBadRequest, "One of the listed hosts is not a valid public hostname.")
					return
				}
			}
			hosts = *req.AllowedHosts
		}
	}

	updated, err := s.store.Approve(id, hosts)
	if err != nil {
		writeJSON(w, http.StatusNotFound, NotFound{strPtr("Spec not found. It may have expired, or the server was restarted.")})
		return
	}

	writeJSON(w, http.StatusOK, SpecSummary{
		Id:            openapi_types.UUID(updated.ID),
		Title:         updated.Title,
		Version:       updated.Version,
		EndpointCount: updated.EndpointCount,
		SchemaCount:   updated.SchemaCount,
		Tags:          &updated.Tags,
		CreatedAt:     &updated.CreatedAt,
		Approved:      updated.Approved,
		AllowedHosts:  updated.AllowedHosts,
	})
}

var proxyMiddleware = Chain(OriginAllowed)

// ProxyRequest gates and forwards requests via the proxy executor using inline middleware.
func (s *Server) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	proxyMiddleware(http.HandlerFunc(s.proxyRequestCore)).ServeHTTP(w, r)
}

func (s *Server) proxyRequestCore(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req ProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Default().Error("proxy: decode request body", "error", err)
		writeError(w, http.StatusBadRequest, "The request body is not valid JSON.")
		return
	}

	if !req.Method.Valid() {
		writeError(w, http.StatusBadRequest, "That HTTP method isn't supported.")
		return
	}

	if len(req.Url) > 2048 {
		writeError(w, http.StatusBadRequest, "The URL is too long (max 2048 characters).")
		return
	}
	parsed, err := url.Parse(req.Url)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		writeError(w, http.StatusBadRequest, "Only http:// and https:// URLs are allowed.")
		return
	}
	if parsed.User != nil {
		writeError(w, http.StatusBadRequest, "URLs with embedded credentials (user:pass@) aren't allowed.")
		return
	}

	stored, err := s.store.Get(req.SpecId)
	if err != nil {
		writeError(w, http.StatusForbidden, "No approved spec is loaded for this request. Open a spec in the Explorer and approve it first.")
		return
	}
	if !stored.Approved {
		writeError(w, http.StatusForbidden, "This spec hasn't been approved yet. Approve it before sending requests.")
		return
	}

	for name, value := range derefHeaders(req.Headers) {
		if !proxy.ValidHeaderName(name) || !proxy.ValidHeaderValue(value) {
			writeError(w, http.StatusBadRequest, "A request header has an invalid name or value.")
			return
		}
	}

	res, err := s.proxy.Execute(r.Context(), proxy.Input{
		Method:       string(req.Method),
		URL:          req.Url,
		Headers:      derefHeaders(req.Headers),
		Body:         req.Body,
		SpecID:       req.SpecId.String(),
		Origin:       r.Header.Get("Origin"),
		AllowedHosts: stored.AllowedHosts,
	})

	if err != nil {
		switch {
		case errors.Is(err, proxy.ErrHostNotAllowed):
			writeError(w, http.StatusForbidden, "This host isn't in the spec's approved list. You can only call the hosts you approved for this spec.")
		case errors.Is(err, proxy.ErrRateLimited):
			writeError(w, http.StatusTooManyRequests, "Too many requests to this host. Wait a few seconds and try again.")
		case errors.Is(err, proxy.ErrSSRFBlocked):
			// Log the real reason; never reveal the SSRF check internals to the client.
			slog.Default().Warn("proxy: request blocked", "error", err)
			writeError(w, http.StatusForbidden, "That address can't be reached. ReqViz only allows requests to public internet hosts, not local or private networks.")
		case errors.Is(err, proxy.ErrBadRequest):
			slog.Default().Warn("proxy: bad request", "error", err)
			writeError(w, http.StatusBadRequest, "The request couldn't be built. Check the URL, method, and body.")
		default:
			slog.Default().Error("proxy: internal error", "error", err)
			writeError(w, http.StatusInternalServerError, "Something went wrong while sending the request.")
		}
		return
	}

	writeJSON(w, http.StatusOK, ProxyResponse{
		Status:     res.Status,
		Headers:    res.Headers,
		Body:       res.Body,
		DurationMs: &res.DurationMs,
	})
}

func (s *Server) parseAndStoreSpec(raw map[string]any) (*SpecSummary, error) {
	rawBytes, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to process spec")
	}

	result, err := parser.FromJSON(rawBytes)
	if err != nil {
		return nil, err
	}

	return s.storeResult(result), nil
}

func (s *Server) storeResult(r *parser.ParseResult) *SpecSummary {
	stored := &store.StoredSpec{
		Title:         r.Title,
		Version:       r.Version,
		EndpointCount: r.EndpointCount,
		SchemaCount:   r.SchemaCount,
		Tags:          r.Tags,
		Raw:           r.Raw,
		AllowedHosts:  proxy.ExtractServerHosts(r.Raw),
	}
	id := s.store.Save(stored)

	now := time.Now()
	return &SpecSummary{
		Id:            openapi_types.UUID(id),
		Title:         r.Title,
		Version:       r.Version,
		EndpointCount: r.EndpointCount,
		SchemaCount:   r.SchemaCount,
		Tags:          &r.Tags,
		CreatedAt:     &now,
		Approved:      stored.Approved,
		AllowedHosts:  stored.AllowedHosts,
	}
}

// detectFormat guesses JSON/YAML by checking if content starts with '{'.
func detectFormat(body []byte) string {
	for _, b := range body {
		switch b {
		case ' ', '\t', '\n', '\r':
			continue
		case '{':
			return "JSON"
		default:
			return "YAML"
		}
	}
	return "YAML"
}

func derefHeaders(h *map[string]string) map[string]string {
	if h == nil {
		return nil
	}
	return *h
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ValidationError{Error: &msg})
}

func strPtr(s string) *string {
	return &s
}
