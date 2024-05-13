package header_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		header          string
		value           string
		expectedHeaders map[string][]string
		expectedHost    string
	}{
		{
			name:   "Add header",
			header: "Foo",
			value:  "Bar",
			expectedHeaders: map[string][]string{
				"Foo": {"Bar"},
			},
			expectedHost: "example.com",
		},
		{
			name:            "Add Host header",
			header:          "Host",
			value:           "example.org",
			expectedHost:    "example.comexample.org",
			expectedHeaders: map[string][]string{},
		},
		{
			name:         "Add empty header",
			header:       "",
			value:        "",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"": {""},
			},
		},
		{
			name:         "Add empty header with value",
			header:       "",
			value:        "Bar",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"": {"Bar"},
			},
		},
		{
			name:         "Add header with empty value",
			header:       "Foo",
			value:        "",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"Foo": {""},
			},
		},
		{
			name:            "Add Host header with empty value",
			header:          "Host",
			value:           "",
			expectedHost:    "example.com",
			expectedHeaders: map[string][]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			header.Add(req, test.header, test.value)

			assert.Equal(t, test.expectedHost, req.Host)

			for hName, hVal := range req.Header {
				assert.Equalf(t, test.expectedHeaders[hName], hVal, "header %q", hName)
			}
		})
	}
}
