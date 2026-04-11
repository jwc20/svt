package rand_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	svtrand "github.com/jwc20/svt/internal/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRandomInt(t *testing.T) {
	t.Run("returns a random number between 1 and 100", func(t *testing.T) {
		result := svtrand.GetRandomInt(100)

		assert.GreaterOrEqual(t, result, 1, fmt.Sprintf("result %d should be >= 1", result))
		assert.LessOrEqual(t, result, 100, fmt.Sprintf("result %d should be <= 100", result))
	})
}

func TestNewGetRandomIntRequest(t *testing.T) {
	url := "https://www.random.org/integers/?num=1&min=1&max=10&col=1&base=10&format=plain&rnd=new"

	t.Run("returns a GET request", func(t *testing.T) {
		req := svtrand.NewGetRandomIntRequest(url)

		assert.Equal(t, http.MethodGet, req.Method)
	})

	t.Run("returns a non-nil request", func(t *testing.T) {
		req := svtrand.NewGetRandomIntRequest(url)

		require.NotNil(t, req)
	})
}

func TestNewGetRandomIntResponseFromClient(t *testing.T) {
	t.Run("returns a 200 response from a reachable server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "42")
		}))
		defer server.Close()

		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)

		resp, _ := svtrand.NewGetRandomIntResponseFromClient(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("response body contains expected content", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "42")
		}))
		defer server.Close()

		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)

		resp, _ := svtrand.NewGetRandomIntResponseFromClient(req)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "42", string(body))
	})
}

func TestExtractRandomInteger(t *testing.T) {
	makeResponse := func(body string) *http.Response {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader(body)),
		}
	}

	t.Run("parses a plain integer", func(t *testing.T) {
		resp := makeResponse("42")
		result := svtrand.ExtractRandomInteger(resp)

		assert.Equal(t, 42, result)
	})

	t.Run("trims trailing newline", func(t *testing.T) {
		resp := makeResponse("7\n")
		result := svtrand.ExtractRandomInteger(resp)

		assert.Equal(t, 7, result)
	})

	t.Run("trims surrounding whitespace", func(t *testing.T) {
		resp := makeResponse("  13  \n")
		result := svtrand.ExtractRandomInteger(resp)

		assert.Equal(t, 13, result)
	})
}
