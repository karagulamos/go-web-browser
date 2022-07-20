package browser

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karagulamos/go-web-browser/internal/renderer"
	"github.com/stretchr/testify/require"
)

func TestBrowser_SiteRequest(t *testing.T) {
	httpserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello, world!</body></html>"))
	}))
	defer httpserver.Close()

	renderer := NewRequest(httpserver.URL)

	result := renderer.Invoke()
	require.NoError(t, result.Error)
	require.Equal(t, "Hello, world!", string(result.Content))
}

func TestBrowser_ViewSource(t *testing.T) {
	httpserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello, world!</body></html>"))
	}))
	defer httpserver.Close()

	renderer := NewRequest("view-source:" + httpserver.URL)

	result := renderer.Invoke()
	require.NoError(t, result.Error)

	require.Equal(t, "<html><body>Hello, world!</body></html>", string(result.Content))
}

func TestBrowser_FileUrl(t *testing.T) {
	r := NewRequest("file:///")

	result := r.Invoke()
	require.NoError(t, result.Error)
	require.NotEmpty(t, result.Content)
}

func TestBrowser_DataURI(t *testing.T) {
	testCases := []struct {
		name    string
		url     string
		content string
		err     error
	}{
		{
			name: "Invalid data URI",
			url:  "data:text/html",
			err:  renderer.ErrInvalidDataURI,
		},
		{
			name: "Unsupported content type",
			url:  "data:base64,PGh0bWw+PGJvZHk+SGVsbG8sIHdvcmxkPC9ib2R5PjwvaHRtbD4=",
			err:  renderer.ErrUnsupportedContentType,
		},
		{
			name:    "Valid data URI",
			url:     "data:text/html,<html><body>Hello, world!</body></html>",
			content: "Hello, world!",
			err:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest(tc.url)

			result := r.Invoke()
			require.Equal(t, tc.err, result.Error)
			require.Equal(t, tc.content, string(result.Content))
		})
	}

}

func TestBrowser_Redirects(t *testing.T) {
	var httpservers [MaxHttpRedirects + 1]*httptest.Server

	t.Cleanup(func() {
		for _, httpserver := range httpservers {
			httpserver.Close()
		}
	})

	httpservers[0] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello, world!</body></html>"))
	}))

	for i := 1; i < MaxHttpRedirects+1; i++ {
		httpservers[i] = func(i int) *httptest.Server {
			return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", httpservers[i-1].URL)
				w.WriteHeader(http.StatusFound)
			}))
		}(i)

	}

	testCases := []struct {
		name      string
		redirects int
		content   string
		err       error
	}{
		{
			name:      "regular redirect",
			redirects: 1,
			content:   "Hello, world!",
			err:       nil,
		},
		{
			name:      "max allowed redirects",
			redirects: MaxHttpRedirects,
			content:   "Hello, world!",
			err:       nil,
		},
		{
			name:      "too many redirects",
			redirects: MaxHttpRedirects + 1,
			content:   "",
			err:       renderer.ErrTooManyRedirects,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest(httpservers[tc.redirects-1].URL)

			result := r.Invoke()
			require.Equal(t, tc.err, result.Error)
			require.Equal(t, tc.content, string(result.Content))
		})
	}
}

func TestBrowser_Unimplemented(t *testing.T) {
	r := NewRequest("ftp://unimplemented.com")

	result := r.Invoke()
	require.Error(t, result.Error)
	require.Equal(t, renderer.ErrNotImplemented, result.Error)
}
