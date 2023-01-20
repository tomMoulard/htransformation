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
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
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
			want: map[string]string{
				"Foo": "Y-Test-12",
			},
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
			want: map[string]string{
				"Bar": "Baz",
				"Foo": "Y-Test-12",
			},
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
			want: map[string]string{
				"Foo": "X-Test",
			},
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
			want: map[string]string{
				"Foo": "X-Test",
			},
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
			want: map[string]string{
				"Foo": "Y-Bla",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			test.rule.Regexp = regexp.MustCompile(test.rule.Header)

			rewrite.Handle(nil, req, test.rule)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name    string
		rule    types.Rule
		wantErr bool
	}{
		{
			name:    "no rules",
			wantErr: true,
		},
		{
			name: "missing ValueReplace value",
			rule: types.Rule{
				Type: types.RewriteValueRule,
			},
			wantErr: true,
		},
		{
			name: "invalid Header regexp",
			rule: types.Rule{
				Header: "(",
				Type:   types.RewriteValueRule,
			},
			wantErr: true,
		},
		{
			name: "invalid Value regexp",
			rule: types.Rule{
				ValueReplace: "not-empty",
				Value:        "(",
				Type:         types.RewriteValueRule,
			},
			wantErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header:       "not-empty",
				ValueReplace: "not-empty",
				Value:        "not-empty",
				Type:         types.RewriteValueRule,
			},
			wantErr: false,
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := rewrite.Validate(test.rule)
			t.Log(err)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
