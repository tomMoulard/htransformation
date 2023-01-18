package rename

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Validate(rule types.Rule) error {
	if _, err := regexp.Compile(rule.Header); err != nil {
		return fmt.Errorf("%s: %w", types.ErrInvalidRegexp.Error(), err)
	}

	if rule.Value == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	for headerName, headerValues := range req.Header {
		if matched := rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		req.Header.Del(headerName)

		for _, val := range headerValues {
			req.Header.Set(rule.Value, val)
		}
	}
}
