package rewrite_test

import (
	"net/http"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/handler/rewrite"
	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/tests/require"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestRewriteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		rule            types.Rule
		requestHeaders  map[string]string
		expectedHeaders map[string]string
		expectedHost    string
	}{
		{
			name: "one transformation",
			rule: types.Rule{
				Header:       "F(.*)",
				Value:        `X-(\d*)-(.*)`,
				ValueReplace: "Y-$2-$1",
			},
			requestHeaders: map[string]string{
				"Foo": "X-12-Test",
			},
			expectedHeaders: map[string]string{
				"Foo": "Y-Test-12",
			},
			expectedHost: "example.com",
		},
		{
			name: "one transformation with 2 headers",
			rule: types.Rule{
				Header:       "F(.*)",
				Value:        `X-(\d*)-(.*)`,
				ValueReplace: "Y-$2-$1",
			},
			requestHeaders: map[string]string{
				"Bar": "Baz",
				"Foo": "X-12-Test",
			},
			expectedHeaders: map[string]string{
				"Bar": "Baz",
				"Foo": "Y-Test-12",
			},
			expectedHost: "example.com",
		},
		{
			// the value doesn't match, we leave the value as is
			name: "no match",
			rule: types.Rule{
				Header:       "F(.*)",
				Value:        `(\d*)`,
				ValueReplace: "Y-$2-$1",
			},
			requestHeaders: map[string]string{
				"Foo": "X-Test",
			},
			expectedHeaders: map[string]string{
				"Foo": "X-Test",
			},
			expectedHost: "example.com",
		},
		{
			name: "no header match and no value match",
			rule: types.Rule{
				Header:       "F(.*)",
				Value:        `toto`,
				ValueReplace: "Y-$2-$1",
			},
			requestHeaders: map[string]string{
				"Foo": "X-Test",
			},
			expectedHeaders: map[string]string{
				"Foo": "X-Test",
			},
			expectedHost: "example.com",
		},
		{
			// no placeholder but the value matches, we replace the value
			name: "no placeholder",
			rule: types.Rule{
				Header:       "F(.*)",
				Value:        `X-(.*)`,
				ValueReplace: "Y-Bla",
			},
			requestHeaders: map[string]string{
				"Foo": "X-Test",
			},
			expectedHeaders: map[string]string{
				"Foo": "Y-Bla",
			},
			expectedHost: "example.com",
		},
		{
			name: "Host header transformation",
			rule: types.Rule{
				Header:       "Host",
				Value:        `(.+).com`,
				ValueReplace: "$1.org",
			},
			expectedHost: "example.org",
		},
		{
			name: "multiple replacements in single header value",
			rule: types.Rule{
				Header:       "Foo",
				Value:        "X-(\\d+)-(\\w+)",
				ValueReplace: "Y-$2-$1",
			},
			requestHeaders: map[string]string{
				"Foo": "X-12-Test;X-34-Prod",
			},
			expectedHeaders: map[string]string{
				"Foo": "Y-Test-12;Y-Prod-34",
			},
			expectedHost: "example.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/foo", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			rewriteHandler, err := rewrite.New(test.rule)
			require.NoError(t, err)

			rewriteHandler.Handle(nil, req)

			for hName, hVal := range test.expectedHeaders {
				actual := req.Header.Get(hName)
				if test.name == "multiple replacements in single header value" {
					t.Logf("DEBUG: actual header value: %q", actual)
				}
				assert.Equal(t, hVal, actual)
			}

			assert.Equal(t, test.expectedHost, req.Host)
			assert.Equal(t, "example.com", req.URL.Host)
		})
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		rule            types.Rule
		wantNewErr      bool
		wantValidateErr bool
	}{
		{
			name:            "no rules",
			wantValidateErr: true,
		},
		{
			name: "missing ValueReplace value",
			rule: types.Rule{
				Type: types.RewriteValueRule,
			},
			wantValidateErr: true,
		},
		{
			name: "invalid Header regexp",
			rule: types.Rule{
				Header: "(",
				Type:   types.RewriteValueRule,
			},
			wantNewErr: true,
		},
		{
			name: "invalid Value regexp",
			rule: types.Rule{
				ValueReplace: "not-empty",
				Value:        "(",
				Type:         types.RewriteValueRule,
			},
			wantNewErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header:       "not-empty",
				ValueReplace: "not-empty",
				Value:        "not-empty",
				Type:         types.RewriteValueRule,
			},
			wantNewErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rewriteHandler, err := rewrite.New(test.rule)
			if test.wantNewErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)

			err = rewriteHandler.Validate()
			t.Log(err)

			if test.wantValidateErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
