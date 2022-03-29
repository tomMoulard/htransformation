package htransformation_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	plug "github.com/tommoulard/htransformation"
	"github.com/tommoulard/htransformation/pkg/types"
)

func TestValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *plug.Config
		wantErr bool
	}{
		{
			name:    "no rules",
			config:  &plug.Config{},
			wantErr: false,
		},
		{
			name: "no rules type",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name: "no rule",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid rules type",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name: "invalid rule",
						Type: "THIS IS NOT A VALID RULE TYPE",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing header",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name: "rule with no header",
						Type: types.Join,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing type",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name:   "rule with no type",
						Header: "not-empty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "join rule without value",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name:   "join rule with no value",
						Header: "not-empty",
						Sep:    "not-empty",
						Type:   types.Join,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "join rule without separator",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name:   "join rule with no sep",
						Header: "not-empty",
						Value:  "not-empty",
						Type:   types.Join,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rewrite rule without ValueReplace",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name:   "rewrite rule with no ValueReplace",
						Type:   types.RewriteValueRule,
						Header: "not-empty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid rule",
			config: &plug.Config{
				Rules: []types.Rule{
					{
						Name:   "delete rule",
						Header: "not-empty",
						Type:   types.Delete,
					},
					{
						Name:   "join Rule",
						Header: "not-empty",
						Values: []string{"not-empty"},
						Sep:    "not-empty",
						Type:   types.Join,
					},
					{
						Name:   "rename rule",
						Header: "not-empty",
						Value:  "not-empty",
						Type:   types.Rename,
					},
					{
						Name:         "rewrite rule",
						Header:       "not-empty",
						ValueReplace: "not-empty",
						Value:        "not-empty",
						Type:         types.RewriteValueRule,
					},
					{
						Name:   "set rule",
						Header: "not-empty",
						Value:  "not-empty",
						Type:   types.Set,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := plug.New(context.Background(), nil, test.config, "test")
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHeaderRules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		rule             types.Rule
		wantErr          bool
		additionalHeader map[string]string
	}{
		{
			name: "set rule",
			rule: types.Rule{
				Name:   "set rule",
				Header: "not-empty",
				Value:  "not-empty",
				Type:   types.Set,
			},
			wantErr: false,
		},
		{
			name: "rename rule",
			rule: types.Rule{
				Name:   "rename rule",
				Header: "not-empty",
				Value:  "not-empty",
				Type:   types.Rename,
			},
			additionalHeader: map[string]string{
				"Referer": "http://foo.bar",
			},
			wantErr: false,
		},
		{
			name: "rewrite rule",
			rule: types.Rule{
				Name:         "rewrite rule",
				Header:       "not-empty",
				Value:        "not-empty",
				ValueReplace: "not-empty",
				Type:         types.RewriteValueRule,
			},
			additionalHeader: map[string]string{
				"Referer": "http://foo.bar",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg := plug.CreateConfig()
			cfg.Rules = []types.Rule{test.rule}

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, err := plug.New(ctx, next, cfg, "demo-plugin")
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for key, value := range test.additionalHeader {
				req.Header.Set(key, value)
			}

			handler.ServeHTTP(recorder, req)
			result := recorder.Result()
			statusCode := result.StatusCode
			require.NoError(t, result.Body.Close())

			if test.wantErr {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			} else {
				assert.Equal(t, http.StatusOK, statusCode)
			}
		})
	}
}
