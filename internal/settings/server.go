// Package settings serves TaskGlass' graphical settings window.
//
// The window is a small web page (embedded in the executable) served on a
// loopback-only port and shown in an app-style browser window, so it needs no
// GUI toolkit, no extra DLLs and no CGO. Every request must carry a random
// per-run token and a loopback Host header, so neither other local programs
// nor web pages (DNS rebinding) can change the settings.
package settings

import (
	"crypto/rand"
	"crypto/subtle"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/i18n"
)

//go:embed index.html
var indexHTML string

// Host is what the settings window needs from the running application.
type Host interface {
	Config() config.Config
	// Apply validates, saves and applies a whole configuration.
	Apply(config.Config) error
	AutostartEnabled() bool
	SetAutostart(on bool) error
	OpenFile(which string) error // "config" or "log"
	Version() string
	ConfigPath() string
	XamlAvailable() bool
}

// Server is the loopback HTTP server behind the settings window.
type Server struct {
	host  Host
	token string

	mu sync.Mutex
	ln net.Listener
	hs *http.Server
}

// New creates a Server. It does not listen until URL is first called.
func New(h Host) *Server {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return &Server{host: h, token: hex.EncodeToString(b[:])}
}

// URL starts the server if needed and returns the address of the settings page.
func (s *Server) URL() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln == nil {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return "", err
		}
		s.ln = ln
		s.hs = &http.Server{Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second}
		go s.hs.Serve(ln)
	}
	return "http://" + s.ln.Addr().String() + "/?t=" + s.token, nil
}

// Close stops the server.
func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.hs != nil {
		_ = s.hs.Close()
		s.hs, s.ln = nil, nil
	}
}

// Handler returns the HTTP handler (exported for tests).
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /api/state", s.guard(s.state))
	mux.HandleFunc("POST /api/config", s.guard(s.saveConfig))
	mux.HandleFunc("POST /api/defaults", s.guard(s.defaults))
	mux.HandleFunc("POST /api/autostart", s.guard(s.autostart))
	mux.HandleFunc("POST /api/open", s.guard(s.open))
	return mux
}

func loopbackHost(h string) bool {
	host, _, err := net.SplitHostPort(h)
	if err != nil {
		host = h
	}
	return host == "127.0.0.1" || host == "localhost" || host == "[::1]" || host == "::1"
}

func (s *Server) tokenOK(t string) bool {
	return subtle.ConstantTimeCompare([]byte(t), []byte(s.token)) == 1
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if !loopbackHost(r.Host) || !s.tokenOK(r.URL.Query().Get("t")) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; connect-src 'self'; img-src data:")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.WriteString(w, strings.Replace(indexHTML, "__TOKEN__", s.token, 1))
}

func (s *Server) guard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !loopbackHost(r.Host) || !s.tokenOK(r.Header.Get("X-TaskGlass-Token")) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			if o := r.Header.Get("Origin"); o != "" && !loopbackHost(strings.TrimPrefix(strings.TrimPrefix(o, "http://"), "https://")) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		}
		w.Header().Set("Cache-Control", "no-store")
		next(w, r)
	}
}

type stateResp struct {
	Config        config.Config                   `json:"config"`
	Defaults      config.Config                   `json:"defaults"`
	Autostart     bool                            `json:"autostart"`
	Version       string                          `json:"version"`
	ConfigPath    string                          `json:"configPath"`
	XamlAvailable bool                            `json:"xamlAvailable"`
	SystemLang    i18n.Lang                       `json:"systemLang"`
	I18n          map[i18n.Lang]map[string]string `json:"i18n"`
	MaxApps       int                             `json:"maxApps"`
}

func (s *Server) state(w http.ResponseWriter, _ *http.Request) {
	tbl := map[i18n.Lang]map[string]string{}
	for _, l := range i18n.Langs {
		tbl[l] = i18n.Table(l)
	}
	writeJSON(w, http.StatusOK, stateResp{
		Config:        s.host.Config(),
		Defaults:      config.Default(),
		Autostart:     s.host.AutostartEnabled(),
		Version:       s.host.Version(),
		ConfigPath:    s.host.ConfigPath(),
		XamlAvailable: s.host.XamlAvailable(),
		SystemLang:    i18n.Resolve(i18n.Auto, systemLocale()),
		I18n:          tbl,
		MaxApps:       config.MaxAppRules,
	})
}

func (s *Server) saveConfig(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	c, err := config.Parse(body) // strict: unknown fields and bad values are rejected
	if err != nil {
		writeErr(w, http.StatusUnprocessableEntity, err)
		return
	}
	if err := s.host.Apply(c); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "config": s.host.Config()})
}

func (s *Server) defaults(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, config.Default())
}

func (s *Server) autostart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		On bool `json:"on"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := s.host.SetAutostart(req.On); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "autostart": s.host.AutostartEnabled()})
}

func (s *Server) open(w http.ResponseWriter, r *http.Request) {
	var req struct {
		What string `json:"what"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.What != "config" && req.What != "log") {
		writeErr(w, http.StatusBadRequest, errors.New("what must be config or log"))
		return
	}
	if err := s.host.OpenFile(req.What); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]any{"ok": false, "error": err.Error()})
}
