package gum

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		prefix, dest string
		in, out      string
	}{
		// common cases
		{
			prefix: "x", dest: "http://example/",
			in: "/x", out: "http://example/",
		},
		{
			prefix: "x", dest: "http://example/",
			in: "/x/", out: "http://example/",
		},
		{
			prefix: "x", dest: "http://example/",
			in: "/x/y", out: "http://example/y",
		},
		{
			prefix: "x", dest: "http://example/",
			in: "/x/y?a=b", out: "http://example/y?a=b",
		},

		// absolute input URL (rare)
		{
			prefix: "x", dest: "http://example/",
			in: "http://foo/x/y", out: "http://example/y",
		},

		// destination URL with path
		{
			prefix: "x", dest: "http://example/a/",
			in: "/x", out: "http://example/a/",
		},
		{
			prefix: "x", dest: "http://example/a/",
			in: "/x/", out: "http://example/a/",
		},
		{
			prefix: "x", dest: "http://example/a/",
			in: "/x/y", out: "http://example/a/y",
		},

		// destination URL with path (no trailing slash)
		{
			prefix: "x", dest: "http://example/a",
			in: "/x", out: "http://example/a",
		},
		{
			prefix: "x", dest: "http://example/a",
			in: "/x/", out: "http://example/a",
		},
		{
			prefix: "x", dest: "http://example/a",
			in: "/x/y", out: "http://example/y",
		},

		// relative destination URL
		{prefix: "x", dest: "/a/", in: "/x", out: "/a/"},
		{prefix: "x", dest: "/a/", in: "/x/", out: "/a/"},
		{prefix: "x", dest: "/a/", in: "/x/y", out: "/a/y"},

		// no prefix
		{
			prefix: "", dest: "http://example/",
			in: "/x", out: "http://example/x",
		},
		{
			prefix: "", dest: "http://example/",
			in: "/x/", out: "http://example/x/",
		},
		{
			prefix: "", dest: "http://example/",
			in: "/x/y", out: "http://example/x/y",
		},

		// no destination (redirects to root)
		{prefix: "x", dest: "", in: "/x", out: "/"},
		{prefix: "x", dest: "", in: "/x/", out: "/"},
		{prefix: "x", dest: "", in: "/x/y", out: "/y"},
	}

	for _, tt := range tests {
		handler, err := NewRedirectHandler(tt.prefix, tt.dest)
		if err != nil {
			t.Fatalf("error constructing handler: %v", err)
		}

		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", tt.in, nil)
		if err != nil {
			t.Fatalf("error constructing request: %v", err)
		}

		handler.ServeHTTP(w, r)

		if got, want := w.Code, http.StatusMovedPermanently; got != want {
			t.Errorf("Response status code got %v, want %v", got, want)
		}

		loc := w.Header().Get("Location")
		if loc == "" {
			t.Errorf("No location header set for input: %q", tt.in)
		}
		if got, want := loc, tt.out; got != want {
			t.Errorf("Location header for input %q got: %v, want: %v", tt.in, got, want)
		}
	}
}

// Test that RedirectHandler registers proper prefixes on mux Router.
func TestRedirectHandler_Register(t *testing.T) {
	tests := []struct {
		prefix string
		in     string
		match  bool
	}{
		{prefix: "x", in: "/x", match: true},
		{prefix: "x", in: "/x/", match: true},
		{prefix: "x", in: "/x/y", match: true},
		{prefix: "x", in: "/xy", match: false},
	}

	for _, tt := range tests {
		router := mux.NewRouter()
		handler, err := NewRedirectHandler(tt.prefix, "")
		if err != nil {
			t.Fatalf("error constructing handler: %v", err)
		}
		handler.Register(router)

		var routeMatch mux.RouteMatch
		req, err := http.NewRequest("GET", tt.in, nil)
		if err != nil {
			t.Errorf("error constructing request for %q: %v", tt.in, err)
		}

		if got, want := router.Match(req, &routeMatch), tt.match; got != want {
			t.Errorf("route match for %q found: %v, want %v", tt.in, got, want)
		}
	}
}