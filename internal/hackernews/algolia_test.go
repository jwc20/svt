package hackernews

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcBonusHype(t *testing.T) {
	tests := []struct {
		name     string
		karma    int
		nbHits   int
		expected int
	}{
		{"karma 100, 0 front page hits", 100, 0, 10},
		{"karma 1000, 1 front page hit", 1000, 1, 17},
		{"karma 150000, 1 front page hit", 150000, 1, 27},
		{"zero karma, zero hits", 0, 0, 0},
		{"high karma and hits caps at 100", 1000000, 50, 100},
		{"karma 1, 0 hits", 1, 0, 0},    // floor(log10(2)) * 5 = 0*5 = 0
		{"karma 9, 0 hits", 9, 0, 5},    // floor(log10(10)) * 5 = 1*5 = 5
		{"karma 0, 5 hits", 0, 5, 10},   // floor(log10(1)) * 5 + 5*2 = 0 + 10
		{"karma 10, 0 hits", 10, 0, 5},  // floor(log10(11)) * 5 = 1*5 = 5
		{"karma 99, 0 hits", 99, 0, 10}, // floor(log10(100)) * 5 = 2*5 = 10
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalcBonusHype(tt.karma, tt.nbHits)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFetchBonusHype_EmptyUsername(t *testing.T) {
	assert.Equal(t, 0, FetchBonusHype(""))
}

func TestFetchBonusHype_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/users/testuser", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"username":"testuser","karma":1000}`)
	})
	mux.HandleFunc("/api/v1/search", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"nbHits":1}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Override package-level URLs via a helper that swaps the HTTP client
	origClient := httpClient
	origUserURL := userURL
	origSearchURL := searchURL
	defer func() {
		httpClient = origClient
	}()

	// Point the fetchers at the test server by replacing the http client
	// with a custom transport that rewrites requests to our test server.
	httpClient = server.Client()
	transport := &rewriteTransport{base: http.DefaultTransport, serverURL: server.URL}
	httpClient.Transport = transport

	// We need to temporarily override the URLs; use the transport rewrite approach instead.
	_ = origUserURL
	_ = origSearchURL

	bonus := FetchBonusHype("testuser")
	// karma=1000, nbHits=1 → floor(log10(1001))*5 + 1*2 = 3*5 + 2 = 17
	assert.Equal(t, 17, bonus)
}

func TestFetchBonusHype_UserNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/users/noone", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	origClient := httpClient
	defer func() { httpClient = origClient }()

	httpClient = server.Client()
	httpClient.Transport = &rewriteTransport{base: http.DefaultTransport, serverURL: server.URL}

	bonus := FetchBonusHype("noone")
	assert.Equal(t, 0, bonus)
}

// rewriteTransport rewrites request URLs to point at the test server.
type rewriteTransport struct {
	base      http.RoundTripper
	serverURL string
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = ""
	req = req.Clone(req.Context())
	req.URL, _ = req.URL.Parse(t.serverURL + req.URL.Path + "?" + req.URL.RawQuery)
	return t.base.RoundTrip(req)
}
