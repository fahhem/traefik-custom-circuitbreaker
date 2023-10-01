package circuitbreaker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fahhem/traefik-custom-circuitbreaker"
)

func TestDemo(t *testing.T) {
	cfg := circuitbreaker.CreateConfig()
	cfg.Expression = "NetworkErrorRatio() >= 0.0 || NetworkErrorRatio() <= 0.0"
  cfg.ResponseCode = 204

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := circuitbreaker.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	// First one always goes through
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	assertEqual(t, recorder.Code, 204)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
