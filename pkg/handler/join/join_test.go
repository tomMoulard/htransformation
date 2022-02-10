package join_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tommoulard/htransformation/pkg/handler/join"
	"github.com/tommoulard/htransformation/pkg/types"
)

func TestJoinHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
	}{
		{
			name: "Join two headers simple value",
			rule: types.Rule{
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"Tested",
				},
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested",
			},
		},
		{
			name: "Join two headers multiple value",
			rule: types.Rule{
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"Tested",
					"Compiled",
					"Working",
				},
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested,Compiled,Working",
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

			err = join.Handle(nil, req, test.rule)
			require.NoError(t, err)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}
