package header_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

func TestSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		header          string
		value           string
		expectedHeaders map[string][]string
		expectedHost    string
	}{
		{
			name:   "Set header",
			header: "Foo",
			value:  "Bar",
			expectedHeaders: map[string][]string{
				"Foo": {"Bar"},
			},
			expectedHost: "example.com",
		},
		{
			name:         "Set Host header",
			header:       "Host",
			value:        "example.org",
			expectedHost: "example.org",
		},
		{
			name:         "Set empty header",
			header:       "",
			value:        "",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"": {""},
			},
		},
		{
			name:         "Set empty header with value",
			header:       "",
			value:        "Bar",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"": {"Bar"},
			},
		},
		{
			name:         "Set header with empty value",
			header:       "Foo",
			value:        "",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"Foo": {""},
			},
		},
		{
			name:         "Set Host header with empty value",
			header:       "Host",
			value:        "",
			expectedHost: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			header.Set(req, test.header, test.value)

			assert.Equal(t, test.expectedHost, req.Host)

			for hName, hVal := range req.Header {
				assert.Equalf(t, test.expectedHeaders[hName], hVal, "header %q", hName)
			}
		})
	}
}
