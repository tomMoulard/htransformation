package set_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomMoulard/htransformation/pkg/handler/set"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestSetHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
	}{
		{
			name: "Set one simple",
			rule: types.Rule{
				Header: "X-Test",
				Value:  "Tested",
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Tested",
			},
		},
		{
			name: "Set already existing simple",
			rule: types.Rule{
				Header: "X-Test",
				Value:  "Tested",
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Tested", // Override
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

			set.Handle(nil, req, test.rule)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name   string
		rule   types.Rule
		expect assert.ErrorAssertionFunc
	}{
		{
			name:   "no rules",
			expect: assert.Error,
		},
		{
			name: "missing Header value",
			rule: types.Rule{
				Type: types.Set,
			},
			expect: assert.Error,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Type:   types.Set,
			},
			expect: assert.NoError,
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := set.Validate(test.rule)
			t.Log(err)
			test.expect(t, err)
		})
	}
}
