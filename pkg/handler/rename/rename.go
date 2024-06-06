package rename

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/tomMoulard/htransformation/pkg/types"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

type Rename struct {
	rule *types.Rule
}

func New(rule types.Rule) (types.Handler, error) {
	re, err := regexp.Compile(rule.Header)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", types.ErrInvalidRegexp, rule.Name)
	}

	rule.Regexp = re

	return &Rename{rule: &rule}, nil
}

func (r *Rename) Validate() error {
	if r.rule.Value == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func (r *Rename) Handle(rw http.ResponseWriter, req *http.Request) {
	originalHost := req.Header.Get("Host") // Eventually X-Forwarded-Host
	req.Header.Set("Host", req.Host)

	for headerName, headerValues := range req.Header {
		if matched := r.rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		if r.rule.SetOnResponse {
			rw.Header().Del(headerName)
		} else {
			header.Delete(req, headerName)
		}

		for _, val := range headerValues {
			if r.rule.SetOnResponse {
				rw.Header().Set(r.rule.Value, val)
			} else {
				header.Set(req, r.rule.Value, val)
			}
		}
	}

	req.Header.Set("Host", originalHost)
}
