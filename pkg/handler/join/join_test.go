package join_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/handler/join"
	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/tests/require"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestJoinHandler(t *testing.T) {
	testCases := []struct {
		name            string
		rule            types.Rule
		requestHeaders  map[string]string
		expectedHeaders map[string]string
		expectedHost    string
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
			expectedHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested",
			},
			expectedHost: "example.com",
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
			expectedHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested,Compiled,Working",
			},
			expectedHost: "example.com",
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
			expectedHeaders: map[string]string{
				"Foo":      "Bar",
				"X-Source": "Tested",
				"X-Test":   "Bar,Tested",
			},
			expectedHost: "example.com",
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
			expectedHeaders: map[string]string{
				"Foo":        "Bar",
				"X-Test":     "Bar,Tested,Compiled,Working",
				"X-Source-1": "Tested",
				"X-Source-3": "Working",
			},
			expectedHost: "example.com",
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
			expectedHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "test,second,test,third",
			},
			expectedHost: "example.com",
		},
		{
			name: "Join value with same HeaderPrefix",
			rule: types.Rule{
				Sep:          ",",
				Header:       "X-Test",
				HeaderPrefix: "Tested",
				Values: []string{
					"Tested",
				},
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			expectedHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested",
			},
			expectedHost: "example.com",
		},
		{
			name: "Join Host header",
			rule: types.Rule{
				Sep:          ",",
				Header:       "Host",
				HeaderPrefix: "Tested",
				Values: []string{
					"Tested",
				},
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			expectedHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			expectedHost: "example.com,Tested",
		},
		{
			name: "Twice Host header",
			rule: types.Rule{
				Sep:    ",",
				Header: "Host",
				Values: []string{
					"^Host",
				},
				HeaderPrefix: "^",
			},
			expectedHost: "example.com,example.com",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/foo", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			joinHandler, err := join.New(test.rule)
			require.NoError(t, err)

			joinHandler.Handle(nil, req)

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
		name    string
		rule    types.Rule
		wantErr bool
	}{
		{
			name:    "no rules",
			wantErr: true,
		},
		{
			name: "missing header",
			rule: types.Rule{
				Type: types.Join,
			},
			wantErr: true,
		},
		{
			name: "without value",
			rule: types.Rule{
				Header: "not-empty",
				Sep:    "not-empty",
				Type:   types.Join,
			},
			wantErr: true,
		},
		{
			name: "join rule without separator",
			rule: types.Rule{
				Header: "not-empty",
				Value:  "not-empty",
				Type:   types.Join,
			},
			wantErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Values: []string{"not-empty"},
				Sep:    "not-empty",
				Type:   types.Join,
			},
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			joinHandler, err := join.New(test.rule)
			require.NoError(t, err)

			err = joinHandler.Validate()
			t.Log(err)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
