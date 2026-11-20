package update

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewer(t *testing.T) {
	cases := []struct {
		cur, lat string
		want     bool
	}{
		{"1.3.1", "1.3.2", true}, {"1.3.1", "1.3.1", false}, {"1.3.1", "1.3.0", false},
		{"1.9.0", "1.10.0", true}, {"1.3.1", "2.0.0", true}, {"v1.3.1", "v1.4.0", true},
		{"dev", "9.9.9", false}, {"", "1.0.0", false}, {"1.3.1-ci", "1.3.2", true},
		{"1.3.1 (abc1234)", "1.3.2", true}, {"1.3.1", "banana", false},
	}
	for _, c := range cases {
		if got := Newer(c.cur, c.lat); got != c.want {
			t.Errorf("Newer(%q,%q)=%v want %v", c.cur, c.lat, got, c.want)
		}
	}
}

func TestParseSums(t *testing.T) {
	h := strings.Repeat("ab", 32)
	m := parseSums(h + " *TaskGlass-1.0.0-Setup.exe\n" + strings.Repeat("cd", 32) + "  other.zip\nnot a line\n")
	if m["TaskGlass-1.0.0-Setup.exe"] != h || len(m) != 2 {
		t.Fatalf("got %v", m)
	}
}

// fake builds a server that looks like GitHub for one release.
func fake(t *testing.T, installer []byte, sumOverride string) (*Client, *httptest.Server) {
	t.Helper()
	sum := sha256.Sum256(installer)
	if sumOverride != "" {
		sum = [32]byte{}
		copy(sum[:], sumOverride)
	}
	mux := http.NewServeMux()
	var srv *httptest.Server
	mux.HandleFunc("/repos/AnderCMD/TaskGlass/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":"v1.4.0","body":"notes","assets":[
			{"name":"TaskGlass-1.4.0-Setup.exe","browser_download_url":"%[1]s/dl/TaskGlass-1.4.0-Setup.exe"},
			{"name":"SHA256SUMS.txt","browser_download_url":"%[1]s/dl/SHA256SUMS.txt"},
			{"name":"taskglass-1.4.0-windows-amd64.zip","browser_download_url":"%[1]s/dl/x.zip"}]}`, srv.URL)
	})
	mux.HandleFunc("/dl/TaskGlass-1.4.0-Setup.exe", func(w http.ResponseWriter, r *http.Request) { w.Write(installer) })
	mux.HandleFunc("/dl/SHA256SUMS.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s  TaskGlass-1.4.0-Setup.exe\n", hex.EncodeToString(sum[:]))
	})
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	c := New("TaskGlass/test")
	c.APIBase = srv.URL
	c.allow = func(u *url.URL) bool { return strings.HasPrefix(u.String(), srv.URL) }
	return c, srv
}

func fakeInstaller() []byte { return bytes.Repeat([]byte("TASKGLASS"), 200_000) } // ~1.8 MB

func TestLatestAndDownload(t *testing.T) {
	inst := fakeInstaller()
	c, _ := fake(t, inst, "")
	r, err := c.Latest(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if r.Version != "1.4.0" || r.InstallerName != "TaskGlass-1.4.0-Setup.exe" || r.Notes != "notes" {
		t.Fatalf("release: %+v", r)
	}
	dir := t.TempDir()
	path, err := c.Download(context.Background(), r, dir)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if !bytes.Equal(got, inst) || filepath.Base(path) != r.InstallerName {
		t.Fatal("downloaded file differs")
	}
	if _, err := os.Stat(path + ".part"); !os.IsNotExist(err) {
		t.Error("partial file left behind")
	}
}

func TestDownloadRejectsChecksumMismatch(t *testing.T) {
	c, _ := fake(t, fakeInstaller(), "wrong-checksum-wrong-checksum-12")
	r, err := c.Latest(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if _, err := c.Download(context.Background(), r, dir); err == nil {
		t.Fatal("a mismatching installer must be rejected")
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("rejected download left files: %v", entries)
	}
}

func TestDownloadRejectsOtherHostsAndTinyFiles(t *testing.T) {
	c := New("TaskGlass/test") // strict GitHub-only policy
	r := Release{InstallerName: "TaskGlass-1.4.0-Setup.exe", InstallerURL: "https://evil.example.com/x.exe", SumsURL: "https://github.com.evil.example.com/s"}
	if _, err := c.Download(context.Background(), r, t.TempDir()); err == nil {
		t.Fatal("non-GitHub hosts must be refused")
	}
	if githubOnly(&url.URL{Scheme: "http", Host: "github.com"}) {
		t.Error("plain http must be refused")
	}
	if !githubOnly(&url.URL{Scheme: "https", Host: "objects.githubusercontent.com"}) || !githubOnly(&url.URL{Scheme: "https", Host: "api.github.com"}) {
		t.Error("GitHub hosts must be accepted")
	}
	if githubOnly(&url.URL{Scheme: "https", Host: "notgithub.com"}) || githubOnly(&url.URL{Scheme: "https", Host: "github.com.evil.com"}) {
		t.Error("look-alike hosts must be refused")
	}
	small, _ := fake(t, []byte("tiny"), "")
	rel, _ := small.Latest(context.Background())
	if _, err := small.Download(context.Background(), rel, t.TempDir()); err == nil {
		t.Fatal("an implausibly small installer must be rejected")
	}
}
