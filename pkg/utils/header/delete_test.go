package header_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomMoulard/htransformation/pkg/tests/assert"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		header          string
		expectedHeaders map[string][]string
		expectedHost    string
	}{
		{
			name:   "Delete header",
			header: "Foo",
			expectedHeaders: map[string][]string{
				"Foo": {""},
			},
			expectedHost: "example.com",
		},
		{
			name:         "Delete Host header",
			header:       "Host",
			expectedHost: "",
			expectedHeaders: map[string][]string{
				"Foo": {"Bar"},
			},
		},
		{
			name:         "Delete empty header",
			header:       "",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"Foo": {"Bar"},
			},
		},
		{
			name:         "Delete header with empty value",
			header:       "Foo",
			expectedHost: "example.com",
			expectedHeaders: map[string][]string{
				"Foo": {"Bar"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			req.Header.Set("Foo", "Bar")

			header.Delete(req, test.header)

			assert.Equal(t, test.expectedHost, req.Host)

			for hName, hVal := range req.Header {
				assert.Equalf(t, test.expectedHeaders[hName], hVal, "header %q", hName)
			}
		})
	}
}
