package rewrite

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

func Validate(rule types.Rule) error {
	if _, err := regexp.Compile(rule.Header); err != nil {
		return fmt.Errorf("%s: %w", types.ErrInvalidRegexp.Error(), err)
	}

	if rule.ValueReplace == "" {
		return types.ErrMissingRequiredFields
	}

	if _, err := regexp.Compile(rule.Value); err != nil {
		return fmt.Errorf("%s: %w", types.ErrInvalidRegexp.Error(), err)
	}

	return nil
}

func replaceHeaderValue(headerValue string, rule types.Rule) string {
	replacedHeaderValue := rule.ValueReplace
	ruleValueRegexp := regexp.MustCompile(rule.Value)
	captures := ruleValueRegexp.FindStringSubmatch(headerValue)

	if len(captures) == 0 || captures[0] == "" {
		return headerValue
	}

	for j, capture := range captures[1:] {
		placeholder := fmt.Sprintf("$%d", j+1)
		replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
	}

	return replacedHeaderValue
}

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	headers := req.Header
	if rule.SetOnResponse {
		headers = rw.Header()
	}

	originalHost := req.Header.Get("Host") // Eventually X-Forwarded-Host
	req.Header.Set("Host", req.Host)

	for headerName, headerValues := range headers {
		if matched := rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		if rule.SetOnResponse {
			rw.Header().Del(headerName)
		} else {
			header.Delete(req, headerName)
		}

		for _, headerValue := range headerValues {
			replacedValue := replaceHeaderValue(headerValue, rule)
			if rule.SetOnResponse {
				rw.Header().Add(rule.Header, replacedValue)
			} else {
				header.Add(req, headerName, replacedValue)
			}
		}
	}

	req.Header.Set("Host", originalHost)
}
