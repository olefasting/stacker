package stacker_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"

	"github.com/olefasting/httpctx"
	"github.com/olefasting/stacker"
)

type testHandler struct{}

func (h *testHandler) ServeHTTPCtx(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	return
}

func testTagger(tag string) stacker.Middleware {
	return func(h httpctx.Handler) httpctx.Handler {
		return httpctx.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(tag))
		})
	}
}

// Test that the stacker are immutable
func TestStackMutability(t *testing.T) {
	t.Parallel()

	// Create test stacker
	s1 := stacker.New(testTagger("1"), testTagger("2"), testTagger("3"))
	s2 := s1.Append(testTagger("4"), testTagger("5"))
	s3 := s2.Append(testTagger("6"))
	// Check tht stacker are different length
	assert.True(t, len(s1) != len(s2))
	assert.True(t, len(s2) != len(s3))
}

// Check that handler order is right
func TestHandlerOrder(t *testing.T) {
	t.Parallel()

	// Create test stacker
	s1 := stacker.New(testTagger("1"), testTagger("2"))
	s2 := s1.Append(testTagger("3"), testTagger("4"), testTagger("5"))
	// Create
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	// Chack that tag order match the strings
	ctx := context.Background()
	rw := httptest.NewRecorder()
	s1.Then(&testHandler{}).ServeHTTPCtx(ctx, rw, req)
	assert.True(t, rw.Body.String() == "12")
	rw = httptest.NewRecorder()
	s2.Then(&testHandler{}).ServeHTTPCtx(ctx, rw, req)
	assert.True(t, rw.Body.String() == "12345")
}

// Test passing nil to Then
func TestNilPassedToThen(t *testing.T) {
	t.Parallel()

	s := stacker.New()
	h := s.Then(nil)
	assert.True(t, h == nil)
}
