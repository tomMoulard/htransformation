package rename_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/handler/rename"
	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/tests/require"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestRenameHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		rule            types.Rule
		requestHeaders  map[string]string
		expectedHeaders map[string]string
		expectedHost    string
	}{
		{
			name: "no transformation",
			rule: types.Rule{
				Header: "not-existing",
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
			name: "one transformation",
			rule: types.Rule{
				Header: "Test",
				Value:  "X-Testing",
			},
			requestHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			expectedHeaders: map[string]string{
				"Foo":       "Bar",
				"X-Testing": "Success",
			},
			expectedHost: "example.com",
		},
		{
			name: "Deletion",
			rule: types.Rule{
				Header: "Test",
			},
			requestHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			expectedHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "",
			},
			expectedHost: "example.com",
		},
		{
			name: "Rename Host to another",
			rule: types.Rule{
				Header: "Host",
				Value:  "Fake-Host",
			},
			requestHeaders: map[string]string{},
			expectedHeaders: map[string]string{
				"Fake-Host": "example.com",
			},
			expectedHost: "",
		},
		{
			name: "Rename another to Host",
			rule: types.Rule{
				Header: "Fake-Host",
				Value:  "Host",
			},
			requestHeaders: map[string]string{
				"Fake-Host": "example.org",
			},
			expectedHeaders: map[string]string{},
			expectedHost:    "example.org",
		},
		{
			name: "Deletion",
			rule: types.Rule{
				Header: "Host",
			},
			requestHeaders:  map[string]string{},
			expectedHeaders: map[string]string{},
			expectedHost:    "",
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

			renameHandler, err := rename.New(test.rule)
			require.NoError(t, err)

			renameHandler.Handle(nil, req)

			for hName, hVal := range test.expectedHeaders {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}

			assert.Equal(t, test.expectedHost, req.Host)
			assert.Equal(t, "example.com", req.URL.Host)
		})
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		rule            types.Rule
		wantNewErr      bool
		wantValidateErr bool
	}{
		{
			name:            "no rules",
			wantValidateErr: true,
		},
		{
			name: "missing header value",
			rule: types.Rule{
				Header: ".",
				Type:   types.Rename,
			},
			wantValidateErr: true,
		},
		{
			name: "invalid regexp",
			rule: types.Rule{
				Header: "(",
				Type:   types.Rename,
			},
			wantNewErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Value:  "not-empty",
				Type:   types.Rename,
			},
			wantNewErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			renameHandler, err := rename.New(test.rule)
			if test.wantNewErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)

			err = renameHandler.Validate()
			t.Log(err)

			if test.wantValidateErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
