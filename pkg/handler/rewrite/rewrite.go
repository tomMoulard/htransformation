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

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	headers := req.Header
	if rule.SetOnResponse {
		headers = rw.Header()
	}

	for headerName, headerValues := range headers {
		if matched := rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		if rule.SetOnResponse {
			rw.Header().Del(headerName)
		} else {
			req.Header.Del(headerName)
		}

		for _, headerValue := range headerValues {
			replacedHeaderValue := rule.ValueReplace
			ruleValueRegexp := regexp.MustCompile(rule.Value)
			captures := ruleValueRegexp.FindStringSubmatch(headerValue)

			if len(captures) == 0 || captures[0] == "" {
				if rule.SetOnResponse {
					rw.Header().Set(rule.Header, replacedHeaderValue)
				} else {
					req.Header.Set(headerName, headerValue)
				}

				continue
			}

			for j, capture := range captures[1:] {
				placeholder := fmt.Sprintf("$%d", j+1)
				replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
			}

			if rule.SetOnResponse {
				rw.Header().Set(rule.Header, replacedHeaderValue)
			} else {
				req.Header.Set(headerName, replacedHeaderValue)
			}
		}
	}
}
