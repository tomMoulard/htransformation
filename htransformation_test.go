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
		{
			name: "[Rename] no transformation with ValueIsHeaderPrefix",
			rule: plug.Rule{
				Type:                "Rename",
				Header:              "not-existing",
				Value:               "^unused",
				ValueIsHeaderPrefix: "^",
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
				Type:                "Rename",
				Header:              "Test",
				Value:               "^X-Dest-Header",
				ValueIsHeaderPrefix: "^",
			},
			headers: map[string]string{
				"Foo":           "Bar",
				"Test":          "Success",
				"X-Dest-Header": "X-Testing",
			},
			want: map[string]string{
				"Foo":           "Bar",
				"X-Dest-Header": "X-Testing",
				"X-Testing":     "Success",
			},
		},
		{
			name: "[Set] new header from existing",
			rule: plug.Rule{
				Type:                "Set",
				Header:              "X-Test",
				Value:               "^X-Source",
				ValueIsHeaderPrefix: "^",
			},
			headers: map[string]string{
				"Foo":      "Bar",
				"X-Source": "SourceHeader",
			},
			want: map[string]string{
				"Foo":      "Bar",
				"X-Source": "SourceHeader",
				"X-Test":   "SourceHeader",
			},
		},
		{
			name: "[Set] existing header from another existing",
			rule: plug.Rule{
				Type:                "Set",
				Header:              "X-Test",
				Value:               "^X-Source",
				ValueIsHeaderPrefix: "^",
			},
			headers: map[string]string{
				"Foo":      "Bar",
				"X-Source": "SourceHeader",
				"X-Test":   "Initial",
			},
			want: map[string]string{
				"Foo":      "Bar",
				"X-Source": "SourceHeader",
				"X-Test":   "SourceHeader",
			},
		},
		{
			name: "[Join] Join two headers simple value",
			rule: plug.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"^X-Source",
				},
				ValueIsHeaderPrefix: "^",
			},
			headers: map[string]string{
				"Foo":      "Bar",
				"X-Source": "Tested",
				"X-Test":   "Bar",
			},
			want: map[string]string{
				"Foo":      "Bar",
				"X-Source": "Tested",
				"X-Test":   "Bar,Tested",
			},
		},
		{
			name: "[Join] Join two headers multiple value",
			rule: plug.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"^X-Source-1",
					"Compiled",
					"^X-Source-3",
				},
				ValueIsHeaderPrefix: "^",
			},
			headers: map[string]string{
				"Foo":        "Bar",
				"X-Test":     "Bar",
				"X-Source-1": "Tested",
				"X-Source-3": "Working",
			},
			want: map[string]string{
				"Foo":        "Bar",
				"X-Test":     "Bar,Tested,Compiled,Working",
				"X-Source-1": "Tested",
				"X-Source-3": "Working",
			},
		},
		{
			name: "[Join] Join two headers multiple value with itself",
			rule: plug.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"second",
					"^X-Test",
					"^X-Source-3",
				},
				ValueIsHeaderPrefix: "^",
			},
			headers: map[string]string{
				"Foo":        "Bar",
				"X-Test":     "test",
				"X-Source-3": "third",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "test,second,test,third",
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
