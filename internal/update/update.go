// Package update finds new TaskGlass releases on GitHub and downloads their
// installer, verifying it against the release's SHA256SUMS.txt before it is
// ever run. It is plain net/http code with no OS dependencies.
//
// Trust model: the connection is HTTPS to GitHub only (redirects may only go to
// GitHub hosts), and the installer must match the checksum published in the
// same release. Without code signing that is integrity against corruption and
// network tampering, not against a compromised GitHub account.
package update

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Repo is the GitHub repository releases are read from.
const Repo = "AnderCMD/TaskGlass"

const (
	apiBase        = "https://api.github.com"
	maxInstaller   = 80 << 20
	maxSums        = 64 << 10
	maxAPIResponse = 2 << 20
	minInstaller   = 1 << 20 // anything smaller is certainly not our installer
)

var (
	installerName = regexp.MustCompile(`^TaskGlass-[0-9][0-9A-Za-z.\-]*-Setup\.exe$`)
	versionRe     = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
)

// Release is a published TaskGlass release.
type Release struct {
	Version       string // without the leading "v", for example "1.3.2"
	Notes         string
	InstallerName string
	InstallerURL  string
	SumsURL       string
}

// Client talks to GitHub.
type Client struct {
	HTTP      *http.Client
	APIBase   string
	Repo      string
	UserAgent string

	// allow decides which URLs may be fetched; tests replace it.
	allow func(*url.URL) bool
}

// New returns a Client for the TaskGlass repository. agent identifies the
// program version only, for example "TaskGlass/1.3.1".
func New(agent string) *Client {
	c := &Client{APIBase: apiBase, Repo: Repo, UserAgent: agent, allow: githubOnly}
	c.HTTP = &http.Client{
		Timeout: 10 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 5 {
				return errors.New("too many redirects")
			}
			if !c.allow(req.URL) {
				return fmt.Errorf("redirect to %s is not allowed", req.URL.Host)
			}
			return nil
		},
	}
	return c
}

// githubOnly accepts https URLs on GitHub's own hosts.
func githubOnly(u *url.URL) bool {
	if u.Scheme != "https" {
		return false
	}
	h := strings.ToLower(u.Hostname())
	return h == "github.com" || h == "api.github.com" || strings.HasSuffix(h, ".githubusercontent.com")
}

func (c *Client) get(ctx context.Context, rawURL, accept string, limit int64) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil || !c.allow(u) {
		return nil, fmt.Errorf("url %q is not allowed", rawURL)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", u.Host, resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, errors.New("response too large")
	}
	return data, nil
}

// Latest returns the newest published release.
func (c *Client) Latest(ctx context.Context) (Release, error) {
	data, err := c.get(ctx, fmt.Sprintf("%s/repos/%s/releases/latest", c.APIBase, c.Repo), "application/vnd.github+json", maxAPIResponse)
	if err != nil {
		return Release{}, err
	}
	var doc struct {
		Tag        string `json:"tag_name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		Assets     []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return Release{}, fmt.Errorf("parse release: %w", err)
	}
	if doc.Draft || doc.Prerelease {
		return Release{}, errors.New("latest release is not a stable release")
	}
	r := Release{Version: strings.TrimPrefix(doc.Tag, "v"), Notes: doc.Body}
	if !versionRe.MatchString(r.Version) {
		return Release{}, fmt.Errorf("unexpected release tag %q", doc.Tag)
	}
	for _, a := range doc.Assets {
		switch {
		case installerName.MatchString(a.Name):
			r.InstallerName, r.InstallerURL = a.Name, a.URL
		case a.Name == "SHA256SUMS.txt":
			r.SumsURL = a.URL
		}
	}
	if r.InstallerURL == "" || r.SumsURL == "" {
		return Release{}, errors.New("release has no installer or checksum file")
	}
	return r, nil
}

// Newer reports whether latest is a higher version than current. A current
// version that is not a plain x.y.z (a development build) is never updated.
func Newer(current, latest string) bool {
	cv, ok1 := parseVersion(current)
	lv, ok2 := parseVersion(latest)
	if !ok1 || !ok2 {
		return false
	}
	for i := range cv {
		if lv[i] != cv[i] {
			return lv[i] > cv[i]
		}
	}
	return false
}

func parseVersion(s string) (v [3]int, ok bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	if i := strings.IndexAny(s, "-+ ("); i >= 0 { // "1.3.1-ci", "1.3.1 (abc)"
		s = s[:i]
	}
	if !versionRe.MatchString(s) {
		return v, false
	}
	for i, p := range strings.Split(s, ".") {
		n, err := strconv.Atoi(p)
		if err != nil {
			return v, false
		}
		v[i] = n
	}
	return v, true
}

// parseSums reads "<sha256>  <name>" lines (sha256sum format).
func parseSums(text string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) == 2 && len(f[0]) == 64 {
			out[strings.TrimPrefix(f[1], "*")] = strings.ToLower(f[0])
		}
	}
	return out
}

// Download fetches the installer of r into dir and returns its path. The file
// only gets its final name after its SHA-256 matches the release's checksum
// file; otherwise it is deleted and an error is returned.
func (c *Client) Download(ctx context.Context, r Release, dir string) (string, error) {
	sums, err := c.get(ctx, r.SumsURL, "text/plain", maxSums)
	if err != nil {
		return "", fmt.Errorf("checksums: %w", err)
	}
	want, ok := parseSums(string(sums))[r.InstallerName]
	if !ok {
		return "", fmt.Errorf("no checksum for %s", r.InstallerName)
	}

	u, err := url.Parse(r.InstallerURL)
	if err != nil || !c.allow(u) || !installerName.MatchString(r.InstallerName) {
		return "", errors.New("installer url is not allowed")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.InstallerURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download: %s", resp.Status)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	final := filepath.Join(dir, r.InstallerName)
	part := final + ".part"
	f, err := os.OpenFile(part, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(resp.Body, maxInstaller+1))
	cerr := f.Close()
	switch {
	case err != nil:
		os.Remove(part)
		return "", err
	case cerr != nil:
		os.Remove(part)
		return "", cerr
	case n > maxInstaller || n < minInstaller:
		os.Remove(part)
		return "", fmt.Errorf("installer has an implausible size (%d bytes)", n)
	case hex.EncodeToString(h.Sum(nil)) != want:
		os.Remove(part)
		return "", errors.New("installer checksum does not match the release; discarded")
	}
	os.Remove(final)
	if err := os.Rename(part, final); err != nil {
		return "", err
	}
	return final, nil
}

// Info is the result of a check, shared by the settings windows.
type Info struct {
	Available bool   `json:"available"`
	Current   string `json:"current"`
	Latest    string `json:"latest"`
	Notes     string `json:"notes"`
}
