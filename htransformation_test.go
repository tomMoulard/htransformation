package htransformation_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	plug "github.com/tommoulard/htransformation"
)

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	h := req.Header.Get(key)
	if h != expected {
		t.Errorf("invalid header value, got '%s', expect: '%s'", h, expected)
	}
}

func TestOneHeader(t *testing.T) {
	tests := []struct {
		name			string
		transformations	plug.Transform
		headers			map[string]string
		want			map[string]string
	}{
		{
			name: "no transformation",
			transformations: plug.Transform{
				Rename: "not-existing",
			},
			headers: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "one transformation",
			transformations: plug.Transform{
				Rename: "Test",
				With:   "Testing",
			},
			headers: map[string]string{
				"Foo": "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo": "Bar",
				"Testing": "Success",
			},
		},
		{
			name: "more transformation",
			transformations: plug.Transform{
				Rename: "Test*",
				With:   "Testing",
			},
			headers: map[string]string{
				"Foo": "Bar",
				"Test1": "Success",
				"Test2": "Pass",
			},
			want: map[string]string{
				"Foo": "Bar",
				"Testing": "Pass",
			},
		},
		{
			name: "DEL",
			transformations: plug.Transform{
				Rename: "Test",
			},
			headers: map[string]string{
				"Foo": "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo": "Bar",
				"Test": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := plug.CreateConfig()
			cfg.Transformations = []plug.Transform{tt.transformations}

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, err := plug.New(ctx, next, cfg, "demo-plugin")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			for hName, hVal := range tt.headers {
				req.Header.Add(hName, hVal)
			}

			handler.ServeHTTP(recorder, req)

			for hName, hVal := range tt.want {
				assertHeader(t, req, hName, hVal)
			}
		})
	}
}
