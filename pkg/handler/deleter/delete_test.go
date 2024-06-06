package deleter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/handler/deleter"
	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/tests/require"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestDeleteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		rule            types.Rule
		requestHeaders  map[string]string
		expectedHeaders map[string]string
		expectedHost    string
	}{
		{
			name: "Remove not existing header",
			rule: types.Rule{
				Header: "X-Test",
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			expectedHeaders: map[string]string{
				"Foo": "Bar",
			},
			expectedHost: "example.com",
		},
		{
			name: "Remove one header",
			rule: types.Rule{
				Header: "X-Test",
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			expectedHeaders: map[string]string{
				"Foo": "Bar",
			},
			expectedHost: "example.com",
		},
		{
			name: "Remove host header",
			rule: types.Rule{
				Header: "Host",
			},
			expectedHost: "",
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

			deleteHandler, err := deleter.New(test.rule)
			require.NoError(t, err)

			deleteHandler.Handle(nil, req)

			for hName, hVal := range test.expectedHeaders {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}

			assert.Equal(t, test.expectedHost, req.Host)
			assert.Equal(t, "example.com", req.URL.Host)
		})
	}
}

func TestDeleteHandlerOnResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
	}{
		{
			name: "Remove not existing header",
			rule: types.Rule{
				Header:        "X-Test",
				SetOnResponse: true,
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "Remove one header",
			rule: types.Rule{
				Header:        "X-Test",
				SetOnResponse: true,
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rw := httptest.NewRecorder()

			for hName, hVal := range test.requestHeaders {
				rw.Header().Add(hName, hVal)
			}

			deleteHandler, err := deleter.New(test.rule)
			require.NoError(t, err)

			deleteHandler.Handle(rw, nil)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, rw.Header().Get(hName))
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
			wantErr: false,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Type:   types.Delete,
			},
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			deleteHandler, err := deleter.New(test.rule)
			require.NoError(t, err)

			err = deleteHandler.Validate()
			t.Log(err)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
