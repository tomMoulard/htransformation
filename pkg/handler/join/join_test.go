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
		{
			name: "Join two headers simple value",
			rule: types.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"^X-Source",
				},
				HeaderPrefix: "^",
			},
			requestHeaders: map[string]string{
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
			name: "Join two headers multiple value",
			rule: types.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"^X-Source-1",
					"Compiled",
					"^X-Source-3",
				},
				HeaderPrefix: "^",
			},
			requestHeaders: map[string]string{
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
			name: "Join two headers multiple value with itself",
			rule: types.Rule{
				Type:   "Join",
				Sep:    ",",
				Header: "X-Test",
				Values: []string{
					"second",
					"^X-Test",
					"^X-Source-3",
				},
				HeaderPrefix: "^",
			},
			requestHeaders: map[string]string{
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

			join.Handle(nil, req, test.rule)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}
