package set_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/handler/set"
	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/tests/require"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestSetHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		wantOnRequest  map[string]string
		wantOnResponse map[string]string
		expectedHost   string
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
			wantOnRequest: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Tested",
			},
			expectedHost: "example.com",
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
			wantOnRequest: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Tested", // Override
			},
			expectedHost: "example.com",
		},
		{
			name: "Set on response",
			rule: types.Rule{
				Header:        "X-Test",
				Value:         "Tested",
				SetOnResponse: true,
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			wantOnRequest: map[string]string{
				"Foo": "Bar",
			},
			wantOnResponse: map[string]string{
				"X-Test": "Tested",
			},
			expectedHost: "example.com",
		},
		{
			name: "Set Host header",
			rule: types.Rule{
				Header: "Host",
				Value:  "example.org",
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

			rw := httptest.NewRecorder()
			set.Handle(rw, req, test.rule)

			for hName, hVal := range test.wantOnRequest {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}

			for hName, hVal := range test.wantOnResponse {
				assert.Equal(t, hVal, rw.Header().Get(hName))
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
			name: "missing Header value",
			rule: types.Rule{
				Type: types.Set,
			},
			wantErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Type:   types.Set,
			},
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := set.Validate(test.rule)
			t.Log(err)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
