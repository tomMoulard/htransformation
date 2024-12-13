package rewrite

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

type Rewrite struct {
	rule            *types.Rule
	ruleValueRegexp *regexp.Regexp
}

func New(rule types.Rule) (types.Handler, error) {
	re, err := regexp.Compile(rule.Header)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %q", types.ErrInvalidRegexp, rule.Name, rule.Header)
	}

	rule.Regexp = re

	re, err = regexp.Compile(rule.Value)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %q", types.ErrInvalidRegexp, rule.Name, rule.Value)
	}

	return &Rewrite{
		rule:            &rule,
		ruleValueRegexp: re,
	}, nil
}

func (r *Rewrite) Validate() error {
	if r.rule.ValueReplace == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func (r *Rewrite) replaceHeaderValue(headerValue string) string {
	replacedHeaderValue := r.rule.ValueReplace
	captures := r.ruleValueRegexp.FindStringSubmatch(headerValue)

	if len(captures) == 0 || captures[0] == "" {
		return headerValue
	}

	for j, capture := range captures[1:] {
		placeholder := fmt.Sprintf("$%d", j+1)
		replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
	}

	return replacedHeaderValue
}

func (r *Rewrite) Handle(rw http.ResponseWriter, req *http.Request) {
	headers := req.Header
	if r.rule.SetOnResponse {
		headers = rw.Header()
	}

	originalHost := req.Header.Get("Host") // Eventually X-Forwarded-Host
	req.Header.Set("Host", req.Host)

	for headerName, headerValues := range headers {
		if matched := r.rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		if r.rule.SetOnResponse {
			rw.Header().Del(headerName)
		} else {
			header.Delete(req, headerName)
		}

		for _, headerValue := range headerValues {
			replacedValue := r.replaceHeaderValue(headerValue)
			if r.rule.SetOnResponse {
				rw.Header().Add(headerName, replacedValue)
			} else {
				header.Add(req, headerName, replacedValue)
			}
		}
	}

	req.Header.Set("Host", originalHost)
}
