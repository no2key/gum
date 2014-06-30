package jekyll

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestJekyllHandler_ServeHTTP(t *testing.T) {
	dir, cleanup, err := newTestSite("")
	if err != nil {
		t.Fatalf("error creating test site: %v", err)
	}
	defer cleanup()

	// write post
	os.MkdirAll(path.Join(dir, "_posts"), 0755)
	postPath := path.Join(dir, "_posts", "2014-05-28-test.md")
	if err := ioutil.WriteFile(postPath, []byte("---\nwordpress_id: 100\n---\n"), 0644); err != nil {
		t.Fatalf("error creating test post: %v", err)
	}

	in := "/b/1f"

	handler, err := NewHandler("b", dir)
	if err != nil {
		t.Fatalf("error constructing handler: %v", err)
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", in, nil)
	if err != nil {
		t.Fatalf("error constructing request: %v", err)
	}

	handler.ServeHTTP(w, r)

	if got, want := w.Code, http.StatusMovedPermanently; got != want {
		t.Errorf("Response status code got %v, want %v", got, want)
	}

	loc := w.Header().Get("Location")
	if loc == "" {
		t.Errorf("No location header set for input: %q", in)
	}
	if got, want := loc, "/2014/05/28/test.html"; got != want {
		t.Errorf("Location header for input %q got: %v, want: %v", in, got, want)
	}
}