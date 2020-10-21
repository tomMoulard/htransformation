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

func TestHeaderRules(t *testing.T) {
	tests := []struct {
		name    string
		rule    plug.Rule
		headers map[string]string
		want    map[string]string
	}{
		{
			name: "[Rename] no transformation",
			rule: plug.Rule{
				Type:   "Rename",
				Header: "not-existing",
			},
			headers: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "[Rename] one transformation",
			rule: plug.Rule{
				Type:   "Rename",
				Header: "Test",
				Value:  "X-Testing",
			},
			headers: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo":       "Bar",
				"X-Testing": "Success",
			},
		},
		{
			name: "[Rename] Deletion",
			rule: plug.Rule{
				Type:   "Rename",
				Header: "Test",
			},
			headers: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo":  "Bar",
				"Test": "",
			},
		},
		{
			name: "[Set] Set one simple",
			rule: plug.Rule{
				Type:   "Set",
				Header: "X-Test",
				Value:  "Tested",
			},
			headers: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Tested",
			},
		},
		{
			name: "[Set] Set already existing simple",
			rule: plug.Rule{
				Type:   "Set",
				Header: "X-Test",
				Value:  "Tested",
			},
			headers: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Tested", // Override
			},
		},
		{
			name: "[Del] Remove not existing header",
			rule: plug.Rule{
				Type:   "Del",
				Header: "X-Test",
			},
			headers: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "[Del] Remove one header",
			rule: plug.Rule{
				Type:   "Del",
				Header: "X-Test",
			},
			headers: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "[Join] Join two headers simple value",
			rule: plug.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"Tested",
				},
			},
			headers: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested",
			},
		},
		{
			name: "[Join] Join two headers multiple value",
			rule: plug.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"Tested",
					"Compiled",
					"Working",
				},
			},
			headers: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested,Compiled,Working",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := plug.CreateConfig()
			cfg.Rules = []plug.Rule{tt.rule}

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

/*

func TestHeaderDeletion(t *testing.T) {
	tests := []struct {
		name    string
		del     plug.Del
		headers map[string]string
		want    map[string]string
	}{
		{
			name: "remove not existing header",
			del: plug.Del{
				Header: "X-Test",
			},
			headers: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "remove one header",
			del: plug.Del{
				Header: "X-Test",
			},
			headers: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "remove all header",
			del: plug.Del{
				Header: "",
			},
			headers: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := plug.CreateConfig()
			cfg.Deletions = []plug.Del{tt.del}

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
*/
