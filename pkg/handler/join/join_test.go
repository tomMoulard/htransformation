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

func TestJoinHandler_Host(t *testing.T) {
	testCases := []struct {
		name            string
		rule            types.Rule
		expectedHost    string
		expectedURLHost string
	}{
		// Wired test cases...
		{
			name: "Join two headers simple value",
			rule: types.Rule{
				Sep:    ",",
				Header: "Host",
				Values: []string{
					"example.org",
				},
			},
			expectedHost:    "example.com,example.org",
			expectedURLHost: "example.com",
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
			expectedHost:    "example.com,example.com",
			expectedURLHost: "example.com",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/foo", nil)
			require.NoError(t, err)

			join.Handle(nil, req, test.rule)

			assert.Equal(t, test.expectedHost, req.Host)
			assert.Equal(t, test.expectedURLHost, req.URL.Host)
		})
	}
}

func TestJoinHandler(t *testing.T) {
	testCases := []struct {
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
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar,Tested",
			},
		},
	}

	for _, test := range testCases {
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

			err := join.Validate(test.rule)
			t.Log(err)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
