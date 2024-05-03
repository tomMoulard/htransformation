package rename_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/handler/rename"
	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/tests/require"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestRenameHandler_Host(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		rule            types.Rule
		requestHeaders  map[string]string
		want            map[string]string
		expectedHost    string
		expectedURLHost string
	}{
		{
			name: "Rename Host to another",
			rule: types.Rule{
				Header: "Host",
				Value:  "Fake-Host",
			},
			requestHeaders: map[string]string{},
			want: map[string]string{
				"Fake-Host": "example.com",
			},
			expectedHost:    "",
			expectedURLHost: "example.com",
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
			want:            map[string]string{},
			expectedHost:    "example.org",
			expectedURLHost: "example.com",
		},
		{
			name: "Deletion",
			rule: types.Rule{
				Header: "Host",
			},
			requestHeaders:  map[string]string{},
			want:            map[string]string{},
			expectedHost:    "",
			expectedURLHost: "example.com",
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

			test.rule.Regexp = regexp.MustCompile(test.rule.Header)

			rename.Handle(nil, req, test.rule)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}

			assert.Equal(t, test.expectedHost, req.Host)
			assert.Equal(t, test.expectedURLHost, req.URL.Host)
		})
	}
}

func TestRenameHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
	}{
		{
			name: "no transformation",
			rule: types.Rule{
				Header: "not-existing",
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
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
			want: map[string]string{
				"Foo":       "Bar",
				"X-Testing": "Success",
			},
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
			want: map[string]string{
				"Foo":  "Bar",
				"Test": "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			test.rule.Regexp = regexp.MustCompile(test.rule.Header)

			rename.Handle(nil, req, test.rule)

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
			name: "missing header value",
			rule: types.Rule{
				Header: ".",
				Type:   types.Rename,
			},
			wantErr: true,
		},
		{
			name: "invalid regexp",
			rule: types.Rule{
				Header: "(",
				Type:   types.Rename,
			},
			wantErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Value:  "not-empty",
				Type:   types.Rename,
			},
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := rename.Validate(test.rule)
			t.Log(err)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
