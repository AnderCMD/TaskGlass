package settings

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

type fakeHost struct {
	cfg  config.Config
	auto bool
}

func (f *fakeHost) Config() config.Config { return f.cfg }
func (f *fakeHost) Apply(c config.Config) error {
	f.cfg = c
	return nil
}
func (f *fakeHost) AutostartEnabled() bool     { return f.auto }
func (f *fakeHost) SetAutostart(on bool) error { f.auto = on; return nil }
func (f *fakeHost) OpenFile(string) error      { return nil }
func (f *fakeHost) Version() string            { return "test" }
func (f *fakeHost) ConfigPath() string         { return "cfg.json" }
func (f *fakeHost) XamlAvailable() bool        { return true }

func do(h http.Handler, method, path, host, token, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Host = host
	if token != "" {
		r.Header.Set("X-TaskGlass-Token", token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestAuthAndHostChecks(t *testing.T) {
	f := &fakeHost{cfg: config.Default()}
	s := New(f)
	h := s.Handler()

	if w := do(h, "GET", "/api/state", "127.0.0.1:1", "", ""); w.Code != 403 {
		t.Errorf("no token: %d", w.Code)
	}
	if w := do(h, "GET", "/api/state", "127.0.0.1:1", "wrong", ""); w.Code != 403 {
		t.Errorf("bad token: %d", w.Code)
	}
	if w := do(h, "GET", "/api/state", "evil.example.com", s.token, ""); w.Code != 403 {
		t.Errorf("rebinding host: %d", w.Code)
	}
	if w := do(h, "GET", "/?t=wrong", "127.0.0.1:1", "", ""); w.Code != 403 {
		t.Errorf("page with bad token: %d", w.Code)
	}
	w := do(h, "GET", "/?t="+s.token, "127.0.0.1:1", "", "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), s.token) || strings.Contains(w.Body.String(), "__TOKEN__") {
		t.Errorf("page not served with token: %d", w.Code)
	}
}

func TestStateAndSave(t *testing.T) {
	f := &fakeHost{cfg: config.Default()}
	s := New(f)
	h := s.Handler()

	w := do(h, "GET", "/api/state", "127.0.0.1:1", s.token, "")
	if w.Code != 200 {
		t.Fatalf("state: %d %s", w.Code, w.Body)
	}
	var st struct {
		I18n map[string]map[string]string `json:"i18n"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil || st.I18n["es"]["ui.save"] != "Guardar" || st.I18n["en"]["ui.save"] != "Save" {
		t.Fatalf("i18n tables missing: %v", err)
	}

	c := config.Default()
	c.Default.Accent = config.AccentAcrylic
	c.Default.Opacity = 55
	c.Language = "es"
	body, _ := json.Marshal(c)
	if w := do(h, "POST", "/api/config", "127.0.0.1:1", s.token, string(body)); w.Code != 200 {
		t.Fatalf("save: %d %s", w.Code, w.Body)
	}
	if f.cfg.Default.Accent != config.AccentAcrylic || f.cfg.Language != "es" {
		t.Errorf("config not applied: %+v", f.cfg.Default)
	}

	// Invalid values and unknown keys are rejected without touching the config.
	bad := strings.Replace(string(body), `"acrylic"`, `"rainbow"`, 1)
	if w := do(h, "POST", "/api/config", "127.0.0.1:1", s.token, bad); w.Code != 422 {
		t.Errorf("invalid accent accepted: %d", w.Code)
	}
	if w := do(h, "POST", "/api/config", "127.0.0.1:1", s.token, `{"bogus":1}`); w.Code != 422 {
		t.Errorf("unknown key accepted: %d", w.Code)
	}
	if f.cfg.Default.Accent != config.AccentAcrylic {
		t.Error("rejected save changed the config")
	}

	if w := do(h, "POST", "/api/autostart", "127.0.0.1:1", s.token, `{"on":true}`); w.Code != 200 || !f.auto {
		t.Errorf("autostart: %d", w.Code)
	}
	if w := do(h, "POST", "/api/open", "127.0.0.1:1", s.token, `{"what":"/etc/passwd"}`); w.Code != 400 {
		t.Errorf("open accepted arbitrary path: %d", w.Code)
	}
}
