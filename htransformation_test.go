package htransformation_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	plug "github.com/tommoulard/htransformation"
)

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
			name: "[ValueRewriteRule] one transformation",
			rule: plug.Rule{
				Type:         "RewriteValueRule",
				Header:       "F(.*)",
				Value:        `X-(\d*)-(.*)`,
				ValueReplace: "Y-$2-$1",
			},
			headers: map[string]string{
				"Foo": "X-12-Test",
			},
			want: map[string]string{
				"Foo": "Y-Test-12",
			},
		},
		{
			// the value doesn't match, we leave the value as is
			name: "[ValueRewriteRule] no match",
			rule: plug.Rule{
				Type:         "RewriteValueRule",
				Header:       "F(.*)",
				Value:        `(\d*)`,
				ValueReplace: "Y-$2-$1",
			},
			headers: map[string]string{
				"Foo": "X-Test",
			},
			want: map[string]string{
				"Foo": "X-Test",
			},
		},
		{
			// no placeholder but the value matches, we replace the value
			name: "[ValueRewriteRule] no placeholder",
			rule: plug.Rule{
				Type:         "RewriteValueRule",
				Header:       "F(.*)",
				Value:        `X-(.*)`,
				ValueReplace: "Y-Bla",
			},
			headers: map[string]string{
				"Foo": "X-Test",
			},
			want: map[string]string{
				"Foo": "Y-Bla",
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
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for hName, hVal := range tt.headers {
				req.Header.Add(hName, hVal)
			}

			handler.ServeHTTP(recorder, req)

			for hName, hVal := range tt.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}
