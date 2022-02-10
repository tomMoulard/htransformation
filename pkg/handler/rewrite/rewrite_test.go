package rewrite_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tommoulard/htransformation/pkg/handler/rewrite"
	"github.com/tommoulard/htransformation/pkg/types"
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

			err = rewrite.Handle(nil, req, test.rule)
			require.NoError(t, err)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}
