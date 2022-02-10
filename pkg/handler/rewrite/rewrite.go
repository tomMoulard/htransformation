package rewrite

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tommoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) error {
	for headerName, headerValues := range req.Header {
		matched, err := regexp.Match(rule.Header, []byte(headerName))
		if err != nil {
			return fmt.Errorf("RewriteValueHandler error: %w", err)
		}

		if !matched {
			continue
		}

		req.Header.Del(headerName)

		for _, headerValue := range headerValues {
			replacedHeaderValue := rule.ValueReplace
			ruleValueRegexp := regexp.MustCompile(rule.Value)
			captures := ruleValueRegexp.FindStringSubmatch(headerValue)

			if len(captures) == 0 || captures[0] == "" {
				req.Header.Add(headerName, headerValue)
			}

			for j, capture := range captures[1:] {
				placeholder := fmt.Sprintf("$%d", j+1)
				replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
			}

			req.Header.Add(headerName, replacedHeaderValue)
		}
	}

	return nil
}
