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

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	for headerName, headerValues := range req.Header {
		if matched := rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		req.Header.Del(headerName)

		for _, headerValue := range headerValues {
			replacedHeaderValue := rule.ValueReplace
			ruleValueRegexp := regexp.MustCompile(rule.Value)
			captures := ruleValueRegexp.FindStringSubmatch(headerValue)

			if len(captures) == 0 || captures[0] == "" {
				req.Header.Add(headerName, headerValue)

				continue
			}

			for j, capture := range captures[1:] {
				placeholder := fmt.Sprintf("$%d", j+1)
				replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
			}

			req.Header.Add(headerName, replacedHeaderValue)
		}
	}
}
