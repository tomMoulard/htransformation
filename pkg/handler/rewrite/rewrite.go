package rewrite

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
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

func applyHeaderRule(headerName string, headerValues []string, rule types.Rule,
	req *http.Request, rw http.ResponseWriter,
) {
	switch {
	case rule.SetOnResponse:
		rw.Header().Del(headerName)

	case strings.EqualFold(headerName, "Host"):
		req.Host = ""

	default:
		req.Header.Del(headerName)
	}

	for _, headerValue := range headerValues {
		replacedValue := replaceHeaderValue(headerValue, rule)
		setHeader(headerName, replacedValue, rule, req, rw)
	}
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

func setHeader(headerName string, headerValue string, rule types.Rule, req *http.Request, rw http.ResponseWriter) {
	switch {
	case rule.SetOnResponse:
		rw.Header().Add(rule.Header, headerValue)

	case strings.EqualFold(headerName, "Host"):
		req.Host = headerValue

	default:
		req.Header.Add(headerName, headerValue)
	}
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

		applyHeaderRule(headerName, headerValues, rule, req, rw)
	}

	req.Header.Set("Host", originalHost)
}
