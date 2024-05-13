package rewrite_test

import (
	"context"
	"net/http"
	"regexp"
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/foo", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			test.rule.Regexp = regexp.MustCompile(test.rule.Header)

			rewrite.Handle(nil, req, test.rule)

			for hName, hVal := range test.expectedHeaders {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}

			assert.Equal(t, test.expectedHost, req.Host)
			assert.Equal(t, "example.com", req.URL.Host)
		})
	}
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name      string
		rule      types.Rule
		expectErr bool
	}{
		{
			name:      "no rules",
			expectErr: true,
		},
		{
			name: "missing ValueReplace value",
			rule: types.Rule{
				Type: types.RewriteValueRule,
			},
			expectErr: true,
		},
		{
			name: "invalid Header regexp",
			rule: types.Rule{
				Header: "(",
				Type:   types.RewriteValueRule,
			},
			expectErr: true,
		},
		{
			name: "invalid Value regexp",
			rule: types.Rule{
				ValueReplace: "not-empty",
				Value:        "(",
				Type:         types.RewriteValueRule,
			},
			expectErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header:       "not-empty",
				ValueReplace: "not-empty",
				Value:        "not-empty",
				Type:         types.RewriteValueRule,
			},
			expectErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := rewrite.Validate(test.rule)
			t.Log(err)

			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
